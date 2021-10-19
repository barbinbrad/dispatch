#!/usr/bin/env python3
import json
import os
import queue
import random
import select
import socket
import threading
import time

from functools import partial
from typing import Any

from jsonrpc import JSONRPCResponseManager, dispatcher
from websocket import ABNF, WebSocketTimeoutException, WebSocketException, create_connection

HANDLER_THREADS = 4
WS_FRAME_SIZE = 4096

dispatcher["echo"] = lambda s: s
recv_queue: Any = queue.Queue()
send_queue: Any = queue.Queue()
log_send_queue: Any = queue.Queue()
log_recv_queue: Any = queue.Queue()


def startLocalProxy(global_end_event, remote_ws_uri, local_port):
  try:
    print("dispatch.startLocalProxy.starting")
    ws = create_connection(remote_ws_uri,
                           #cookie="jwt=" + identity_token,
                           enable_multithread=True)

    ssock, csock = socket.socketpair()
    local_sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    local_sock.connect(('127.0.0.1', local_port))
    local_sock.setblocking(0)

    proxy_end_event = threading.Event()
    threads = [
      threading.Thread(target=ws_proxy_recv, args=(ws, local_sock, ssock, proxy_end_event, global_end_event)),
      threading.Thread(target=ws_proxy_send, args=(ws, local_sock, csock, proxy_end_event))
    ]
    for thread in threads:
      thread.start()

    print("dispatch.startLocalProxy.started")
    return {"success": 1}
  except Exception as e:
    print("dispatchd.startLocalProxy.exception")
    raise e

def jsonrpc_handler(end_event):
    dispatcher["startLocalProxy"] = partial(startLocalProxy, end_event)
    while not end_event.is_set():
        try:
            data = recv_queue.get(timeout=1)

            if "method" in data:
                print(f"dispatch.jsonrpc_handler.call_method {data}")
                response = JSONRPCResponseManager.handle(data, dispatcher)
                send_queue.put_nowait(response.json)
            elif "id" in data and ("result" in data or "error" in data):
                log_recv_queue.put_nowait(data)
            else:
                raise Exception("not a valid request or response")
        except queue.Empty:
            pass
        except Exception as e:
            print("dispatch jsonrpc handler failed")
            send_queue.put_nowait(json.dumps({"error": str(e)}))

@dispatcher.add_method
def example():
  return {"success": 1}

@dispatcher.add_method
def exampleWithParams(sleep):
  time.sleep(sleep)
  return {
    "slept": sleep
  }


def ws_proxy_recv(ws, local_sock, ssock, end_event, global_end_event):
  while not (end_event.is_set() or global_end_event.is_set()):
    try:
      data = ws.recv()
      local_sock.sendall(data)
    except WebSocketTimeoutException:
      pass
    except Exception:
      print("dispatchd.ws_proxy_recv.exception")
      break

  print("dispatch.ws_proxy_recv closing sockets")
  ssock.close()
  local_sock.close()
  print("dispatch.ws_proxy_recv done closing sockets")

  end_event.set()


def ws_proxy_send(ws, local_sock, signal_sock, end_event):
  while not end_event.is_set():
    try:
      r, _, _ = select.select((local_sock, signal_sock), (), ())
      if r:
        if r[0].fileno() == signal_sock.fileno():
          # got end signal from ws_proxy_recv
          end_event.set()
          break
        data = local_sock.recv(4096)
        if not data:
          # local_sock is dead
          end_event.set()
          break

        ws.send(data, ABNF.OPCODE_BINARY)
    except Exception:
      print("dispatchd.ws_proxy_send.exception")
      end_event.set()

  print("dispatch.ws_proxy_send closing sockets")
  signal_sock.close()
  print("dispatch.ws_proxy_send done closing sockets")

def ws_recv(ws, end_event):
    #last_ping = int(sec_since_boot() * 1e9)
    while not end_event.is_set():
        try:
            opcode, data = ws.recv_data(control_frame=True)
            
            if opcode in (ABNF.OPCODE_TEXT, ABNF.OPCODE_BINARY):
                if opcode == ABNF.OPCODE_TEXT:
                    data = data.decode("utf-8")
                recv_queue.put_nowait(data)

            elif opcode == ABNF.OPCODE_PING:
                #last_ping = int(sec_since_boot() * 1e9)
                pass
        except WebSocketTimeoutException:
            print("dispatchd.ws_recv.timeout")
            end_event.set()
        except Exception:
            print("dispatchd.ws_recv.exception")
            end_event.set()


def ws_send(ws, end_event):
  while not end_event.is_set():
    try:
      try:
        data = send_queue.get_nowait()
      except queue.Empty:
        data = log_send_queue.get(timeout=1)
      for i in range(0, len(data), WS_FRAME_SIZE):
        frame = data[i:i+WS_FRAME_SIZE]
        last = i + WS_FRAME_SIZE >= len(data)
        opcode = ABNF.OPCODE_TEXT if i == 0 else ABNF.OPCODE_CONT
        ws.send_frame(ABNF.create_frame(frame, opcode, last))
    except queue.Empty:
      pass
    except Exception:
      print("dispatchd.ws_send.exception")
      end_event.set()

def handle_long_poll(ws):
  end_event = threading.Event()

  threads = [
    threading.Thread(target=ws_recv, args=(ws, end_event), name='ws_recv'),
    threading.Thread(target=ws_send, args=(ws, end_event), name='ws_send'),
  ] + [
    threading.Thread(target=jsonrpc_handler, args=(end_event,), name=f'worker_{x}')
    for x in range(HANDLER_THREADS)
  ]

  for thread in threads:
    thread.start()
  try:
    while not end_event.is_set():
      time.sleep(0.1)
  except (KeyboardInterrupt, SystemExit):
    end_event.set()
    raise
  finally:
    for thread in threads:
      thread.join()


def backoff(retries):
  return random.randrange(0, min(128, int(2 ** retries)))

def main(host, dongleid):
    ws_uri = host + '/' + dongleid
    print(ws_uri)

    while 1:
        try:
            ws = create_connection(ws_uri, cookie="jwt=test", enable_multithread=True, timeout=30.0)
            conn_retries = 0

            handle_long_poll(ws)

        except (KeyboardInterrupt, SystemExit):
            break
        except (ConnectionError, TimeoutError, WebSocketException):
            conn_retries += 1
        except socket.timeout:
            pass
        except Exception as ex:
            print(ex)
            print("dispatch.main.excpetion")

        time.sleep(backoff(conn_retries))


if __name__ == "__main__":
  import argparse

  parser = argparse.ArgumentParser()

  parser.add_argument('--host', default='ws://localhost:4720', dest='host', help='Location of host')
  parser.add_argument('--dongleid', default='e3a435de', dest='dongleid', help='Dongle ID')
  args = parser.parse_args()

  main(args.host, args.dongleid)
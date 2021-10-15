# Dispatch

Dispatch is a JSON RPC proxy server written in go that allows an edge device to serve as an RPC server. To run the examples, first start the proxy server. Then the edge device. Then the client.

![diagram](https://raw.githubusercontent.com/barbinbrad/dispatch/master/assets/diagram.png)

## The Problem

A server is normally used to expose an RPC interface. Mobile edge devices face the following problems when trying to act as an RPC server:

- Connection is intermittent
- IP addresses change
- Web servers cost a lot of CPU

## The Solution

A proxy server can be used to provide a consitent endpoint for the RPC client. The JSON RPC spec allows the proxy server to function as a relay without worrying about the specific methods implemented by the edge device.

All requests to the edge have the same format with varying `method` and `params` values:

```json
{
    "method": "add", 
    "id": 0,
    "jsonrpc": "2.0",
    "params": { 
        "a": 1,
        "b": 2,
    }
}

```

Similarly, all results follow the same format:

```json
{
    "jsonrpc": "2.0",
    "result": { 
        "sum": 3
    },
    "error": {}, 
    "id": 3
}
```



# Dispatch

[!diagram](https://raw.githubusercontent.com/barbinbrad/dispatch/master/assets/diagram.png)

## The Problem

A server is normally used to expose an RPC interface. A mobile edge devices faces the following problems when trying to act as an RPC server:

- Connection is intermittent
- IP addresses can change
- Limited computational resources are spent on edge computing

## The Solution

A proxy server can be used to provide a consitent endpoint for the RPC client. The JSON RPC spec allows the proxy server to function as a relay without worrying about the specific methods implemented by the edge device.

All requests to the edge have the same format with varying `method` and `params` values:

```json
{
    "method": "add", // the name of the function
    "id": 0,
    "jsonrpc": "2.0",
    "params": { // the optional arguments of the function
        "a": 1,
        "b": 2,
    }
}

```

Similarly, all results follow the same format:

```json
{
    "jsonrpc": "2.0",
    "result": { // returned by function
        "sum": 3
    },
    "error": {}, // if needed
    "id": 3
}
```



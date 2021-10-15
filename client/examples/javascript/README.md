# Deno Client

This is a console app that creates many concurrent JSON RPC requests with race conditions. The race conditions are created by generating a random `sleep` time that's passed to the device's `exampleWithParameters` method.

Here's the basics of any client JSON RPC Request:

```ts
type JsonRPCRequest = {
    method: string, // "exampleWithParameters"
    id: number, // 0
    jsonrpc: string, // "2.0"
    params?: unknown // { "sleep": 2 }
};

const response = await fetch(url, {
                method: 'POST',
                headers: { 'Content-type': 'application/json; charset=UTF-8' },
                body: JSON.stringify(jsonRPCRequest),
});

```

## Running the tests

To run the tests, you can run:

```sh
deno run --allow-net --allow-read  stress-test.ts
```

- `-host` is where the proxy server is located
- `-clients` is the number of concurrent clients to be used
- `-requests` is how many requests each client should make
- `-dongleid` is the `dongle_id` of the device to be proxied

All of the flags are optional and will default to the values provided above.


## Installing Deno

You can use any kind of client to make these HTTP requests, but [Deno](https://deno.land) seemed like something fun to try. If you don't have Deno installed you can run:

Shell (Mac, Linux):

```sh
curl -fsSL https://deno.land/x/install/install.sh | sh
```

PowerShell (Windows):

```powershell
iwr https://deno.land/x/install/install.ps1 -useb | iex
```

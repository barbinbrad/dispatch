import { parse } from "https://deno.land/std@0.83.0/flags/mod.ts";

const args = parse(Deno.args);

const clients:number = args.clients || 3;
const requests:number = args.requests || 5;
const host:string = args.host || "http://localhost:4720";
const dongleId:string = args.dongleid || "e3a435de";

const workers:Array<Worker> = [];

// initialize the concurrent clients
for(let client = 0; client<clients; client++){
    workers.push(new Worker(new URL("./lib/worker.ts", import.meta.url).href, { type: "module" }));
}

// start the concurrent clients
for(let worker = 0; worker<workers.length; worker++){
    workers[worker].postMessage({
        clientId: worker + 1,
        requests: requests,
        host: host,
        dongleId: dongleId
    });
}

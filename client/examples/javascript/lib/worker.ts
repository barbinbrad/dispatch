// deno-lint-ignore no-explicit-any
self.onmessage = async (message:any) => {

    const { clientId, requests, host, dongleId } = message.data;
    const url = `${host}/${dongleId}`;

    try{
        for(let request = 0; request<requests; request++){
            const jsonRPCRequest = generateRequest();

            performance.mark("start");

            const response = await fetch(url, {
                method: 'POST',
                headers: { 'Content-type': 'application/json; charset=UTF-8' },
                body: JSON.stringify(jsonRPCRequest),
            });
        
            performance.mark("end");
            const data = await response.text()
            //const data = await response.json();
            const duration = performance.measure("request", "start", "end").duration;
            console.log(clientId, request, duration, data);
        }
    
    } catch(ex){
        console.error(ex);
        self.close();
    }
    
    self.close();
};

function generateRequest() :JsonRPCRequest{
    const request: JsonRPCRequest = {
        method: "example",
        id: 0,
        jsonrpc: "2.0",
        params: {
            "sleep": randomBetweenRange(0.1, 2.5)
        }
    };
    return request
}

function randomBetweenRange(min:number, max:number) :number{
    return Math.random() * (max - min) + min;
}

type JsonRPCRequest = {
    method: string,
    id: number,
    jsonrpc: string,
    params?: unknown
};
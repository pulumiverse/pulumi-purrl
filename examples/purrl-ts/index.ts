import * as purrl from "@pulumiverse/purrl";

let purrlCommand = new purrl.Purrl("purrl", {
    name: "httpbin",
    url: "https://httpbin.org/get",
    method: "GET",
    headers: {
        "test": "test",
    },
    responseCodes: [
        "200"
    ],
    deleteMethod: "DELETE",
    deleteUrl: "https://httpbin.org/delete",
    deleteResponseCodes: [
        "200"
    ],
});

export const response = purrlCommand.responseCode

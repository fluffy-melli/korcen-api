import { request } from "http";
import { resolve } from "path";

const url = "https://korcen.shibadogs.net/api/v1/korcen";
const data = JSON.stringify({
    input: "욕설이 포함될수 있는 메시지",
    replace_front: "감지된 욕설 앞부분에 넣을 메시지 (옵션)",
    replace_end: "감지된 욕설 뒷부분에 넣을 메시지 (옵션)"
    });

const options = {
    method: "POST",
    headers: {
        "Accept": "application/json",
        "Content-Type": "application/json",
        "Content-Length": Buffer.byteLength(data)
    }
};

function callApi() {
    return new Promise<void>((resolve, reject) => {
        const req = request(url, options, (res) => {
            let responseData = "";

            res.on("data", (chunk) => {
                responseData += chunk;
            });

            res.on("end", () => {
                if (res.statusCode === 200) {
                    console.log("Response:", responseData);
                    resolve();
                } else {
                    console.error(`Error: HTTP ${res.statusCode}`);
                    reject(new Error(`HTTP ${res.statusCode}`));
                }
            });
        });

        req.on("error", (error) => {
            console.error("Request Error:", error.message);
            reject(error);
        });

        req.write(data);
        req.end();
    });
}

callApi().catch((error) => {
    console.error("API 호출 중 오류 발생:", error);
});


import Network.HTTP
import Network.URI

main :: IO ()
main = do
    let url = "https://korcen.shibadogs.net/api/v1/korcen"
    let jsonData = "{\"input\": \"욕설이 포함될수 있는 메시지\", \"replace_front\": \"감지된 욕설 앞부분에 넣을 메시지 (옵션)\", \"replace_end\": \"감지된 욕설 뒷부분에 넣을 메시지 (옵션)\"}"
    case parseURI url of
        Nothing -> putStrLn "Invalid URL"
        Just uri -> do
            let req = Request {
                rqURI = uri,
                rqMethod = POST,
                rqHeaders = [mkHeader HdrContentType "application/json"],
                rqBody = jsonData
            }
            response <- simpleHTTP req >>= getResponseBody
            putStrLn ("Response: " ++ response)

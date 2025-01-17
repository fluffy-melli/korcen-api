import Foundation

let url = URL(string: "https://korcen.shibadogs.net/api/v1/korcen")!
var request = URLRequest(url: url)
request.httpMethod = "POST"
request.setValue("application/json", forHTTPHeaderField: "Accept")
request.setValue("application/json", forHTTPHeaderField: "Content-Type")

let jsonData: [String: String] = [
    "input": "욕설이 포함될수 있는 메시지",
    "replace_front": "감지된 욕설 앞부분에 넣을 메시지 (옵션)",
    "replace_end": "감지된 욕설 뒷부분에 넣을 메시지 (옵션)"
]

do {
    request.httpBody = try JSONSerialization.data(withJSONObject: jsonData, options: [])
} catch {
    print("Error encoding JSON:", error)
    exit(1)
}

let task = URLSession.shared.dataTask(with: request) { data, response, error in
    if let error = error {
        print("Error: \(error.localizedDescription)")
        return
    }

    guard let httpResponse = response as? HTTPURLResponse else {
        print("Invalid response")
        return
    }

    if httpResponse.statusCode == 200, let data = data {
        let responseData = String(data: data, encoding: .utf8) ?? "Invalid response encoding"
        print("Response:", responseData)
    } else {
        print("HTTP Error:", httpResponse.statusCode)
    }

    exit(0)
}

task.resume()

RunLoop.main.run()

import requests
import threading

# url = "https://korcen.shibadogs.net/api/v1/korcen"
url = "http://localhost:7777/api/v1/korcen"

data = {
    "input": "string",
    "replace-end": "string",
    "replace-front": "string"
}

def send_request():
    try:
        response = requests.post(url, json=data)
        print(f"Status Code: {response.status_code}, Response: {response.text}")
    except Exception as e:
        print(f"Error: {e}")

threads = []
for _ in range(200):
    thread = threading.Thread(target=send_request)
    threads.append(thread)
    thread.start()

for thread in threads:
    thread.join()

print("completed")
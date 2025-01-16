import requests

url = 'https://korcen.shibadogs.net/api/v1/korcen'

headers = {
    'accept': 'application/json',
    'Content-Type': 'application/json'
}

data = {
    'input'        : '욕설이 포함될수 있는 메시지',
    'replace-front': '감지된 욕설 앞부분에 넣을 메시지 (옵션)',
    'replace-end'  : '감지된 욕설 뒷부분에 넣을 메시지 (옵션)'
}

response = requests.post(url, headers=headers, json=data)

print(response.status_code)
print(response.text)

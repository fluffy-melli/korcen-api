import requests

url = 'https://korcen.shibadogs.net/api/v1/korcen'

headers = {
    'accept': 'application/json',
    'Content-Type': 'application/json'
}

data = {
    'input'        : '욕설이 포함될수 있는 메시지',
    'replace-front': '<',
    'replace-end'  : '>'
}

response = requests.post(url, headers=headers, json=data)

print(response.status_code)
print(response.text)

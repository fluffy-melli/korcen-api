const url = 'https://korcen.shibadogs.net/api/v1/korcen'

const headers = {
    'accept': 'application/json',
    'Content-Type': 'application/json'
}

const data = {
    input: '욕설이 포함될수 있는 메시지',
    'replace-front': '감지된 욕설 앞부분에 넣을 메시지 (옵션)',
    'replace-end': '감지된 욕설 뒷부분에 넣을 메시지 (옵션)'
}

fetch(url, {
    method: 'POST',
    headers: headers,
    body: JSON.stringify(data)
})
.then(response => {
    console.log(response.status)
    return response.json()
})
.then(data => {
    console.log(data)
})
.catch(error => {
    console.error('Error:', error)
})

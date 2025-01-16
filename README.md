<div align="center">
  <h1>Korcen.api</h1>

  [![Go Version](https://github.com/fluffy-melli/korcen-go/blob/main/docs/asset/go_version.svg)](https://go.dev/)
  [![Module Version](https://github.com/fluffy-melli/korcen-go/blob/main/docs/asset/module_version.svg)](https://pkg.go.dev/github.com/fluffy-melli/korcen-go)
</div>

![131_20220604170616](https://user-images.githubusercontent.com/85154556/171998341-9a7439c8-122f-4a9f-beb6-0e0b3aad05ed.png)

# 🛠 제작자

>[Tanat05](https://github.com/Tanat05) / [korcen](https://github.com/Tanat05/korcen)
```
https://github.com/Tanat05/korcen
---------------------------------
이 프로젝트는 원본 `korcen` 프로젝트를 수정하여 배포한 것입니다.
원본 프로젝트는 `https://github.com/Tanat05/korcen`에서 확인할 수 있습니다.

이 프로젝트의 라이센스 또한 Apache-2.0 을 따르고 있습니다
```

>[gyarang](https://github.com/gyarang) / [gohangul](https://github.com/gyarang/gohangul)
```
https://github.com/gyarang/gohangul
-----------------------------------
이 프로젝트는 `gohangul` 을 사용해서 배포한 것입니다.
해당 프로젝트는 `https://github.com/gyarang/gohangul`에서 확인할 수 있습니다.

해당 프로젝트의 라이센스인 MIT License 또한 따르고 있습니다
```

>[Feralthedogg](https://github.com/Feralthedogg)
```
https://github.com/Feralthedogg
-------------------------------
같이 API 프로젝트를 구현한 제작자 입니다
```

![Apache-2.0](https://github.com/fluffy-melli/korcen-go/blob/main/docs/asset/Apache-2.0.png)

`korcen-api`는 `Apache-2.0` & `MIT License` 라이선스를 `모두` 따르고 있습니다.
코드를 사용할 경우 라이선스 내용을 준수해주세요. 

Copyright© All rights reserved.

---

# ❓ 코드 예제

>python
```py
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
```

>javascript
```js
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
```
>Rust
```rs
use reqwest::Client;
use serde::Serialize;

#[derive(Serialize)]
struct KorcenRequest {
    input: String,
    #[serde(rename = "replace-front")]
    replace_front: String,
    #[serde(rename = "replace-end")]
    replace_end: String,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let url = "https://korcen.shibadogs.net/api/v1/korcen";

    let data = KorcenRequest {
        input: "욕설이 포함될수 있는 메시지".to_string(),
        replace_front: "감지된 욕설 앞부분에 넣을 메시지 (옵션)".to_string(),
        replace_end: "감지된 욕설 뒷부분에 넣을 메시지 (옵션)".to_string(),
    };

    let client = Client::new();

    let response = client
        .post(url)
        .header("accept", "application/json")
        .header("Content-Type", "application/json")
        .json(&data)
        .send()
        .await?;

    println!("{}", response.status());

    let text = response.text().await?;
    println!("{}", text);

    Ok(())
}
```

# ⬇️ 설치 방법

>docker
```sh
$ docker build -t korcen-api .
$ docker run -d -p 7777:7777 korcen-api
```

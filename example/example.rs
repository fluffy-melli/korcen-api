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

program KorcenAPI;
uses fphttpclient, fpjson, jsonparser;

var
  Client: TFPHTTPClient;
  Response: String;
  JsonData: String;
begin
  JsonData := '{' +
              '"input": "욕설이 포함될수 있는 메시지",' +
              '"replace_front": "감지된 욕설 앞부분에 넣을 메시지 (옵션)",' +
              '"replace_end": "감지된 욕설 뒷부분에 넣을 메시지 (옵션)"' +
              '}';

  Client := TFPHTTPClient.Create(nil);
  try
    Client.AddHeader('Content-Type', 'application/json');
    Client.AddHeader('Accept', 'application/json');
    Response := Client.SimplePost('https://korcen.shibadogs.net/api/v1/korcen', JsonData);
    Writeln('Response: ', Response);
  except
    on E: Exception do
      Writeln('Error: ', E.Message);
  end;
  Client.Free;
end.

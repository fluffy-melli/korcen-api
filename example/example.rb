require 'net/http'
require 'uri'
require 'json'

uri = URI.parse('https://korcen.shibadogs.net/api/v1/korcen')
header = {'Content-Type' => 'application/json', 'Accept' => 'application/json'}
data = {
  'input' => '욕설이 포함될수 있는 메시지',
  'replace-front' => '감지된 욕설 앞부분에 넣을 메시지 (옵션)',
  'replace-end' => '감지된 욕설 뒷부분에 넣을 메시지 (옵션)'
}

http = Net::HTTP.new(uri.host, uri.port)
request = Net::HTTP::Post.new(uri.path, header)
request.body = data.to_json

response = http.request(request)
puts response.body

using System;
using System.Net.Http;
using System.Text;
using System.Threading.Tasks;

class Program
{
    static async Task Main(string[] args)
    {
        using (HttpClient client = new HttpClient())
        {
            client.DefaultRequestHeaders.Add("Accept", "application/json");
            client.DefaultRequestHeaders.Add("Content-Type", "application/json");

            var data = new
            {
                input = "욕설이 포함될수 있는 메시지",
                replace_front = "감지된 욕설 앞부분에 넣을 메시지 (옵션)",
                replace_end = "감지된 욕설 뒷부분에 넣을 메시지 (옵션)"
            };

            string json = System.Text.Json.JsonSerializer.Serialize(data);
            StringContent content = new StringContent(json, Encoding.UTF8, "application/json");

            HttpResponseMessage response = await client.PostAsync("https://korcen.shibadogs.net/api/v1/korcen", content);

            if (response.IsSuccessStatusCode)
            {
                string responseData = await response.Content.ReadAsStringAsync();
                Console.WriteLine(responseData);
            }
            else
            {
                Console.WriteLine($"Error: {response.StatusCode}");
            }
        }
    }
}

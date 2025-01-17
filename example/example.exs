# mix.exs
#
# defp deps do
#  [
#    {:httpoison, "~> 1.8"},
#    {:jason, "~> 1.4"}
# ]
#end

defmodule KorcenClient do
  @url "https://korcen.shibadogs.net/api/v1/korcen"
  @headers [
    {"Accept", "application/json"},
    {"Content-Type", "application/json"}
  ]

  def call_api do
    body = %{
      "input" => "욕설이 포함될수 있는 메시지",
      "replace_front" => "감지된 욕설 앞부분에 넣을 메시지 (옵션)",
      "repace_end" => "감지된 욕설 뒷부분에 넣을 메시지 (옵션)"
    }
    |> Jason.endcode!()

    case HTTPoison.post(@url, body, @headers) do
      {:ok, %HTTPoison.Response{status_code: 200, body: response_body}} ->
        IO.puts("Response: #{response_body}")

      {:ok, %HTTPoison.Response{status_code: status_code}} ->
        IO.puts("Error: Received HTTP status #{status_code}")

      {:error, HTTPoison.Response{reason: reason}} ->
        IO.puts("HTTP Request failed: #{inspect(reason)}")
    end
  end
end

KorcenClient.call_api()

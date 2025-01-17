import java.net.{HttpURLConnection, URL}
import java.io.{BufferedReader, InputStreamReader, OutputStreamWriter}

object KorcenAPI {
  def main(args: Array[String]): Unit = {
    val url = new URL("https://korcen.shibadogs.net/api/v1/korcen")
    val connection = url.openConnection().asInstanceOf[HttpURLConnection]

    connection.setRequestMethod("POST")
    connection.setRequestProperty("Accept", "application/json")
    connection.setRequestProperty("Content-Type", "application/json")
    connection.setDoOutput(true)

    val jsonData =
      """{
        |  "input": "욕설이 포함될수 있는 메시지",
        |  "replace_front": "감지된 욕설 앞부분에 넣을 메시지 (옵션)",
        |  "replace_end": "감지된 욕설 뒷부분에 넣을 메시지 (옵션)"
        |}""".stripMargin

    val writer = new OutputStreamWriter(connection.getOutputStream, "UTF-8")
    writer.write(jsonData)
    writer.flush()
    writer.close()

    val responseCode = connection.getResponseCode
    println(s"Response Code: $responseCode")

    if (responseCode == HttpURLConnection.HTTP_OK) {
      val reader = new BufferedReader(new InputStreamReader(connection.getInputStream))
      val response = reader.lines().toArray.mkString("\n")
      println(s"Response: $response")
      reader.close()
    } else {
      println(s"Error: HTTP $responseCode")
    }

    connection.disconnect()
  }
}

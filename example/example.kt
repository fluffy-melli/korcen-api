import java.io.*
import java.net.HttpURLConnection
import java.net.URL
import org.json.JSONObject

fun main() {
    val url = URL("https://korcen.shibadogs.net/api/v1/korcen")
    val connection = url.openConnection() as HttpURLConnection
    connection.requestMethod = "POST"
    connection.setRequestProperty("Accept", "application/json")
    connection.setRequestProperty("Content-Type", "application/json")
    connection.doOutput = true

    val data = JSONObject()
    data.put("input", "욕설이 포함될수 있는 메시지")
    data.put("replace-front", "감지된 욕설 앞부분에 넣을 메시지 (옵션)")
    data.put("replace-end", "감지된 욕설 뒷부분에 넣을 메시지 (옵션)")

    val outputStream: OutputStream = connection.outputStream
    val writer = BufferedWriter(OutputStreamWriter(outputStream, "UTF-8"))
    writer.write(data.toString())
    writer.flush()
    writer.close()

    val responseCode = connection.responseCode
    println("Response Code: $responseCode")

    if (responseCode == HttpURLConnection.HTTP_OK) {
        val inputStream = BufferedReader(InputStreamReader(connection.inputStream))
        val response = StringBuffer()
        var inputLine: String?

        while (inputStream.readLine().also { inputLine = it } != null) {
            response.append(inputLine)
        }

        println("Response: ${response.toString()}")
    } else {
        println("Error: $responseCode")
    }

    connection.disconnect()
}

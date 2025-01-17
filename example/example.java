import java.io.*;
import java.net.HttpURLConnection;
import java.net.URL;

public class example {
    public static void main(String[] args) {
        try {
            URL url = new URL("https://korcen.shibadogs.net/api/v1/korcen");
            HttpURLConnection connection = (HttpURLConnection) url.openConnection();
            connection.setRequestMethod("POST");
            connection.setRequestProperty("Accept", "application/json");
            connection.setRequestProperty("Content-Type", "application/json");
            connection.setDoOutput(true);

            String jsonInput = "{"
                    + "\"input\":\"욕설이 포함될수 있는 메시지\","
                    + "\"replace_front\":\"감지된 욕설 앞부분에 넣을 메시지 (옵션)\","
                    + "\"replace_end\":\"감지된 욕설 뒷부분에 넣을 메시지 (옵션)\""
                    + "}";

            try (OutputStream outputStream = connection.getOutputStream();
                 BufferedWriter writer = new BufferedWriter(new OutputStreamWriter(outputStream, "UTF-8"))) {
                writer.write(jsonInput);
                writer.flush();
            }

            int responseCode = connection.getResponseCode();
            System.out.println("Response Code: " + responseCode);

            if (responseCode == HttpURLConnection.HTTP_OK) {
                try (BufferedReader inputStream = new BufferedReader(new InputStreamReader(connection.getInputStream()));
                     StringWriter response = new StringWriter()) {
                    String inputLine;
                    while ((inputLine = inputStream.readLine()) != null) {
                        response.write(inputLine);
                    }
                    System.out.println("Response: " + response.toString());
                }
            } else {
                System.out.println("Error: " + responseCode);
            }

            connection.disconnect();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
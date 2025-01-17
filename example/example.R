url <- "https://korcen.shibadogs.net/api/v1/korcen"
json_data <- '{"input": "욕설이 포함될수 있는 메시지",
               "replace_front": "감지된 욕설 앞부분에 넣을 메시지 (옵션)",
               "replace_end": "감지된 욕설 뒷부분에 넣을 메시지 (옵션)"}'

headers <- c("Accept: application/json", "Content-Type: application/json")
con <- url(url, "rb")
write(json_data, con)
response <- readLines(con, warn=FALSE)
close(con)

cat("Response:", response, "\n")

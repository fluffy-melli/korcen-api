using Sockets

function call_api()
    url = "https://korcen.shibadogs.net/api/v1/korcen"
    json_data = """
    {
        "input": "욕설이 포함될수 있는 메시지",
        "replace_front": "감지된 욕설 앞부분에 넣을 메시지 (옵션)",
        "replace_end": "감지된 욕설 뒷부분에 넣을 메시지 (옵션)"
    }
    """

    host, port = split(url[9:end], "/")[1], 443
    socket = connect(host, port)

    request = """
    POST /api/v1/korcen HTTP/1.1
    Host: korcen.shibadogs.net
    Accept: application/json
    Content-Type: application/json
    Content-Length: $(length(json_data))

    $(json_data)
    """

    write(socket, request)
    response = read(socket, String)
    close(socket)

    println("Response: ", response)
end

call_api()
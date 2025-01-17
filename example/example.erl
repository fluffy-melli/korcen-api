-module(example).
-export([call_api/0]).

call_api() ->
    URL = "https://korcen.shibadogs.net/api/v1/korcen",
    Headers = [
        {"Accept", "application/json"},
        {"Content-Type", "application/json"}
    ],
    Body = "{"
           "\"input\":\"욕설이 포함될수 있는 메시지\","
           "\"replace_front\":\"감지된 욕설 앞부분에 넣을 메시지 (옵션)\","
           "\"replace_end\":\"감지된 욕설 뒷부분에 넣을 메시지 (옵션)\""
           "}",
    
    case httpc:request(post, {URL, Headers, "application/json", Body}, [], []) of
        {ok, {{_, 200, _}, _, ResponseBody}} ->
            io:format("Response: ~s~n", [ResponseBody]);
        {ok, {{_, StatusCode, _}, _, _}} ->
            io:format("Error: HTTP ~p~n", [StatusCode]);
        {error, Reason} ->
            io:format("Request failed: ~p~n", [Reason])
    end.

local json_data = [[
{
    "input": "욕설이 포함될수 있는 메시지",
    "replace_front": "감지된 욕설 앞부분에 넣을 메시지 (옵션)",
    "replace_end": "감지된 욕설 뒷부분에 넣을 메시지 (옵션)"
}
]]

local command = "curl -s -X POST " ..
    "-H \"Accept: application/json\" "
    "-H \"Content-Type: application/json\" "
    "-d '" .. json_data:gsub("'", "\\'") .. "' " ..
    "\"https://korcen.shibadogs.net/api/v1/korcen\""

local handle = io.popen(command)
local response = handle:read("*a")
handle:close()

print("Response: " .. response)
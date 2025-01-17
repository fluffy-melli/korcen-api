#include <iostream>
#include <cstdlib>

int main() {
    std::string json_data =
        "{"
        "\"input\": \"욕설이 포함될수 있는 메시지\","
        "\"replace_front\": \"감지된 욕설 앞부분에 넣을 메시지 (옵션)\","
        "\"replace_end\": \"감지된 욕설 뒷부분에 넣을 메시지 (옵션)\""
        "}";

    std::string command = "echo \"" + json_data + "\" | curl -X POST -H \"Accept: application/json\" -H \"Content-Type: application/json\" -d @- \"https://korcen.shibadogs.net/api/v1/korcen\"";

    system(command.c_str());

    return 0;
}

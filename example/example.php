<?php
$url = "https://korcen.shibadogs.net/api/v1/korcen";
$data = json_encode([
    "input" => "욕설이 포함될수 있는 메시지",
    "replace_front" => "감지된 욕설 앞부분에 넣을 메시지 (옵션)",
    "replace_end" => "감지된 욕설 뒷부분에 넣을 메시지 (옵션)"
]);

$options = [
    "http" => [
        "header"  => "Content-Type: application/json\r\n" .
                     "Accept: application/json\r\n",
        "method"  => "POST",
        "content" => $data
    ]
];

$context = stream_context_create($options);
$response = file_get_contents($url, false, $context);

if ($response === false) {
    echo "Error: Request failed.\n";
} else {
    echo "Response: " . $response . "\n";
}
?>

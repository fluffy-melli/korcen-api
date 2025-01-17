use strict;
use warnings;
use IO::Socket::INET;

my $host = "korcen.shibadogs.net";
my $port = 443;
my $url = "/api/v1/korcen";

my $json_data = qq|{
    "input": "욕설이 포함될수 있는 메시지",
    "replace_front": "감지된 욕설 앞부분에 넣을 메시지 (옵션)",
    "replace_end": "감지된 욕설 뒷부분에 넣을 메시지 (옵션)"
}|;

my $socket = IO::Socket::INET->new(
    PeerAddr => $host,
    PeerPort => $port,
    Proto    => "tcp"
) or die "Could not connect: $!\n";

print $socket "POST $url HTTP/1.1\r\n";
print $socket "Host: $host\r\n";
print $socket "Content-Type: application/json\r\n";
print $socket "Accept: application/json\r\n";
print $socket "Content-Length: " . length($json_data) . "\r\n";
print $socket "\r\n";
print $socket $json_data;

while (<$socket>) {
    print $_;
}

close($socket);

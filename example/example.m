#import <Foundation/Foundation.h>

void sendRequest() {
    NSURL *url = [NSURL URLWithString:@"https://korcen.shibadogs.net/api/v1/korcen"];
    NSMutableURLRequest *request = [NSMutableURLRequest requestWithURL:url];

    [request setHTTPMethod:@"POST"];
    [request setValue:@"application/json" forHTTPHeaderField:@"Accept"];
    [request setValue:@"application/json" forHTTPHeaderField:@"Content-Type"];

    NSDictionary *data = @{
        @"input": @"욕설이 포함될수 있는 메시지",
        @"replace_front": @"감지된 욕설 앞부분에 넣을 메시지 (옵션)",
        @"replace_end": @"감지된 욕설 뒷부분에 넣을 메시지 (옵션)"
    };

    NSError *error;
    NSData *jsonData = [NSJSONSerialization dataWithJSONObject:data options:0 error:&error];

    if (!jsonData) {
        NSLog(@"JSON 변환 실패: %@", error.localizedDescription);
        return;
    }

    [request setHTTPBody:jsonData];

    NSURLSessionDataTask *task = [[NSURLSession sharedSession] dataTaskWithRequest:request completionHandler:^(NSData *data, NSURLResponse *response, NSError *error) {
        if (error) {
            NSLog(@"요청 실패: %@", error.localizedDescription);
            return;
        }

        NSHTTPURLResponse *httpResponse = (NSHTTPURLResponse *)response;
        if (httpResponse.statusCode == 200) {
            NSString *responseString = [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding];
            NSLog(@"응답: %@", responseString);
        } else {
            NSLog(@"에러 발생: HTTP %ld", (long)httpResponse.statusCode);
        }
    }];

    [task resume];
}

int main(int argc, const char * argv[]) {
    @autoreleasepool {
        sendRequest();

        [[NSRunLoop currentRunLoop] run];
    }
    return 0;
}

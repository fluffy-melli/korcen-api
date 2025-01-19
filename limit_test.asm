; limit_test.asm
; NASM syntax x86-64 Windows

; 아직 미완

BITS 64

SECTION .data
    jsonData db '{"input":"string","replace-end":"string","replace-front":"string"}',0
    jsonDataLen equ $ - jsonData - 1

    httpRequestHeader db "POST /api/v1/korcen HTTP/1.1\r\nHost: localhost:7777\r\nContent-Type: application/json\r\nContent-Length: 64\r\nConnection: close\r\n\r\n",0
    httpRequestHeaderLen equ $ - httpRequestHeader -1

    completeMsg db "200 HTTP POST Complete",10,0
    wsaStartMsg db "WSAStartup successful",10,0
    socketCreatedMsg db "Socket created",10,0
    connectedMsg db "Connected to server",10,0
    sendHeaderMsg db "HTTP header sent",10,0
    sendJsonMsg db "JSON data sent",10,0
    recvMsg db "Response received",10,0
    connectFailedMsg db "Connect failed",10,0
    sendFailedMsg db "Send failed",10,0
    recvFailedMsg db "Recv failed",10,0
    cleanupMsg db "WSACleanup called",10,0
    exitErrorMsg db "WSAStartup failed",10,0

sockaddr_in:
    dw 2
    dw 0x1E61
    dd 0x0100007F
    times 8 db 0

SECTION .bss
    wsaData resb 104
    sock resq 1
    recvBuffer resb 4096

SECTION .text
    extern ExitProcess
    extern WSAStartup
    extern WSACleanup
    extern socket
    extern connect
    extern send
    extern recv
    extern closesocket
    extern printf

    global main

main:
    mov ecx, 0x0202
    lea rdx, [rel wsaData]
    sub rsp, 40
    call WSAStartup
    add rsp, 40
    cmp eax, 0
    jne exit_error

    lea rcx, [rel wsaStartMsg]
    sub rsp, 40
    call printf
    add rsp, 40

    mov r12d, 0

send_loop:
    cmp r12d, 200
    jge end_loop

    mov ecx, 2
    mov edx, 1
    mov r8d, 6
    sub rsp, 40
    call socket
    add rsp, 40
    cmp rax, -1
    je loop_continue
    mov [rel sock], rax

    lea rcx, [rel socketCreatedMsg]
    sub rsp, 40
    call printf
    add rsp, 40

    mov rax, [rel sock]
    mov rcx, rax
    lea rdx, [rel sockaddr_in]
    mov r8d, 16
    sub rsp, 40
    call connect
    add rsp, 40
    cmp eax, -1
    je connect_failed

    lea rcx, [rel connectedMsg]
    sub rsp, 40
    call printf
    add rsp, 40

    mov rax, [rel sock]
    mov rcx, rax
    lea rdx, [rel httpRequestHeader]
    mov r8d, httpRequestHeaderLen
    mov r9d, 0
    sub rsp, 40
    call send
    add rsp, 40
    cmp eax, -1
    je send_failed

    lea rcx, [rel sendHeaderMsg]
    sub rsp, 40
    call printf
    add rsp, 40

    mov rax, [rel sock]
    mov rcx, rax
    lea rdx, [rel jsonData]
    mov r8d, jsonDataLen
    mov r9d, 0
    sub rsp, 40
    call send
    add rsp, 40
    cmp eax, -1
    je send_failed

    lea rcx, [rel sendJsonMsg]
    sub rsp, 40
    call printf
    add rsp, 40

    mov rax, [rel sock]
    mov rcx, rax
    lea rdx, [rel recvBuffer]
    mov r8d, 4096
    mov r9d, 0
    sub rsp, 40
    call recv
    add rsp, 40
    cmp eax, -1
    je recv_failed

    lea rcx, [rel recvMsg]
    sub rsp, 40
    call printf
    add rsp, 40

    mov rax, [rel sock]
    mov rcx, rax
    sub rsp, 40
    call closesocket
    add rsp, 40

    inc r12d

    jmp send_loop

connect_failed:
    lea rcx, [rel connectFailedMsg]
    sub rsp, 40
    call printf
    add rsp, 40

    mov rax, [rel sock]
    mov rcx, rax
    sub rsp, 40
    call closesocket
    add rsp, 40

    inc r12d

    jmp send_loop

send_failed:
    lea rcx, [rel sendFailedMsg]
    sub rsp, 40
    call printf
    add rsp, 40

    mov rax, [rel sock]
    mov rcx, rax
    sub rsp, 40
    call closesocket
    add rsp, 40

    inc r12d

    jmp send_loop

recv_failed:
    lea rcx, [rel recvFailedMsg]
    sub rsp, 40
    call printf
    add rsp, 40

    mov rax, [rel sock]
    mov rcx, rax
    sub rsp, 40
    call closesocket
    add rsp, 40

    inc r12d

    jmp send_loop

loop_continue:
    inc r12d
    jmp send_loop

end_loop:
    sub rsp, 40
    call WSACleanup
    add rsp, 40

    lea rcx, [rel cleanupMsg]
    sub rsp, 40
    call printf
    add rsp, 40

    lea rcx, [rel completeMsg]
    sub rsp, 40
    call printf
    add rsp, 40

    xor ecx, ecx
    sub rsp, 40
    call ExitProcess
    add rsp, 40

exit_error:
    lea rcx, [rel exitErrorMsg]
    sub rsp, 40
    call printf
    add rsp, 40

    mov ecx, 1
    sub rsp, 40
    call ExitProcess
    add rsp, 40

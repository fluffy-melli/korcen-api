; limit_test.asm
; NASM syntax x86-64 Windows

; Made By Tjdmin1

section .data
    url db "127.0.0.1", 0
    port dw 7777
    http_request db "POST /api/v1/korcen HTTP/1.1", 0x0d, 0x0a
                 db "Host: localhost", 0x0d, 0x0a
                 db "Content-Type: application/json", 0x0d, 0x0a
                 db "Content-Length: 58", 0x0d, 0x0a, 0x0d, 0x0a
                 db '{"input":"string","replace-end":"string","replace-front":"string"}', 0

section .bss
    response_buffer resb 1024
    socket_fd resq 1
    loop_count resq 1       ; 반복 카운트 저장소

section .text
    global _start

_start:
    mov qword [loop_count], 200  ; 초기 반복 횟수 설정

.loop:
    call send_request            ; send_request 함수 호출

    ; 반복 횟수 감소
    mov rax, [loop_count]        ; loop_count 값을 rax로 로드
    dec rax                      ; rax 감소
    mov [loop_count], rax        ; 감소된 값을 loop_count에 저장

    cmp rax, 0                   ; rax가 0인지 확인
    jne .loop                    ; 0이 아니면 .loop로 이동

    ; 프로그램 종료
    mov rax, 60                  ; exit system call
    xor rdi, rdi                 ; exit code 0
    syscall

send_request:
    mov rax, 41             ; socket system call
    xor rdi, rdi
    xor rsi, rsi
    xor rdx, rdx
    syscall
    mov [socket_fd], rax    ; socket 파일 디스크립터 저장

    mov rdi, [socket_fd]    ; connect system call
    lea rsi, [rel url]
    mov rdx, 16
    mov rax, 42
    syscall

    mov rdi, [socket_fd]    ; send system call
    lea rsi, [rel http_request]
    mov rdx, 128
    mov rax, 44
    syscall

    mov rdi, [socket_fd]    ; recv system call
    lea rsi, [response_buffer]
    mov rdx, 1024
    mov rax, 45
    syscall

    mov rdi, 1              ; write system call
    lea rsi, [response_buffer]
    mov rdx, rax
    mov rax, 1
    syscall

    mov rdi, [socket_fd]    ; close system call
    mov rax, 3
    syscall
    ret

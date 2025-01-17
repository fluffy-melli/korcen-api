program korcen_api
    use iso_fortran_env, only: output_unit
    implicit none
    character(len=1024) :: json_data, response
    integer :: status

    json_data = '{' // &
        '"input": "욕설이 포함될수 있는 메시지",' // &
        '"replace_front": "감지된 욕설 앞부분에 넣을 메시지 (옵션)",' // &
        '"replace_end": "감지된 욕설 뒷부분에 넣을 메시지 (옵션)"' // &
        '}'

    call execute_command_line('curl -s -X POST ' // &
        '-H "Accept: application/json" ' // &
        '-H "Content-Type: application/json" ' // &
        '-d ''' // trim(json_data) // ''' ' // &
        '"https://korcen.shibadogs.net/api/v1/korcen" > response.json', status)

    open(unit=10, file="response.json", status="old", action="read", iostat=status)
    if (status == 0) then
        read(10, '(A)') response
        print *, "Response:", trim(response)
        close(10)
    else
        print *, "Error: Failed to read response."
    end if
end program korcen_api

@256
D=A
@SP
M=D
@RETURN_ADDRESS_1
D=A
@SP
A=M
M=D
@SP
M=M+1
@LCL
D=M
@SP
A=M
M=D
@SP
M=M+1
@ARG
D=M
@SP
A=M
M=D
@SP
M=M+1
@THIS
D=M
@SP
A=M
M=D
@SP
M=M+1
@THAT
D=M
@SP
A=M
M=D
@SP
M=M+1
@SP
D=M
@5
D=D-A
@0
D=D-A
@ARG
M=D
@SP
D=M
@LCL
M=D
@Sys.init
0;JMP
(RETURN_ADDRESS_1)
// function Sys.init 0
(Sys.init)
@0
D=A
@tmp
M=D
(Sys.init_init)
@tmp
D=M
@Sys.init_end_init
D;JEQ
@tmp
M=M-1
@SP
A=M
M=0
@SP
M=M+1
@Sys.init_init
0;JMP
(Sys.init_end_init)
// push constant 4000
@4000
D=A
@SP
A=M
M=D
@SP
M=M+1
// pop pointer 0
@SP
M=M-1
A=M
D=M
@THIS
M=D
// push constant 5000
@5000
D=A
@SP
A=M
M=D
@SP
M=M+1
// pop pointer 1
@SP
M=M-1
A=M
D=M
@THAT
M=D
// call Sys.main 0
@RETURN_ADDRESS_2
D=A
@SP
A=M
M=D
@SP
M=M+1
@LCL
D=M
@SP
A=M
M=D
@SP
M=M+1
@ARG
D=M
@SP
A=M
M=D
@SP
M=M+1
@THIS
D=M
@SP
A=M
M=D
@SP
M=M+1
@THAT
D=M
@SP
A=M
M=D
@SP
M=M+1
@SP
D=M
@5
D=D-A
@0
D=D-A
@ARG
M=D
@SP
D=M
@LCL
M=D
@Sys.main
0;JMP
(RETURN_ADDRESS_2)
// pop temp 1
@SP
M=M-1
A=M
D=M
@6
M=D
// label LOOP
(LOOP)
// goto LOOP
@LOOP
0;JMP
// function Sys.main 5
(Sys.main)
@5
D=A
@tmp
M=D
(Sys.main_init)
@tmp
D=M
@Sys.main_end_init
D;JEQ
@tmp
M=M-1
@SP
A=M
M=0
@SP
M=M+1
@Sys.main_init
0;JMP
(Sys.main_end_init)
// push constant 4001
@4001
D=A
@SP
A=M
M=D
@SP
M=M+1
// pop pointer 0
@SP
M=M-1
A=M
D=M
@THIS
M=D
// push constant 5001
@5001
D=A
@SP
A=M
M=D
@SP
M=M+1
// pop pointer 1
@SP
M=M-1
A=M
D=M
@THAT
M=D
// push constant 200
@200
D=A
@SP
A=M
M=D
@SP
M=M+1
// pop local 1
@SP
M=M-1
@LCL
D=M
@1
D=D+A
@Addr
M=D
@SP
A=M
D=M
@Addr
A=M
M=D
// push constant 40
@40
D=A
@SP
A=M
M=D
@SP
M=M+1
// pop local 2
@SP
M=M-1
@LCL
D=M
@2
D=D+A
@Addr
M=D
@SP
A=M
D=M
@Addr
A=M
M=D
// push constant 6
@6
D=A
@SP
A=M
M=D
@SP
M=M+1
// pop local 3
@SP
M=M-1
@LCL
D=M
@3
D=D+A
@Addr
M=D
@SP
A=M
D=M
@Addr
A=M
M=D
// push constant 123
@123
D=A
@SP
A=M
M=D
@SP
M=M+1
// call Sys.add12 1
@RETURN_ADDRESS_3
D=A
@SP
A=M
M=D
@SP
M=M+1
@LCL
D=M
@SP
A=M
M=D
@SP
M=M+1
@ARG
D=M
@SP
A=M
M=D
@SP
M=M+1
@THIS
D=M
@SP
A=M
M=D
@SP
M=M+1
@THAT
D=M
@SP
A=M
M=D
@SP
M=M+1
@SP
D=M
@5
D=D-A
@1
D=D-A
@ARG
M=D
@SP
D=M
@LCL
M=D
@Sys.add12
0;JMP
(RETURN_ADDRESS_3)
// pop temp 0
@SP
M=M-1
A=M
D=M
@5
M=D
// push local 0
@LCL
D=M
@0
A=D+A
D=M
@SP
A=M
M=D
@SP
M=M+1
// push local 1
@LCL
D=M
@1
A=D+A
D=M
@SP
A=M
M=D
@SP
M=M+1
// push local 2
@LCL
D=M
@2
A=D+A
D=M
@SP
A=M
M=D
@SP
M=M+1
// push local 3
@LCL
D=M
@3
A=D+A
D=M
@SP
A=M
M=D
@SP
M=M+1
// push local 4
@LCL
D=M
@4
A=D+A
D=M
@SP
A=M
M=D
@SP
M=M+1
// add
@SP
M=M-1
A=M
D=M
A=A-1
M=M+D
// add
@SP
M=M-1
A=M
D=M
A=A-1
M=M+D
// add
@SP
M=M-1
A=M
D=M
A=A-1
M=M+D
// add
@SP
M=M-1
A=M
D=M
A=A-1
M=M+D
// return
@LCL
D=M
@endFrame
M=D
@5
A=D-A
D=M
@retAddr
M=D
@SP
M=M-1
A=M
D=M
@ARG
A=M
M=D
@ARG
D=M
D=D+1
@SP
M=D
@endFrame
D=M
A=D-1
D=M
@THAT
M=D
@endFrame
D=M
@2
A=D-A
D=M
@THIS
M=D
@endFrame
D=M
@3
A=D-A
D=M
@ARG
M=D
@endFrame
D=M
@4
A=D-A
D=M
@LCL
M=D
@retAddr
A=M;JMP
// function Sys.add12 0
(Sys.add12)
@0
D=A
@tmp
M=D
(Sys.add12_init)
@tmp
D=M
@Sys.add12_end_init
D;JEQ
@tmp
M=M-1
@SP
A=M
M=0
@SP
M=M+1
@Sys.add12_init
0;JMP
(Sys.add12_end_init)
// push constant 4002
@4002
D=A
@SP
A=M
M=D
@SP
M=M+1
// pop pointer 0
@SP
M=M-1
A=M
D=M
@THIS
M=D
// push constant 5002
@5002
D=A
@SP
A=M
M=D
@SP
M=M+1
// pop pointer 1
@SP
M=M-1
A=M
D=M
@THAT
M=D
// push argument 0
@ARG
D=M
@0
A=D+A
D=M
@SP
A=M
M=D
@SP
M=M+1
// push constant 12
@12
D=A
@SP
A=M
M=D
@SP
M=M+1
// add
@SP
M=M-1
A=M
D=M
A=A-1
M=M+D
// return
@LCL
D=M
@endFrame
M=D
@5
A=D-A
D=M
@retAddr
M=D
@SP
M=M-1
A=M
D=M
@ARG
A=M
M=D
@ARG
D=M
D=D+1
@SP
M=D
@endFrame
D=M
A=D-1
D=M
@THAT
M=D
@endFrame
D=M
@2
A=D-A
D=M
@THIS
M=D
@endFrame
D=M
@3
A=D-A
D=M
@ARG
M=D
@endFrame
D=M
@4
A=D-A
D=M
@LCL
M=D
@retAddr
A=M;JMP
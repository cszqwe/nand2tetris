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
// function Main.fibonacci 0
(Main.fibonacci)
@0
D=A
@tmp
M=D
(Main.fibonacci_init)
@tmp
D=M
@Main.fibonacci_end_init
D;JEQ
@tmp
M=M-1
@SP
A=M
M=0
@SP
M=M+1
@Main.fibonacci_init
0;JMP
(Main.fibonacci_end_init)
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
// push constant 2
@2
D=A
@SP
A=M
M=D
@SP
M=M+1
// lt
@tmp
M=-1
@SP
M=M-1
A=M
D=M
A=A-1
D=M-D
@CONTINUE2
D;JLT
@tmp
M=0
(CONTINUE2)
@tmp
D=M
@SP
A=M-1
M=D
// if-goto IF_TRUE
@SP
M=M-1
A=M
D=M
@IF_TRUE
D;JNE
// goto IF_FALSE
@IF_FALSE
0;JMP
// label IF_TRUE
(IF_TRUE)
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
// label IF_FALSE
(IF_FALSE)
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
// push constant 2
@2
D=A
@SP
A=M
M=D
@SP
M=M+1
// sub
@SP
M=M-1
A=M
D=M
A=A-1
M=M-D
// call Main.fibonacci 1
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
@Main.fibonacci
0;JMP
(RETURN_ADDRESS_3)
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
// push constant 1
@1
D=A
@SP
A=M
M=D
@SP
M=M+1
// sub
@SP
M=M-1
A=M
D=M
A=A-1
M=M-D
// call Main.fibonacci 1
@RETURN_ADDRESS_4
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
@Main.fibonacci
0;JMP
(RETURN_ADDRESS_4)
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
// push constant 4
@4
D=A
@SP
A=M
M=D
@SP
M=M+1
// call Main.fibonacci 1
@RETURN_ADDRESS_5
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
@Main.fibonacci
0;JMP
(RETURN_ADDRESS_5)
// label WHILE
(WHILE)
// goto WHILE
@WHILE
0;JMP
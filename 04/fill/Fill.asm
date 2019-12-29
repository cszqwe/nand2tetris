// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// Put your code here.

// main function
(LOOP)
    @KBD
    D=M
    // if KBD != 0 
    //  goto Black
    // goto White
    @BLACK
    D;JNE
    @WHITE
    0;JMP

(BLACK)
    @i
    M=0 // i = 0
    (BLACKLOOP)
        @i
        D=M    
        @8192
        D=D-A
        @LOOP
        D;JGE // if i >= 8192, come back to the loop
        @i
        D=M
        @SCREEN
        A=A+D // screen[i] = BLACK
        M=-1
        @i
        M=M+1 // i = i + 1
        @BLACKLOOP
        0;JMP //go back to BLACKLOOP
  
(WHITE)
    @i
    M=0 // i = 0
    (WHITELOOP)
        @i
        D=M    
        @8192
        D=D-A
        @LOOP
        D;JGE // if i >= 8192, come back to the loop
        @i
        D=M
        @SCREEN
        A=A+D // screen[i] = WHITE
        M=0
        @i
        M=M+1 // i = i + 1
        @WHITELOOP
        0;JMP //go back to WHITELOOP
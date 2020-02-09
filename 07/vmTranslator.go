package main

import (
	"bufio"
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fileName := os.Args[1]
	f, err := os.Open(fileName)
	check(err)
	in := bufio.NewReader(f)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	lines := []string{}
	for true {
		line, err := in.ReadBytes('\n')
		if err != nil {
			break
		}
		lineStr := trim(string(line))
		if len(lineStr) == 0 {
			continue
		}
		if len(lineStr) >= 2 && lineStr[:2] == "//" {
			continue
		}
		lines = append(lines, lineStr)
	}
	commands := parse(lines)
	for _, command := range commands {
		tmp := command.write()
		//fmt.Fprintln(out, command)
		for _, s := range tmp {
			fmt.Fprintln(out, s)
		}
	}
	return
}

func process(lines []string) []string {

	return nil
}

func trim(s string) string {
	ans := []byte{}
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' || s[i] == 13 {
			continue
		}
		if s[i] == '/' {
			break
		}
		ans = append(ans, s[i])
	}
	return string(ans)
}

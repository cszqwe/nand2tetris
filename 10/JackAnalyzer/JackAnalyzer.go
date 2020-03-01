package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	name := os.Args[1]
	fi, err := os.Stat(name)
	if err != nil {
		panic(err)
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		processDir(name)
	case mode.IsRegular():
		processFile(name)
	}

	return
}

func processDir(dirName string) {
	if dirName[len(dirName)-1] != '/' {
		dirName += "/"
	}
	filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".jack") {
			f, err := os.Open(path)
			check(err)
			in := bufio.NewReader(f)
			outFileName := path[:len(path)-5] + ".xml"
			outFile, err := os.Create(outFileName)
			check(err)
			out := bufio.NewWriter(outFile)
			outTokenFileName := path[:len(path)-5] + "T.xml"
			outTokenFile, err := os.Create(outTokenFileName)
			check(err)
			outToken := bufio.NewWriter(outTokenFile)
			processSingleFile(in, outToken, out)
		}
		return nil
	})
}

func processFile(fileName string) {
	f, err := os.Open(fileName)
	check(err)
	outputFileName := strings.Replace(fileName, ".jack", ".xml", 1)
	outputTokenFileName := strings.Replace(fileName, ".jack", "T.xml", 1)
	of, err := os.Create(outputFileName)
	check(err)
	otf, err := os.Create(outputTokenFileName)
	in := bufio.NewReader(f)
	out := bufio.NewWriter(of)
	outToken := bufio.NewWriter(otf)
	processSingleFile(in, outToken, out)
}

func processSingleFile(in *bufio.Reader, outToken *bufio.Writer, out *bufio.Writer) {
	defer out.Flush()
	defer outToken.Flush()
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
	tokens := []token{}
	for _, line := range lines {
		tmp := getTokens(line)
		tokens = append(tokens, tmp...)
	}
	fmt.Fprintln(outToken, "<tokens>")
	for _, token := range tokens {
		fmt.Fprintln(outToken, "	"+token.print())
	}
	fmt.Fprintln(outToken, "</tokens>")
	c, _ := compileClass(tokens)
	for _, str := range c.print() {
		fmt.Fprintln(out, str)
	}
}

func trim(s string) string {
	ans := []byte{}
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' || s[i] == 13 {
			continue
		}
		if s[i] == '/' && i < len(s)-1 && (s[i+1] == '/' || s[i+1] == '*') {
			break
		}
		ans = append(ans, s[i])
	}
	for len(ans) > 0 && (ans[len(ans)-1] == ' ' || ans[len(ans)-1] == '\t') {
		ans = ans[:len(ans)-1]
	}
	return string(ans)
}

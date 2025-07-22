package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solveXML(s string) string {
	indent := 0
	var out strings.Builder
	for i := 0; i < len(s); {
		if s[i] == '<' {
			j := i + 1
			for j < len(s) && s[j] != '>' {
				j++
			}
			if j >= len(s) {
				break
			}
			token := s[i : j+1]
			if len(token) >= 2 && token[1] == '/' {
				indent--
				if indent < 0 {
					indent = 0
				}
			}
			out.WriteString(strings.Repeat(" ", indent*2))
			out.WriteString(token)
			out.WriteByte('\n')
			if len(token) >= 2 && token[1] != '/' {
				indent++
			}
			i = j + 1
		} else {
			i++
		}
	}
	return strings.TrimSpace(out.String())
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Scan()
	t, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		s := strings.TrimSpace(scan.Text())
		expected := solveXML(s)
		out, err := run(bin, s+"\n")
		if err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Printf("case %d failed: expected:\n%s\n----\ngot:\n%s\n", i+1, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

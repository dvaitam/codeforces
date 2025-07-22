package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runRef(input string) (string, error) {
	cmd := exec.Command("go", "run", "303E.go")
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runBin(bin string, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not open testcasesE.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	if !sc.Scan() {
		fmt.Fprintln(os.Stderr, "empty test file")
		os.Exit(1)
	}
	var t int
	fmt.Sscan(sc.Text(), &t)
	for i := 0; i < t; i++ {
		if !sc.Scan() {
			fmt.Fprintf(os.Stderr, "case %d missing n line\n", i+1)
			os.Exit(1)
		}
		nLine := sc.Text()
		var n int
		fmt.Sscan(nLine, &n)
		var block strings.Builder
		block.WriteString(nLine)
		block.WriteByte('\n')
		for j := 0; j < n; j++ {
			if !sc.Scan() {
				fmt.Fprintf(os.Stderr, "case %d incomplete interval lines\n", i+1)
				os.Exit(1)
			}
			line := sc.Text()
			block.WriteString(line)
			block.WriteByte('\n')
		}
		input := block.String()
		exp, err := runRef(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected:%s\ngot:%s\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

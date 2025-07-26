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

func runRef(input string) (string, error) {
	cmd := exec.Command("go", "run", "919F.go")
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runBin(bin, input string) (string, error) {
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
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) == 3 {
		bin = os.Args[2]
	}
	data, err := os.ReadFile("testcasesF.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not read testcasesF.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Fprintln(os.Stderr, "bad file")
			os.Exit(1)
		}
		fval := scan.Text()
		alice := make([]string, 8)
		for i := 0; i < 8; i++ {
			if !scan.Scan() {
				fmt.Fprintln(os.Stderr, "bad file")
				os.Exit(1)
			}
			alice[i] = scan.Text()
		}
		bob := make([]string, 8)
		for i := 0; i < 8; i++ {
			if !scan.Scan() {
				fmt.Fprintln(os.Stderr, "bad file")
				os.Exit(1)
			}
			bob[i] = scan.Text()
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fval + "\n")
		sb.WriteString(strings.Join(alice, " ") + "\n")
		sb.WriteString(strings.Join(bob, " ") + "\n")
		input := sb.String()
		exp, err := runRef(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on case %d: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
		got, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", caseIdx+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

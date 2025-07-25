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

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func expected(n int) string {
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		if i%2 == 1 {
			sb.WriteString("I hate")
		} else {
			sb.WriteString("I love")
		}
		if i == n {
			sb.WriteString(" it")
		} else {
			sb.WriteString(" that ")
		}
	}
	return sb.String()
}

func loadTests() ([]int, error) {
	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		return nil, err
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("invalid test file")
	}
	t, _ := strconv.Atoi(scan.Text())
	res := make([]int, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			return nil, fmt.Errorf("not enough values in test file")
		}
		res[i], _ = strconv.Atoi(scan.Text())
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	tests, err := loadTests()
	if err != nil {
		fmt.Println("could not read testcasesA.txt:", err)
		os.Exit(1)
	}
	for idx, n := range tests {
		input := fmt.Sprintf("%d\n", n)
		exp := expected(n) + "\n"
		out, err := runBinary(exe, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", idx+1, err, out)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Printf("case %d failed:\ninput:%sexpected:%s\ngot:%s", idx+1, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

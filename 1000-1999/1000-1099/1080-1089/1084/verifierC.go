package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const MOD int64 = 1000000007

func solveCase(s string) int64 {
	ans := int64(1)
	cnt := int64(0)
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case 'a':
			cnt++
		case 'b':
			ans = ans * (cnt + 1) % MOD
			cnt = 0
		}
	}
	ans = ans * (cnt + 1) % MOD
	res := ans - 1
	if res < 0 {
		res += MOD
	}
	return res % MOD
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	f, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open testcasesC.txt: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	reader := bufio.NewReader(f)

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read test count: %v\n", err)
		os.Exit(1)
	}

	for caseNum := 1; caseNum <= t; caseNum++ {
		var s string
		if _, err := fmt.Fscan(reader, &s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: read string: %v\n", caseNum, err)
			os.Exit(1)
		}
		input := s + "\n"
		want := fmt.Sprintf("%d", solveCase(s))
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", caseNum, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}

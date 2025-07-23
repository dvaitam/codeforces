package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func run(bin, input string) (string, error) {
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

type testCaseF struct {
	input string
}

func parseTestcases(path string) ([]testCaseF, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	in := bufio.NewReader(f)
	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return nil, err
	}
	cases := make([]testCaseF, T)
	for i := 0; i < T; i++ {
		var n, p int
		fmt.Fscan(in, &n, &p)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, p))
		for j := 0; j < n; j++ {
			var c, h, t int
			fmt.Fscan(in, &c, &h, &t)
			sb.WriteString(fmt.Sprintf("%d %d %d\n", c, h, t))
		}
		var q int
		fmt.Fscan(in, &q)
		sb.WriteString(fmt.Sprintf("%d\n", q))
		for j := 0; j < q; j++ {
			var a, b int
			fmt.Fscan(in, &a, &b)
			sb.WriteString(fmt.Sprintf("%d %d\n", a, b))
		}
		cases[i] = testCaseF{input: sb.String()}
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases("testcasesF.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		expected, err := run("500F.go", tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ref failed on case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, strings.TrimSpace(expected), got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

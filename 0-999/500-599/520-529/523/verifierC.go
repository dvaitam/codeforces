package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runCandidate(bin string, input string) (string, error) {
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

func expectedOutput(s, t string) string {
	m := len(s)
	n := len(t)
	L := make([]int, m)
	pj := 0
	for i := 0; i < n && pj < m; i++ {
		if t[i] == s[pj] {
			L[pj] = i
			pj++
		}
	}
	if pj < m {
		return "0"
	}
	R := make([]int, m)
	pj = m - 1
	for i := n - 1; i >= 0 && pj >= 0; i-- {
		if t[i] == s[pj] {
			R[pj] = i
			pj--
		}
	}
	if pj >= 0 {
		return "0"
	}
	ans := R[0] - L[m-1]
	if ans < 0 {
		ans = 0
	}
	return fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var tcases int
	if _, err := fmt.Fscan(reader, &tcases); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read test count: %v\n", err)
		os.Exit(1)
	}
	for caseNum := 1; caseNum <= tcases; caseNum++ {
		var s, t string
		fmt.Fscan(reader, &s)
		fmt.Fscan(reader, &t)
		input := s + "\n" + t + "\n"
		want := expectedOutput(s, t)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", caseNum, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", tcases)
}

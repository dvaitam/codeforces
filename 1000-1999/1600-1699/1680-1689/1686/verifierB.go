package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func expected(perm []int) int {
	res := 0
	for i := 0; i < len(perm)-1; i++ {
		if perm[i] > perm[i+1] {
			res++
			i++
		}
	}
	return res
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	file, err := os.Open(filepath.Join(dir, "testcasesB.txt"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open testcases:", err)
		os.Exit(1)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		fmt.Fprintln(os.Stderr, "failed to read test count:", err)
		os.Exit(1)
	}
	for caseNum := 1; caseNum <= t; caseNum++ {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: failed to read n: %v\n", caseNum, err)
			os.Exit(1)
		}
		perm := make([]int, n)
		for i := 0; i < n; i++ {
			if _, err := fmt.Fscan(reader, &perm[i]); err != nil {
				fmt.Fprintf(os.Stderr, "case %d: failed to read value: %v\n", caseNum, err)
				os.Exit(1)
			}
		}
		input := fmt.Sprintf("1\n%d\n", n)
		for i, v := range perm {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		want := fmt.Sprintf("%d", expected(perm))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", caseNum, err)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\n got: %s\n", caseNum, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}

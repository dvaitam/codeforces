package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func solve(n, m, k int) string {
	if n*m-1 == k {
		return "YES"
	}
	return "NO"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	in := bufio.NewReader(file)
	var t int
	fmt.Fscan(in, &t)
	for i := 1; i <= t; i++ {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)
		var input strings.Builder
		input.WriteString("1\n")
		fmt.Fprintf(&input, "%d %d %d\n", n, m, k)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n%s", i, err, errBuf.String())
			os.Exit(1)
		}
		got := strings.ToUpper(strings.TrimSpace(out.String()))
		exp := solve(n, m, k)
		if got != exp {
			fmt.Printf("case %d failed: expected %s got %s\n", i, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}

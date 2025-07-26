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

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n int) string {
	res := 2*n*(n-1) + 1
	return fmt.Sprintf("%d", res)
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open testcases:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "empty test file")
		os.Exit(1)
	}
	T, _ := strconv.Atoi(scan.Text())
	for tc := 0; tc < T; tc++ {
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "bad test %d\n", tc+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		input := fmt.Sprintf("%d\n", n)
		want := expected(n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", tc+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput: %sexpected: %s\ngot: %s\n", tc+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", T)
}

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

func solveCase(n, k, t int64) int64 {
	var ans int64
	switch {
	case t <= k:
		ans = t
	case t <= n:
		ans = k
	default:
		ans = n + k - t
	}
	if ans < 0 {
		ans = 0
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Println("could not read testcasesA.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	tCases, _ := strconv.Atoi(scan.Text())
	for i := 0; i < tCases; i++ {
		var vals [3]int64
		for j := 0; j < 3; j++ {
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			v, _ := strconv.ParseInt(scan.Text(), 10, 64)
			vals[j] = v
		}
		n, k, t := vals[0], vals[1], vals[2]
		input := fmt.Sprintf("%d %d %d\n", n, k, t)
		expected := solveCase(n, k, t)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		gotVal, err := strconv.ParseInt(strings.TrimSpace(got), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: output not an integer: %s\n", i+1, got)
			os.Exit(1)
		}
		if gotVal != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:%s", i+1, expected, gotVal, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

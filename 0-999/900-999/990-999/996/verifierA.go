package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func expected(n int64) int64 {
	bills := []int64{100, 20, 10, 5, 1}
	cnt := int64(0)
	for _, b := range bills {
		cnt += n / b
		n %= b
	}
	return cnt
}

func runCase(bin, input string) (string, error) {
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	fixed := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		21, 23, 50, 57, 99, 100, 101, 123, 199, 200, 555, 999, 1000, 1001, 1020, 1023,
		10000, 1000000, 500000000, 1000000000}

	caseNum := 0
	for _, n := range fixed {
		input := fmt.Sprintf("%d\n", n)
		expect := fmt.Sprintf("%d", expected(n))
		got, err := runCase(bin, input)
		caseNum++
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", caseNum, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", caseNum, expect, got, input)
			os.Exit(1)
		}
	}

	for ; caseNum < 100; caseNum++ {
		n := rng.Int63n(1000000000) + 1
		input := fmt.Sprintf("%d\n", n)
		expect := fmt.Sprintf("%d", expected(n))
		got, err := runCase(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", caseNum+1, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", caseNum+1, expect, got, input)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed")
}

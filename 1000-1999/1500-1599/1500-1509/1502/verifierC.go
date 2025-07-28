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

func fib(n int) uint64 {
	if n <= 1 {
		return uint64(n)
	}
	a, b := uint64(0), uint64(1)
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}
func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesC.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	idx := 0
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		idx++
		n, _ := strconv.Atoi(line)
		expect := fib(n)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(fmt.Sprintf("%d\n", n))
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx, err, out.String())
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(out.String())
		got, err2 := strconv.ParseUint(gotStr, 10, 64)
		if err2 != nil || got != expect {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d got %s\n", idx, expect, gotStr)
			os.Exit(1)
		}
	}
	if err := sc.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

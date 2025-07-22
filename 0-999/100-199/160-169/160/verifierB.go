package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCaseB struct {
	n int
	s string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(43)
	var tests []testCaseB
	for i := 0; i < 100; i++ {
		n := rand.Intn(30) + 1
		b := make([]byte, 2*n)
		for j := range b {
			b[j] = byte('0' + rand.Intn(10))
		}
		tests = append(tests, testCaseB{n, string(b)})
	}
	tests = append(tests, testCaseB{1, "12"})
	tests = append(tests, testCaseB{2, "0000"})

	for i, t := range tests {
		input := fmt.Sprintf("%d\n%s\n", t.n, t.s)
		expect := solveB(strings.NewReader(input))
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, strings.TrimSpace(expect), strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func solveB(r io.Reader) string {
	in := bufio.NewReader(r)
	var n int
	var s string
	fmt.Fscan(in, &n)
	fmt.Fscan(in, &s)
	a := []byte(s[:n])
	b := []byte(s[n:])
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	less := true
	greater := true
	for i := 0; i < n; i++ {
		if a[i] >= b[i] {
			less = false
		}
		if a[i] <= b[i] {
			greater = false
		}
	}
	if less || greater {
		return "YES\n"
	}
	return "NO\n"
}

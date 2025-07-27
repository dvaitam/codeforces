package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveCase(s string) int {
	n := len(s)
	runs := make([]int, 0, n)
	cur := 1
	for i := 1; i < n; i++ {
		if s[i] == s[i-1] {
			cur++
		} else {
			runs = append(runs, cur)
			cur = 1
		}
	}
	runs = append(runs, cur)
	if len(runs) > 1 && s[0] == s[n-1] {
		runs[0] += runs[len(runs)-1]
		runs = runs[:len(runs)-1]
	}
	ans := 0
	for _, r := range runs {
		ans += r / 2
	}
	return ans
}

func generateTest() (string, string) {
	n := rand.Intn(10) + 1
	b := make([]byte, n)
	for i := range b {
		if rand.Intn(2) == 0 {
			b[i] = 'L'
		} else {
			b[i] = 'R'
		}
	}
	s := string(b)
	inp := fmt.Sprintf("%d\n%s\n", n, s)
	exp := fmt.Sprintf("%d\n", solveCase(s))
	return inp, exp
}

func referenceIO(t int) (string, string) {
	var in strings.Builder
	var out strings.Builder
	fmt.Fprintf(&in, "%d\n", t)
	for i := 0; i < t; i++ {
		ti, to := generateTest()
		in.WriteString(ti)
		out.WriteString(to)
	}
	return in.String(), out.String()
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return buf.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	rand.Seed(4)
	in, exp := referenceIO(100)
	out, err := runBinary(os.Args[1], in)
	if err != nil {
		fmt.Println("Runtime error:", err)
		os.Exit(1)
	}
	if strings.TrimSpace(out) != strings.TrimSpace(exp) {
		fmt.Println("Wrong Answer")
		fmt.Println("Expected:\n" + exp)
		fmt.Println("Got:\n" + out)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveCase(h []int64) []int64 {
	n := len(h)
	res := make([]int64, n)
	copy(res, h)
	for {
		changed := false
		for i := 0; i < n-1; i++ {
			if res[i+1]-res[i] >= 2 {
				m := (res[i+1] - res[i]) / 2
				res[i] += m
				res[i+1] -= m
				changed = true
			}
		}
		if !changed {
			break
		}
	}
	return res
}

func generateTest() (string, string) {
	n := rand.Intn(6) + 1
	h := make([]int64, n)
	base := int64(rand.Intn(5))
	for i := 0; i < n; i++ {
		base += int64(rand.Intn(3) + 1)
		h[i] = base
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range h {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')

	res := solveCase(h)
	var out strings.Builder
	for i, v := range res {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprintf(&out, "%d", v)
	}
	out.WriteByte('\n')
	return sb.String(), out.String()
}

func referenceIO(t int) (string, string) {
	var in strings.Builder
	var out strings.Builder
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	rand.Seed(6)
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

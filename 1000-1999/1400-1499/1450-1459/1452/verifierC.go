package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func solveC(s string) int {
	openRound := 0
	openSquare := 0
	ans := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			openRound++
		case '[':
			openSquare++
		case ')':
			if openRound > 0 {
				ans++
				openRound--
			}
		case ']':
			if openSquare > 0 {
				ans++
				openSquare--
			}
		}
	}
	return ans
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("run error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func randomString(n int) string {
	chars := []byte{'(', ')', '[', ']'}
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	const t = 100
	strs := make([]string, t)
	exp := make([]int, t)
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", t)
	for i := 0; i < t; i++ {
		s := randomString(rand.Intn(20) + 1)
		strs[i] = s
		exp[i] = solveC(s)
		b.WriteString(s)
		b.WriteByte('\n')
	}
	out, err := runBinary(binary, b.String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fields := strings.Fields(strings.TrimSpace(out))
	if len(fields) != t {
		fmt.Printf("expected %d lines, got %d\noutput:\n%s\n", t, len(fields), out)
		os.Exit(1)
	}
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil || v != exp[i] {
			fmt.Printf("test %d failed: input=%s expected=%d got=%s\n", i+1, strs[i], exp[i], f)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

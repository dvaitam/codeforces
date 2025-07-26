package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type test struct {
	s string
}

func genTests() []test {
	rand.Seed(1)
	tests := make([]test, 100)
	for i := range tests {
		n := rand.Intn(20) + 1
		var sb strings.Builder
		for j := 0; j < n; j++ {
			if rand.Intn(2) == 0 {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		}
		tests[i] = test{sb.String()}
	}
	return tests
}

func runBinary(path, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", path)
	} else {
		cmd = exec.CommandContext(ctx, path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// reference solve from 949A.go
func solveRef(s string) ([][]int, bool) {
	n := len(s)
	nxt := make([]int, n+2)
	vst := make([]bool, n+2)
	a0 := make([]int, 0, n)
	a1 := make([]int, 0, n)
	m := 0
	for i := 1; i <= n; i++ {
		c := s[i-1]
		if c == '0' {
			if len(a1) > 0 {
				u := a1[len(a1)-1]
				a1 = a1[:len(a1)-1]
				nxt[u] = i
				a0 = append(a0, i)
			} else {
				a0 = append(a0, i)
				m++
			}
		} else {
			if len(a0) > 0 {
				u := a0[len(a0)-1]
				a0 = a0[:len(a0)-1]
				nxt[u] = i
				a1 = append(a1, i)
			} else {
				return nil, false
			}
		}
	}
	if len(a1) > 0 {
		return nil, false
	}
	res := make([][]int, 0, m)
	for i := 1; i <= n; i++ {
		if vst[i] {
			continue
		}
		seq := []int{}
		for u := i; u != 0; u = nxt[u] {
			seq = append(seq, u)
			vst[u] = true
		}
		res = append(res, seq)
	}
	return res, true
}

func checkOutput(s string, out string, possible bool) bool {
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(out)))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return false
	}
	tok := scanner.Text()
	if tok == "-1" {
		if possible {
			return false
		}
		if scanner.Scan() {
			return false
		}
		return true
	}
	k, err := strconv.Atoi(tok)
	if err != nil || k <= 0 || k > len(s) {
		return false
	}
	used := make([]bool, len(s)+1)
	for i := 0; i < k; i++ {
		if !scanner.Scan() {
			return false
		}
		li, err := strconv.Atoi(scanner.Text())
		if err != nil || li <= 0 {
			return false
		}
		prev := 0
		for j := 0; j < li; j++ {
			if !scanner.Scan() {
				return false
			}
			idx, err := strconv.Atoi(scanner.Text())
			if err != nil || idx < 1 || idx > len(s) || idx <= prev {
				return false
			}
			if used[idx] {
				return false
			}
			used[idx] = true
			prev = idx
		}
	}
	if scanner.Scan() {
		return false
	}
	for i := 1; i <= len(s); i++ {
		if !used[i] {
			return false
		}
	}
	// validate zebra property
	scanner = bufio.NewScanner(strings.NewReader(strings.TrimSpace(out)))
	scanner.Split(bufio.ScanWords)
	scanner.Scan() // skip k
	for i := 0; i < k; i++ {
		scanner.Scan()
		li, _ := strconv.Atoi(scanner.Text())
		prev := -1
		for j := 0; j < li; j++ {
			scanner.Scan()
			x, _ := strconv.Atoi(scanner.Text())
			if j == 0 && s[x-1] != '0' {
				return false
			}
			if j == li-1 && s[x-1] != '0' {
				return false
			}
			if prev != -1 && s[x-1] == s[prev-1] {
				return false
			}
			prev = x
		}
	}
	if !possible {
		return false
	}
	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		input := t.s + "\n"
		out, err := runBinary(candidate, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		_, ok := solveRef(t.s)
		if !checkOutput(t.s, out, ok) {
			fmt.Printf("test %d failed\ninput:\n%s\noutput:\n%s\n", i+1, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

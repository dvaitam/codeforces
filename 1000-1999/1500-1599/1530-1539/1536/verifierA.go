package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func solveA(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	fmt.Fscan(in, &t)
	var out strings.Builder
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		minv := math.MaxInt64
		s := make(map[int]struct{})
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if a[i] < minv {
				minv = a[i]
			}
			s[a[i]] = struct{}{}
		}
		if minv < 0 {
			out.WriteString("NO\n")
			continue
		}
		out.WriteString("YES\n")
		sort.Ints(a)
		for i := 0; i < len(a); i++ {
			for j := i + 1; j < len(a); j++ {
				diff := a[j] - a[i]
				if diff < 0 {
					diff = -diff
				}
				if _, ok := s[diff]; !ok {
					s[diff] = struct{}{}
					a = append(a, diff)
					sort.Ints(a)
					i = -1
					break
				}
			}
		}
		res := make([]int, 0, len(s))
		for v := range s {
			res = append(res, v)
		}
		sort.Ints(res)
		out.WriteString(fmt.Sprintf("%d\n", len(res)))
		for i, v := range res {
			if i > 0 {
				out.WriteByte(' ')
			}
			out.WriteString(fmt.Sprintf("%d", v))
		}
		out.WriteByte('\n')
	}
	return strings.TrimSpace(out.String())
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() []string {
	r := rand.New(rand.NewSource(1))
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := r.Intn(10) + 2
		perm := r.Perm(201)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", perm[j]-100))
		}
		sb.WriteByte('\n')
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		expected := solveA(t)
		got, err := runBinary(bin, t)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("test %d failed. input: %sexpected %s got %s\n", i+1, t, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

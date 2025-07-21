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

func solveJ(input string) string {
	fields := bytes.Fields([]byte(input))
	if len(fields) < 2 {
		return "0\n"
	}
	s := fields[0]
	t := fields[1]
	n := len(s)
	if len(t) != n-1 {
		return "0\n"
	}
	pre := make([]bool, n)
	pre[0] = true
	for i := 1; i < n; i++ {
		if pre[i-1] && s[i-1] == t[i-1] {
			pre[i] = true
		}
	}
	suf := make([]bool, n)
	suf[n-1] = true
	for i := n - 2; i >= 0; i-- {
		if s[i+1] == t[i] && suf[i+1] {
			suf[i] = true
		}
	}
	var pos []int
	for i := 0; i < n; i++ {
		if pre[i] && suf[i] {
			pos = append(pos, i+1)
		}
	}
	if len(pos) == 0 {
		return "0\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(pos)))
	for i, v := range pos {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	return sb.String()
}

func randString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(3))
	}
	return string(b)
}

func generateCaseJ(rng *rand.Rand) string {
	n := rng.Intn(4) + 2
	s := randString(rng, n)
	p := rng.Intn(n)
	t := s[:p] + s[p+1:]
	return fmt.Sprintf("%s %s\n", s, t)
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierJ.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		cases[i] = generateCaseJ(rng)
	}
	for i, tc := range cases {
		expect := solveJ(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:%sq\ngot:%sq\n", i+1, tc, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

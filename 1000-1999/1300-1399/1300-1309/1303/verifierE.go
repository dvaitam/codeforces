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

func solveCase(s, t string) string {
	n := len(s)
	m := len(t)
	for i := 0; i <= m; i++ {
		used := make([]bool, n)
		p := 0
		for j := 0; j < n && p < i; j++ {
			if s[j] == t[p] {
				used[j] = true
				p++
			}
		}
		if p < i {
			continue
		}
		q := i
		for j := 0; j < n && q < m; j++ {
			if used[j] {
				continue
			}
			if s[j] == t[q] {
				q++
			}
		}
		if q == m {
			return "YES\n"
		}
	}
	return "NO\n"
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	m := rng.Intn(6) + 1
	sb1 := make([]byte, n)
	sb2 := make([]byte, m)
	for i := range sb1 {
		sb1[i] = byte(rng.Intn(26)) + 'a'
	}
	for i := range sb2 {
		sb2[i] = byte(rng.Intn(26)) + 'a'
	}
	s := string(sb1)
	t := string(sb2)
	in := fmt.Sprintf("1\n%s\n%s\n", s, t)
	out := solveCase(s, t)
	return in, out
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	got := strings.TrimSpace(buf.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected \n%s\ngot \n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

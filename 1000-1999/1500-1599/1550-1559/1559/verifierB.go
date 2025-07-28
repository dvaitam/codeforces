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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n int, s string) string {
	bs := []byte(s)
	l := -1
	for i := 0; i < n; i++ {
		if bs[i] != '?' {
			l = i
		}
	}
	if l == -1 {
		bs[0] = 'B'
		l = 0
	}
	for i := l - 1; i >= 0; i-- {
		if bs[i] == '?' {
			if bs[i+1] == 'B' {
				bs[i] = 'R'
			} else {
				bs[i] = 'B'
			}
		}
	}
	for i := l + 1; i < n; i++ {
		if bs[i] == '?' {
			if bs[i-1] == 'B' {
				bs[i] = 'R'
			} else {
				bs[i] = 'B'
			}
		}
	}
	return string(bs)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		r := rng.Intn(3)
		switch r {
		case 0:
			b[i] = 'B'
		case 1:
			b[i] = 'R'
		default:
			b[i] = '?'
		}
	}
	s := string(b)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %s\n", n, s))
	return sb.String(), expected(n, s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, want := generateCase(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

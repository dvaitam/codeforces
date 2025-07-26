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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(n int, m int, s, t string) string {
	best := n + 1
	var bestPos []int
	for i := 0; i <= m-n; i++ {
		diff := make([]int, 0)
		for j := 0; j < n; j++ {
			if s[j] != t[i+j] {
				diff = append(diff, j+1)
			}
		}
		if len(diff) < best {
			best = len(diff)
			bestPos = diff
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", best))
	for i, p := range bestPos {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", p))
	}
	if len(bestPos) > 0 {
		sb.WriteByte('\n')
	}
	return strings.TrimSpace(sb.String())
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + n // ensure m>=n
	letters := "abcdefghijklmnopqrstuvwxyz"
	makeStr := func(length int) string {
		b := make([]byte, length)
		for i := range b {
			b[i] = letters[rng.Intn(len(letters))]
		}
		return string(b)
	}
	s := makeStr(n)
	t := makeStr(m)
	input := fmt.Sprintf("%d %d\n%s\n%s\n", n, m, s, t)
	expected := solveCase(n, m, s, t)
	return input, expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

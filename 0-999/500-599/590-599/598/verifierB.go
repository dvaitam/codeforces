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

func applyQueries(s []byte, queries [][3]int) string {
	for _, q := range queries {
		l, r, k := q[0], q[1], q[2]
		l--
		length := r - l
		k %= length
		if k > 0 {
			tmp := append([]byte(nil), s[l:r]...)
			copy(s[l:r], append(tmp[length-k:], tmp[:length-k]...))
		}
	}
	return string(s)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(15) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	s := string(b)
	m := rng.Intn(20) + 1
	queries := make([][3]int, m)
	var sb strings.Builder
	sb.WriteString(s + "\n")
	sb.WriteString(strconv.Itoa(m) + "\n")
	for i := 0; i < m; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		k := rng.Intn(1000) + 1
		queries[i] = [3]int{l, r, k}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", l, r, k))
	}
	expected := applyQueries([]byte(s), queries)
	return sb.String(), expected + "\n"
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	exe := filepath.Join(dir, "oracleG")
	cmd := exec.Command("go", "build", "-o", exe, filepath.Join(dir, "1593G.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle: %v\n%s", err, out)
	}
	return exe, nil
}

func generateCase(rng *rand.Rand) (string, [][2]int) {
	n := rng.Intn(10) + 2
	chars := []byte{'(', ')', '[', ']'}
	s := make([]byte, n)
	for i := range s {
		s[i] = chars[rng.Intn(len(chars))]
	}
	q := rng.Intn(5) + 1
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(n-1) + 1
		r := rng.Intn(n-l) + l + 1
		if (r-l+1)%2 == 1 {
			if r < n {
				r++
			} else if l > 1 {
				l--
			}
		}
		queries[i] = [2]int{l, r}
	}
	return string(s), queries
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		s, queries := generateCase(rng)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(s)
		sb.WriteByte('\n')
		fmt.Fprintf(&sb, "%d\n", len(queries))
		for _, q := range queries {
			fmt.Fprintf(&sb, "%d %d\n", q[0], q[1])
		}
		input := sb.String()
		exp, err := runProg("./"+oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failure on case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := runProg(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

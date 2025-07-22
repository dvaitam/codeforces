package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleD")
	cmd := exec.Command("go", "build", "-o", oracle, "201D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func randWord(rng *rand.Rand) string {
	l := rng.Intn(3) + 1
	b := make([]byte, l)
	for i := 0; i < l; i++ {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 1
	words := make([]string, 0, n)
	used := map[string]bool{}
	for len(words) < n {
		w := randWord(rng)
		if !used[w] {
			used[w] = true
			words = append(words, w)
		}
	}
	m := rng.Intn(3) + 1
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(words[i])
	}
	b.WriteByte('\n')
	fmt.Fprintf(&b, "%d\n", m)
	for i := 0; i < m; i++ {
		k := rng.Intn(6) + 1
		fmt.Fprintf(&b, "%d", k)
		for j := 0; j < k; j++ {
			b.WriteByte(' ')
			b.WriteString(randWord(rng))
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		expected, err := runBinary(oracle, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if expected != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "case %d wrong answer\nexpected:\n%s\n\ngot:\n%s\ninput:\n%s", i+1, expected, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

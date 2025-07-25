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

func runCandidate(bin, input string) (string, error) {
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

func expected(memory []string, queries []string) string {
	var sb strings.Builder
	for _, q := range queries {
		found := false
		for _, m := range memory {
			if len(m) != len(q) {
				continue
			}
			diff := 0
			for i := 0; i < len(q) && diff <= 1; i++ {
				if q[i] != m[i] {
					diff++
				}
			}
			if diff == 1 {
				found = true
				break
			}
		}
		if found {
			sb.WriteString("YES\n")
		} else {
			sb.WriteString("NO\n")
		}
	}
	return strings.TrimSpace(sb.String())
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	memory := make([]string, n)
	for i := 0; i < n; i++ {
		memory[i] = randString(rng)
	}
	queries := make([]string, m)
	for i := 0; i < m; i++ {
		queries[i] = randString(rng)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, s := range memory {
		sb.WriteString(s + "\n")
	}
	for _, s := range queries {
		sb.WriteString(s + "\n")
	}
	expect := expected(memory, queries)
	return sb.String(), expect
}

func randString(rng *rand.Rand) string {
	l := rng.Intn(5) + 1
	b := make([]byte, l)
	for i := range b {
		b[i] = byte(rng.Intn(3)) + 'a'
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// deterministic edge cases
	mem := []string{"a", "b"}
	queries := []string{"c", "aa"}
	in := "2 2\na\nb\nc\naa\n"
	exp := expected(mem, queries)
	out, err := runCandidate(bin, in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "edge case failed: %v\n", err)
		os.Exit(1)
	}
	if out != exp {
		fmt.Fprintf(os.Stderr, "edge case failed: expected %s got %s\n", exp, out)
		os.Exit(1)
	}

	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

func countOccurrences(s, x string) int {
	m := len(x)
	if m == 0 || m > len(s) {
		return 0
	}
	doubled := x + x
	cnt := 0
	for i := 0; i+m <= len(s); i++ {
		sub := s[i : i+m]
		if strings.Contains(doubled, sub) {
			cnt++
		}
	}
	return cnt
}

func solveCase(s string, queries []string) string {
	res := make([]string, len(queries))
	for i, q := range queries {
		res[i] = fmt.Sprintf("%d", countOccurrences(s, q))
	}
	return strings.Join(res, "\n")
}

func generateCase(rng *rand.Rand) (string, string) {
	letters := "abc"
	ls := rng.Intn(8) + 1
	var sb strings.Builder
	for i := 0; i < ls; i++ {
		sb.WriteByte(letters[rng.Intn(len(letters))])
	}
	s := sb.String()
	q := rng.Intn(3) + 1
	queries := make([]string, q)
	for i := 0; i < q; i++ {
		lq := rng.Intn(5) + 1
		var qb strings.Builder
		for j := 0; j < lq; j++ {
			qb.WriteByte(letters[rng.Intn(len(letters))])
		}
		queries[i] = qb.String()
	}
	input := s + "\n" + fmt.Sprintf("%d\n", q) + strings.Join(queries, "\n") + "\n"
	expect := solveCase(s, queries)
	return input, expect
}

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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
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

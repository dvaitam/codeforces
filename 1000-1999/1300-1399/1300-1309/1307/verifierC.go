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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(input string) string {
	in := strings.NewReader(input)
	var s string
	fmt.Fscan(in, &s)
	cnt := make([]int64, 26)
	pairs := make([][]int64, 26)
	for i := 0; i < 26; i++ {
		pairs[i] = make([]int64, 26)
	}
	for i := 0; i < len(s); i++ {
		c := s[i] - 'a'
		for j := 0; j < 26; j++ {
			pairs[j][c] += cnt[j]
		}
		cnt[c]++
	}
	var ans int64
	for i := 0; i < 26; i++ {
		if cnt[i] > ans {
			ans = cnt[i]
		}
	}
	for i := 0; i < 26; i++ {
		for j := 0; j < 26; j++ {
			if pairs[i][j] > ans {
				ans = pairs[i][j]
			}
		}
	}
	return fmt.Sprintf("%d", ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(5))
	}
	s := string(b)
	return s + "\n", solve(s + "\n")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

func runCandidate(bin string, input string) (string, error) {
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

func expected(s string) string {
	n := len(s)
	ans := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		c := s[i]
		m := len(ans)
		if m >= 2 && ans[m-1] == c && ans[m-2] == c {
			continue
		}
		if m >= 3 && ans[m-3] == ans[m-2] && ans[m-1] == c {
			continue
		}
		ans = append(ans, c)
	}
	return string(ans)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		s := generateCase(rng)
		input := fmt.Sprintf("%s\n", s)
		want := expected(s)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

const mod int64 = 1000000007

// countSub returns the number of subsequences of t equal to sub.
func countSub(t, sub string) int64 {
	m := len(sub)
	dp := make([]int64, m+1)
	dp[0] = 1
	for i := 0; i < len(t); i++ {
		for j := m - 1; j >= 0; j-- {
			if sub[j] == t[i] {
				dp[j+1] += dp[j]
			}
		}
	}
	return dp[m]
}

func expectedAnswer(s, t string) int64 {
	var ans int64
	for i := 0; i < len(s); i++ {
		for j := i + 1; j <= len(s); j++ {
			sub := s[i:j]
			ans = (ans + countSub(t, sub)) % mod
		}
	}
	return ans % mod
}

// run executes the binary or go source with the given input and returns trimmed output.
func run(bin, input string) (string, error) {
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

func randomString(rng *rand.Rand, n int) string {
	letters := []byte{"a", "b", "c"}
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		slen := rng.Intn(5) + 1
		tlen := rng.Intn(6) + 1
		s := randomString(rng, slen)
		t := randomString(rng, tlen)
		input := fmt.Sprintf("%s\n%s\n", s, t)
		expected := expectedAnswer(s, t)
		gotStr, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(gotStr), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: cannot parse output: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got%mod != expected%mod {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, expected%mod, got%mod, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

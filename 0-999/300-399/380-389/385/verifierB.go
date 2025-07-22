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

func expectedCount(s string) int64 {
	n := len(s)
	nextOcc := make([]int, n)
	const inf = int(1e9)
	last := inf
	for i := n - 1; i >= 0; i-- {
		if i <= n-4 && s[i] == 'b' && s[i+1] == 'e' && s[i+2] == 'a' && s[i+3] == 'r' {
			last = i
		}
		nextOcc[i] = last
	}
	var ans int64
	for i := 0; i < n; i++ {
		p := nextOcc[i]
		if p <= n-4 {
			ans += int64(n - (p + 3))
		}
	}
	return ans
}

func randString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(100) + 1
	s := randString(rng, n)
	return s + "\n", expectedCount(s)
}

func runCase(bin, input string, expected int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

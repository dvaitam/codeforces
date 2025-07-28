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

func solveB(n, k int, s string) string {
	lower := make([]int, 26)
	upper := make([]int, 26)
	for _, ch := range s {
		if ch >= 'a' && ch <= 'z' {
			lower[ch-'a']++
		} else if ch >= 'A' && ch <= 'Z' {
			upper[ch-'A']++
		}
	}
	ans := 0
	for i := 0; i < 26; i++ {
		pairs := min(lower[i], upper[i])
		ans += pairs
		diff := abs(lower[i] - upper[i])
		extra := min(k, diff/2)
		ans += extra
		k -= extra
	}
	return fmt.Sprint(ans)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(50) + 1
	k := rng.Intn(n + 1)
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(letters[rng.Intn(len(letters))])
	}
	s := sb.String()
	input := fmt.Sprintf("1\n%d %d\n%s\n", n, k, s)
	expected := solveB(n, k, s)
	return input, expected
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

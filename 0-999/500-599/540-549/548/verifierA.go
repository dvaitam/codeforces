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

func isPalindrome(s string) bool {
	i, j := 0, len(s)-1
	for i < j {
		if s[i] != s[j] {
			return false
		}
		i++
		j--
	}
	return true
}

func solveCase(s string, k int) string {
	n := len(s)
	if k <= 0 || n%k != 0 {
		return "NO"
	}
	m := n / k
	for i := 0; i < k; i++ {
		start := i * m
		end := start + m
		if !isPalindrome(s[start:end]) {
			return "NO"
		}
	}
	return "YES"
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	sb := make([]byte, n)
	for i := range sb {
		sb[i] = byte('a' + rng.Intn(26))
	}
	k := rng.Intn(20) + 1
	input := fmt.Sprintf("%s\n%d\n", string(sb), k)
	expect := solveCase(string(sb), k)
	return input, expect
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := genCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

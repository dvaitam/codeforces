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

func maxBeauty(s string, n int) int {
	freq := make(map[byte]int)
	for i := 0; i < len(s); i++ {
		freq[s[i]]++
	}
	maxFreq := 0
	for _, v := range freq {
		if v > maxFreq {
			maxFreq = v
		}
	}
	l := len(s)
	if maxFreq == l && n == 1 {
		return l - 1
	}
	if maxFreq+n > l {
		return l
	}
	return maxFreq + n
}

func solve(n int, s1, s2, s3 string) string {
	b1 := maxBeauty(s1, n)
	b2 := maxBeauty(s2, n)
	b3 := maxBeauty(s3, n)
	if b1 > b2 && b1 > b3 {
		return "Kuro"
	} else if b2 > b1 && b2 > b3 {
		return "Shiro"
	} else if b3 > b1 && b3 > b2 {
		return "Katie"
	}
	return "Draw"
}

func randString(rng *rand.Rand, l int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, l)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5)
	l := rng.Intn(10) + 1
	s1 := randString(rng, l)
	s2 := randString(rng, l)
	s3 := randString(rng, l)
	input := fmt.Sprintf("%d\n%s\n%s\n%s\n", n, s1, s2, s3)
	expected := solve(n, s1, s2, s3)
	return input, expected
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
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, in, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

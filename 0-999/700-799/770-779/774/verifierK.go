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

func isVowel(c byte) bool {
	switch c {
	case 'a', 'e', 'i', 'o', 'u', 'y':
		return true
	}
	return false
}

func expected(n int, s string) string {
	var res []byte
	i := 0
	for i < n {
		ch := s[i]
		j := i + 1
		for j < n && s[j] == ch {
			j++
		}
		runLen := j - i
		if isVowel(ch) {
			if (ch == 'e' || ch == 'o') && runLen == 2 {
				res = append(res, ch, ch)
			} else {
				res = append(res, ch)
			}
		} else {
			for k := 0; k < runLen; k++ {
				res = append(res, ch)
			}
		}
		i = j
	}
	return string(res)
}

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

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	if rng.Float64() < 0.1 {
		n = rng.Intn(50) + 1
	}
	letters := []byte("abcdefghijklmnopqrstuvwxyz")
	sb := make([]byte, n)
	for i := 0; i < n; i++ {
		sb[i] = letters[rng.Intn(len(letters))]
	}
	s := string(sb)
	input := fmt.Sprintf("%d\n%s\n", n, s)
	return input, expected(n, s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierK.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := genCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, strings.TrimSpace(out), input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

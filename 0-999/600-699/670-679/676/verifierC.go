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

func maxLen(s string, k int, ch byte) int {
	left := 0
	cnt := 0
	best := 0
	for right := 0; right < len(s); right++ {
		if s[right] != ch {
			cnt++
		}
		for cnt > k {
			if s[left] != ch {
				cnt--
			}
			left++
		}
		if cur := right - left + 1; cur > best {
			best = cur
		}
	}
	return best
}

func solve(n, k int, s string) int {
	a := maxLen(s, k, 'a')
	b := maxLen(s, k, 'b')
	if a > b {
		return a
	}
	return b
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	k := rng.Intn(n + 1)
	b := make([]byte, n)
	for i := range b {
		if rng.Intn(2) == 0 {
			b[i] = 'a'
		} else {
			b[i] = 'b'
		}
	}
	s := string(b)
	input := fmt.Sprintf("%d %d\n%s\n", n, k, s)
	expected := fmt.Sprintf("%d\n", solve(n, k, s))
	return input, expected
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, strings.TrimSpace(expect), strings.TrimSpace(out), input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

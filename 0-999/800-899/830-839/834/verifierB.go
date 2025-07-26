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

func expect(n, k int, s string) string {
	first := make([]int, 26)
	last := make([]int, 26)
	for i := 0; i < 26; i++ {
		first[i] = -1
		last[i] = -1
	}
	for i, ch := range s {
		idx := int(ch - 'A')
		if first[idx] == -1 {
			first[idx] = i
		}
		last[idx] = i
	}
	open := make([]bool, 26)
	openCount := 0
	for i, ch := range s {
		idx := int(ch - 'A')
		if i == first[idx] {
			if !open[idx] {
				open[idx] = true
				openCount++
			}
		}
		if openCount > k {
			return "YES"
		}
		if i == last[idx] {
			if open[idx] {
				open[idx] = false
				openCount--
			}
		}
	}
	return "NO"
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
	n := rng.Intn(20) + 1
	if rng.Float64() < 0.2 {
		n = rng.Intn(100) + 1
	}
	k := rng.Intn(26) + 1
	sb := make([]byte, n)
	for i := 0; i < n; i++ {
		sb[i] = byte('A' + rng.Intn(26))
	}
	s := string(sb)
	input := fmt.Sprintf("%d %d\n%s\n", n, k, s)
	exp := expect(n, k, s)
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

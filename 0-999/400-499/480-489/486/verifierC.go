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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func expected(n, p int, s string) int {
	// convert to left half
	half := n / 2
	if p > half {
		p = n - p + 1
	}
	totalChange := 0
	left, right := n, 0
	for i := 1; i <= half; i++ {
		a := s[i-1]
		b := s[n-i]
		diff := int(a) - int(b)
		if diff < 0 {
			diff = -diff
		}
		cost := min(diff, 26-diff)
		if cost > 0 {
			totalChange += cost
			if i < left {
				left = i
			}
			if i > right {
				right = i
			}
		}
	}
	if totalChange == 0 {
		return 0
	}
	span := right - left
	mv := span + min(abs(p-left), abs(p-right))
	return totalChange + mv
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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

func randomString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	const cases = 100
	for i := 0; i < cases; i++ {
		n := rng.Intn(20) + 1
		p := rng.Intn(n) + 1
		s := randomString(rng, n)
		input := fmt.Sprintf("%d %d\n%s\n", n, p, s)
		want := fmt.Sprintf("%d", expected(n, p, s))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}

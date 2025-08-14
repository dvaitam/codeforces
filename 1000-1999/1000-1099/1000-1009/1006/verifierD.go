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

// get returns the minimum number of changes required to make the
// four characters form two equal pairs. Only the first and third
// characters (originally from string "a") can be changed.
func get(a, b, c, d byte) int {
	// If the characters from string "b" are the same, we only need to
	// ensure the characters from string "a" match each other.
	if b == d {
		if a == c {
			return 0
		}
		return 1
	}
	// Otherwise the characters of string "b" are different, so the final
	// configuration must contain both b and d once in each string. We can
	// choose which position gets which letter and minimise the number of
	// modifications in "a".
	cost1 := 0
	if a != b {
		cost1++
	}
	if c != d {
		cost1++
	}
	cost2 := 0
	if a != d {
		cost2++
	}
	if c != b {
		cost2++
	}
	if cost1 < cost2 {
		return cost1
	}
	return cost2
}

func solve(s0, s1 string) int {
	n := len(s0)
	ans := 0
	for i := 0; i < n/2; i++ {
		ans += get(s0[i], s1[i], s0[n-i-1], s1[n-i-1])
	}
	if n%2 == 1 && s0[n/2] != s1[n/2] {
		ans++
	}
	return ans
}

func randString(rng *rand.Rand, n int) string {
	letters := []byte("abcdefghijklmnopqrstuvwxyz")
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(30) + 1
		s0 := randString(rng, n)
		s1 := randString(rng, n)
		input := fmt.Sprintf("%d\n%s\n%s\n", n, s0, s1)
		expected := fmt.Sprintf("%d", solve(s0, s1))
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

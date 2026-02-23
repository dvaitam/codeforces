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

var prime []bool

func init() {
	const N = 1000
	prime = make([]bool, N+1)
	for i := 2; i <= N; i++ {
		prime[i] = true
	}
	for i := 2; i*i <= N; i++ {
		if prime[i] {
			for j := i * i; j <= N; j += i {
				prime[j] = false
			}
		}
	}
}

func solveCase(s string) (int, string) {
	n := len(s)
	cnt := make([]int, 10)
	digits := make([]int, n)
	for i, ch := range s {
		d := int(ch - '0')
		digits[i] = d
		cnt[d]++
	}
	if cnt[1] > 0 || cnt[4] > 0 || cnt[6] > 0 || cnt[8] > 0 || cnt[9] > 0 {
		switch {
		case cnt[1] > 0:
			return 1, "1"
		case cnt[4] > 0:
			return 1, "4"
		case cnt[6] > 0:
			return 1, "6"
		case cnt[8] > 0:
			return 1, "8"
		default:
			return 1, "9"
		}
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			num := digits[i]*10 + digits[j]
			if !prime[num] {
				return 2, fmt.Sprintf("%d%d", digits[i], digits[j])
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				num := digits[i]*100 + digits[j]*10 + digits[k]
				if !prime[num] {
					return 3, fmt.Sprintf("%d%d%d", digits[i], digits[j], digits[k])
				}
			}
		}
	}
	return 0, ""
}

func validateOutput(s, output string) error {
	lines := strings.Fields(output)
	if len(lines) != 2 {
		return fmt.Errorf("expected 2 lines of output, got %d", len(lines))
	}

	var c int
	if _, err := fmt.Sscan(lines[0], &c); err != nil {
		return fmt.Errorf("invalid count line: %v", err)
	}

	numStr := lines[1]
	if c != len(numStr) {
		return fmt.Errorf("count %d doesn't match digit string length %d (%q)", c, len(numStr), numStr)
	}

	// Check numStr is a subsequence of s
	j := 0
	for i := 0; i < len(s) && j < len(numStr); i++ {
		if s[i] == numStr[j] {
			j++
		}
	}
	if j != len(numStr) {
		return fmt.Errorf("%q is not a subsequence of %q", numStr, s)
	}

	// Parse number and check it's not prime
	num := 0
	for _, ch := range numStr {
		num = num*10 + int(ch-'0')
	}
	if len(numStr) < len(prime) && prime[num] {
		return fmt.Errorf("output number %d is prime", num)
	}

	return nil
}

func generateCase(rng *rand.Rand) (string, string) {
	for {
		k := rng.Intn(9) + 2 // 2..10 digits
		var digits strings.Builder
		for i := 0; i < k; i++ {
			digits.WriteByte(byte(rng.Intn(9) + '1'))
		}
		s := digits.String()
		lenAns, _ := solveCase(s)
		if lenAns == 0 {
			continue // skip unsolvable cases
		}
		input := fmt.Sprintf("1\n%d\n%s\n", k, s)
		return input, s
	}
}

func runCase(exe, input, s string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	return validateOutput(s, outStr)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, s := generateCase(rng)
		if err := runCase(exe, in, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

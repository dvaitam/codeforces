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

func generateCase(rng *rand.Rand) (string, string) {
	k := rng.Intn(9) + 2 // 2..10 digits
	var digits strings.Builder
	for i := 0; i < k; i++ {
		digits.WriteByte(byte(rng.Intn(9) + '1'))
	}
	s := digits.String()
	input := fmt.Sprintf("1\n%d\n%s\n", k, s)
	lenAns, ans := solveCase(s)
	expected := fmt.Sprintf("%d\n%s\n", lenAns, ans)
	return input, expected
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if outStr != exp {
		return fmt.Errorf("expected %q got %q", exp, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

func count(n, y int64) int64 {
	length := int64(1)
	if y%2 == 0 {
		length = 2
	}
	var res int64
	cur := y
	l := length
	for cur <= n {
		end := cur + l - 1
		if end > n {
			res += n - cur + 1
		} else {
			res += l
		}
		if cur > n/2 {
			break
		}
		cur <<= 1
		l <<= 1
	}
	return res
}

func maxY(n, k, parity int64) int64 {
	var low int64
	if parity == 1 {
		low = 1
	} else {
		low = 2
	}
	high := n
	if high%2 != parity {
		high--
	}
	var ans int64
	for low <= high {
		mid := (low + high) / 2
		if mid%2 != parity {
			mid++
			if mid > high {
				break
			}
		}
		if count(n, mid) >= k {
			if mid > ans {
				ans = mid
			}
			low = mid + 2
		} else {
			high = mid - 2
		}
	}
	return ans
}

func expected(n, k int64) string {
	odd := maxY(n, k, 1)
	even := maxY(n, k, 0)
	if odd > even {
		return fmt.Sprintf("%d", odd)
	}
	return fmt.Sprintf("%d", even)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Int63n(1_000_000_000_000) + 1
	k := rng.Int63n(n) + 1
	input := fmt.Sprintf("%d %d\n", n, k)
	return input, expected(n, k)
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
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		if err := runCase(exe, input, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

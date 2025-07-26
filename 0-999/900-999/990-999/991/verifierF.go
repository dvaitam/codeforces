package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const limitF int64 = 10000000000

var memoF = make(map[int64]string)

func powInt(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if a > limitF/res {
			return limitF + 1
		}
		res *= a
		b--
	}
	return res
}

func root(n int64, b int) int64 {
	if b <= 1 {
		return -1
	}
	x := int64(math.Pow(float64(n), 1.0/float64(b)) + 0.5)
	if x < 2 {
		return -1
	}
	for {
		p := powInt(x, int64(b))
		if p == n {
			return x
		}
		if p > n {
			x--
		} else {
			x++
		}
		if x < 2 {
			return -1
		}
		if powInt(x, int64(b)) == n {
			return x
		}
		if powInt(x, int64(b)) > n && powInt(x-1, int64(b)) < n {
			break
		}
	}
	return -1
}

func solve(n int64) string {
	if v, ok := memoF[n]; ok {
		return v
	}
	best := strconv.FormatInt(n, 10)
	bestLen := len(best)
	for b := 2; b <= 34; b++ {
		base := root(n, b)
		if base > 1 {
			s1 := solve(base)
			cand := s1 + "^" + strconv.Itoa(b)
			if len(cand) < bestLen {
				best = cand
				bestLen = len(cand)
			}
		}
	}
	for i := int64(2); i*i <= n; i++ {
		if n%i == 0 {
			s1 := solve(i)
			s2 := solve(n / i)
			cand := s1 + "*" + s2
			if len(cand) < bestLen {
				best = cand
				bestLen = len(cand)
			}
		}
	}
	for k := 1; k <= 10; k++ {
		p := int64(math.Pow10(k))
		if p > n {
			break
		}
		q := n / p
		r := n % p
		if q == 0 {
			continue
		}
		left := solve(q) + "*" + fmt.Sprintf("10^%d", k)
		if r == 0 {
			if len(left) < bestLen {
				best = left
				bestLen = len(left)
			}
		} else {
			right := solve(r)
			cand := left + "+" + right
			if len(cand) < bestLen {
				best = cand
				bestLen = len(cand)
			}
		}
	}
	for r := int64(1); r <= 9 && r < n; r++ {
		s1 := solve(n - r)
		cand := s1 + "+" + strconv.FormatInt(r, 10)
		if len(cand) < bestLen {
			best = cand
			bestLen = len(cand)
		}
	}
	memoF[n] = best
	return best
}

func generateCaseF(rng *rand.Rand) (string, string) {
	n := rng.Int63n(1_000_000_0000) + 1
	input := fmt.Sprintf("%d\n", n)
	expected := solve(n)
	return input, expected
}

func runCaseF(bin, input, expected string) error {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseF(rng)
		if err := runCaseF(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

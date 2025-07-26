package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func mulDiv(a, b, div uint64) uint64 {
	hi, lo := bits.Mul64(a, b)
	q, _ := bits.Div64(hi, lo, div)
	return q
}

func expected(n int, k uint64, arr []uint64) string {
	rev := make([]uint64, n)
	for i := 0; i < n; i++ {
		rev[n-1-i] = arr[i]
	}
	limit := k
	has := func(t uint64) bool {
		var sum uint64
		comb := uint64(1)
		for j := 0; j < n; j++ {
			if j > 0 {
				if t == 0 {
					break
				}
				comb = mulDiv(comb, t+uint64(j-1), uint64(j))
				if comb > limit {
					comb = limit + 1
				}
			}
			if rev[j] != 0 {
				if comb > limit/rev[j] {
					return true
				}
				sum += comb * rev[j]
				if sum >= limit {
					return true
				}
			}
			if comb == 0 {
				break
			}
		}
		return sum >= limit
	}
	lo, hi := uint64(0), k
	for lo < hi {
		mid := (lo + hi) / 2
		if has(mid) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return fmt.Sprintf("%d", lo)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	k := uint64(rng.Intn(1000) + 1)
	arr := make([]uint64, n)
	for i := range arr {
		arr[i] = uint64(rng.Intn(1000))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", arr[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), expected(n, k, arr)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

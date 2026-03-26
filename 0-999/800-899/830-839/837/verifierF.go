package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func reaches(i uint64, rev []uint64, k uint64) bool {
	var sum uint64
	t := uint64(1)
	m := len(rev)

	for d, a := range rev {
		if a != 0 {
			rem := k - sum
			hi, lo := bits.Mul64(t, a)
			if hi != 0 || lo >= rem {
				return true
			}
			sum += lo
		}

		if t == k && d+1 < m {
			return true
		}

		if d+1 < m && t < k {
			num := i + uint64(d)
			den := uint64(d + 1)
			hi, lo := bits.Mul64(t, num)
			if hi >= den {
				t = k
			} else {
				q, _ := bits.Div64(hi, lo, den)
				if q >= k {
					t = k
				} else {
					t = q
				}
			}
		}
	}

	return false
}

func expected(n int, k uint64, arr []uint64) string {
	var maxv uint64
	left := n
	for i := 0; i < n; i++ {
		if arr[i] > maxv {
			maxv = arr[i]
		}
		if arr[i] > 0 && left == n {
			left = i
		}
	}

	if maxv >= k {
		return "0"
	}

	m := n - left
	rev := make([]uint64, m)
	for i := 0; i < m; i++ {
		rev[i] = arr[n-1-i]
	}

	l, r := uint64(1), k
	for l < r {
		mid := l + (r-l)/2
		if reaches(mid, rev, k) {
			r = mid
		} else {
			l = mid + 1
		}
	}

	return strconv.FormatUint(l, 10)
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

package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func mulMod(a, b, mod uint64) uint64 {
	res := new(big.Int).Mul(new(big.Int).SetUint64(a), new(big.Int).SetUint64(b))
	res.Mod(res, new(big.Int).SetUint64(mod))
	return res.Uint64()
}

func powMod(a, d, mod uint64) uint64 {
	res := uint64(1)
	for d > 0 {
		if d&1 == 1 {
			res = mulMod(res, a, mod)
		}
		a = mulMod(a, a, mod)
		d >>= 1
	}
	return res
}

func isPrime(n uint64) bool {
	if n < 2 {
		return false
	}
	small := []uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37}
	for _, p := range small {
		if n%p == 0 {
			return n == p
		}
	}
	d := n - 1
	s := 0
	for d&1 == 0 {
		d >>= 1
		s++
	}
	bases := []uint64{2, 325, 9375, 28178, 450775, 9780504, 1795265022}
	for _, a := range bases {
		if a%n == 0 {
			continue
		}
		x := powMod(a%n, d, n)
		if x == 1 || x == n-1 {
			continue
		}
		skip := false
		for r := 1; r < s; r++ {
			x = mulMod(x, x, n)
			if x == n-1 {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		return false
	}
	return true
}

func gcd(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func pollardsRho(n uint64) uint64 {
	if n%2 == 0 {
		return 2
	}
	if n%3 == 0 {
		return 3
	}
	for {
		c := uint64(rand.Int63n(int64(n-1))) + 1
		x := uint64(rand.Int63n(int64(n)))
		y := x
		d := uint64(1)
		for d == 1 {
			x = (mulMod(x, x, n) + c) % n
			y = (mulMod(y, y, n) + c) % n
			y = (mulMod(y, y, n) + c) % n
			if x > y {
				d = gcd(x-y, n)
			} else {
				d = gcd(y-x, n)
			}
			if d == n {
				break
			}
		}
		if d > 1 && d < n {
			return d
		}
	}
}

func factor(n uint64, res *[]uint64) {
	if n == 1 {
		return
	}
	if isPrime(n) {
		*res = append(*res, n)
		return
	}
	d := pollardsRho(n)
	factor(d, res)
	factor(n/d, res)
}

func solve(a []uint64) int {
	n := len(a)
	omega := make([]int, n)
	for i := 0; i < n; i++ {
		var fs []uint64
		factor(a[i], &fs)
		omega[i] = len(fs)
	}
	isChildable := make([]bool, n)
	for j := 0; j < n; j++ {
		for i := 0; i < n; i++ {
			if i != j && a[i]%a[j] == 0 {
				isChildable[j] = true
				break
			}
		}
	}
	B := make([]int, 0, n)
	unCount := 0
	for i := 0; i < n; i++ {
		if isChildable[i] {
			B = append(B, i)
		} else {
			unCount++
		}
	}
	sumOmegaComp := 0
	for i := 0; i < n; i++ {
		if omega[i] > 1 {
			sumOmegaComp += omega[i]
		}
	}
	best := int(1e18)
	m := len(B)
	for mask := 0; mask < (1 << m); mask++ {
		sumSel := 0
		bits := 0
		for k := 0; k < m; k++ {
			if mask&(1<<k) != 0 {
				sumSel += omega[B[k]]
				bits++
			}
		}
		roots := unCount + (m - bits)
		penalty := 0
		if roots > 1 {
			penalty = 1
		}
		total := n + sumOmegaComp - sumSel + penalty
		if total < best {
			best = total
		}
	}
	return best
}

func generateTest(rng *rand.Rand) (string, int) {
	n := rng.Intn(8) + 1
	a := make([]uint64, n)
	for i := 0; i < n; i++ {
		a[i] = uint64(rng.Intn(1000000) + 2)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	ans := solve(a)
	return sb.String(), ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	rand.Seed(time.Now().UnixNano())
	const tests = 100
	for t := 1; t <= tests; t++ {
		inp, want := generateTest(rng)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(inp)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Test %d: runtime error: %v\n%s", t, err, out.String())
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(out.String())
		var got int
		if _, err := fmt.Sscan(gotStr, &got); err != nil {
			fmt.Fprintf(os.Stderr, "Test %d: failed to parse output: %v\nOutput: %s\n", t, err, gotStr)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected: %d\nGot: %d\n", t, inp, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

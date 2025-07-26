package main

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
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

func factorPrimePower(m uint64) (uint64, int) {
	mBig := new(big.Int).SetUint64(m)
	for t := 60; t >= 2; t-- {
		r := uint64(math.Round(math.Pow(float64(m), 1.0/float64(t))))
		for d := int64(r) - 1; d <= int64(r)+1; d++ {
			if d <= 1 {
				continue
			}
			cand := new(big.Int).SetInt64(d)
			pow := new(big.Int).Exp(cand, big.NewInt(int64(t)), nil)
			if pow.Cmp(mBig) == 0 {
				return uint64(d), t
			}
		}
	}
	return m, 1
}

func factorDistinct(x uint64) []uint64 {
	var fs []uint64
	if x%2 == 0 {
		fs = append(fs, 2)
		for x%2 == 0 {
			x /= 2
		}
	}
	for i := uint64(3); i*i <= x; i += 2 {
		if x%i == 0 {
			fs = append(fs, i)
			for x%i == 0 {
				x /= i
			}
		}
	}
	if x > 1 {
		fs = append(fs, x)
	}
	return fs
}

func modPow(a, e, m uint64) uint64 {
	return new(big.Int).Exp(new(big.Int).SetUint64(a), new(big.Int).SetUint64(e), new(big.Int).SetUint64(m)).Uint64()
}

func solve(n, m, p uint64) []uint64 {
	q, t := factorPrimePower(m)
	phi := uint64(1)
	for i := 0; i < t-1; i++ {
		phi *= q
	}
	phi *= (q - 1)
	if p%q == 0 {
		if phi < 1+n {
			return []uint64{math.MaxUint64}
		}
		res := make([]uint64, 0, n)
		for x := uint64(2); x < m && uint64(len(res)) < n; x++ {
			if x%q == 0 {
				continue
			}
			res = append(res, x)
		}
		return res
	}
	var primes []uint64
	if t > 1 {
		primes = append(primes, q)
	}
	primes = append(primes, factorDistinct(q-1)...)
	mset := make(map[uint64]bool)
	for _, v := range primes {
		mset[v] = true
	}
	primes = primes[:0]
	for v := range mset {
		primes = append(primes, v)
	}
	ord := phi
	for _, f := range primes {
		for ord%f == 0 {
			if modPow(p, ord/f, m) == 1 {
				ord /= f
			} else {
				break
			}
		}
	}
	if phi < ord+n {
		return []uint64{math.MaxUint64}
	}
	mBig := new(big.Int).SetUint64(m)
	ordBig := new(big.Int).SetUint64(ord)
	oneBig := big.NewInt(1)
	xBig := new(big.Int)
	res := make([]uint64, 0, n)
	for x := uint64(1); x < m && uint64(len(res)) < n; x++ {
		if x%q == 0 {
			continue
		}
		xBig.SetUint64(x)
		if new(big.Int).Exp(xBig, ordBig, mBig).Cmp(oneBig) == 0 {
			continue
		}
		res = append(res, x)
	}
	return res
}

func genTest(rng *rand.Rand) (string, string) {
	n := uint64(rng.Intn(3) + 1)
	m := uint64(rng.Intn(50) + 2)
	p := uint64(rng.Intn(int(m-1)) + 1)
	ans := solve(n, m, p)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, p))
	var out strings.Builder
	if len(ans) == 1 && ans[0] == math.MaxUint64 {
		out.WriteString("-1")
	} else {
		for i, v := range ans {
			if i > 0 {
				out.WriteByte('\n')
			}
			out.WriteString(fmt.Sprintf("%d", v))
		}
	}
	return sb.String(), out.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expected := genTest(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

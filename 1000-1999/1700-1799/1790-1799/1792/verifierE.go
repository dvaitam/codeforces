package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
)

func factorize(x int64) map[int64]int {
	res := make(map[int64]int)
	for p := int64(2); p*p <= x; p++ {
		for x%p == 0 {
			res[p]++
			x /= p
		}
	}
	if x > 1 {
		res[x]++
	}
	return res
}

func factorWithPrimes(d int64, primes []int64) []int {
	exps := make([]int, len(primes))
	for i, p := range primes {
		for d%p == 0 {
			exps[i]++
			d /= p
		}
	}
	return exps
}

func findMinDiv(primes []int64, exps []int, idx int, cur, lo, hi int64, best *int64) {
	if cur > hi || cur >= *best {
		return
	}
	if idx == len(primes) {
		if cur >= lo && cur < *best {
			*best = cur
		}
		return
	}
	val := int64(1)
	for i := 0; i <= exps[idx]; i++ {
		findMinDiv(primes, exps, idx+1, cur*val, lo, hi, best)
		val *= primes[idx]
		if cur*val > hi {
			break
		}
	}
}

func minimalRow(n int64, primes []int64, d int64) int64 {
	if d > n*n {
		return 0
	}
	if d <= n {
		return 1
	}
	lo := (d + n - 1) / n
	exps := factorWithPrimes(d, primes)
	best := int64(1<<63 - 1)
	findMinDiv(primes, exps, 0, 1, lo, n, &best)
	if best == int64(1<<63-1) {
		return 0
	}
	return best
}

func solveCaseE(n, m1, m2 int64) (int64, int64) {
	fac := factorize(m1)
	for p, e := range factorize(m2) {
		fac[p] += e
	}
	type pe struct {
		p int64
		e int
	}
	pes := make([]pe, 0, len(fac))
	for p, e := range fac {
		pes = append(pes, pe{p, e})
	}
	sort.Slice(pes, func(i, j int) bool { return pes[i].p < pes[j].p })
	primes := make([]int64, len(pes))
	exps := make([]int, len(pes))
	for i, v := range pes {
		primes[i] = v.p
		exps[i] = v.e
	}
	divisors := []int64{1}
	for i, p := range primes {
		cnt := exps[i]
		curSize := len(divisors)
		pow := int64(1)
		for e := 1; e <= cnt; e++ {
			pow *= p
			for j := 0; j < curSize; j++ {
				divisors = append(divisors, divisors[j]*pow)
			}
		}
	}
	var present int64
	var xorVal int64
	for _, d := range divisors {
		row := minimalRow(n, primes, d)
		if row != 0 {
			present++
		}
		xorVal ^= row
	}
	return present, xorVal
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s /path/to/binary\n", os.Args[0])
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(46)
	const t = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	expected := make([][2]int64, t)
	for i := 0; i < t; i++ {
		n := int64(rand.Intn(50) + 1)
		m1 := int64(rand.Intn(50) + 1)
		m2 := int64(rand.Intn(50) + 1)
		fmt.Fprintf(&input, "%d %d %d\n", n, m1, m2)
		a, b := solveCaseE(n, m1, m2)
		expected[i][0] = a
		expected[i][1] = b
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	outBytes, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "binary execution failed:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(outBytes))
	scanner.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		for j := 0; j < 2; j++ {
			if !scanner.Scan() {
				fmt.Printf("not enough output on test %d value %d\n", i+1, j+1)
				os.Exit(1)
			}
			got, err := strconv.ParseInt(scanner.Text(), 10, 64)
			if err != nil {
				fmt.Printf("invalid integer on test %d value %d: %v\n", i+1, j+1, err)
				os.Exit(1)
			}
			if got != expected[i][j] {
				fmt.Printf("mismatch on test %d value %d: expected %d got %d\n", i+1, j+1, expected[i][j], got)
				os.Exit(1)
			}
		}
	}
	if scanner.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed.")
}

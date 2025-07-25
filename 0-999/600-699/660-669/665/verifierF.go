package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var primes []uint64
var piSmall []uint64
var phiMemo map[phiKey]uint64
var piMemo map[uint64]uint64

type phiKey struct {
	x uint64
	s int
}

func sieve(limit int) {
	isPrime := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= limit; i++ {
		if isPrime[i] {
			for j := i * i; j <= limit; j += i {
				isPrime[j] = false
			}
		}
	}
	for i := 2; i <= limit; i++ {
		if isPrime[i] {
			primes = append(primes, uint64(i))
		}
	}
	piSmall = make([]uint64, limit+1)
	idx := 0
	for i := 1; i <= limit; i++ {
		if idx < len(primes) && primes[idx] == uint64(i) {
			idx++
		}
		piSmall[i] = uint64(idx)
	}
}

func isqrt(n uint64) uint64 {
	r := uint64(math.Sqrt(float64(n)))
	for (r+1)*(r+1) <= n {
		r++
	}
	for r*r > n {
		r--
	}
	return r
}

func icbrt(n uint64) uint64 {
	r := uint64(math.Cbrt(float64(n)))
	for (r+1)*(r+1)*(r+1) <= n {
		r++
	}
	for r*r*r > n {
		r--
	}
	return r
}

func phi(x uint64, s int) uint64 {
	if s == 0 {
		return x
	}
	if s == 1 {
		return x - x/primes[0]
	}
	if s == 2 {
		p1, p2 := primes[0], primes[1]
		return x - x/p1 - x/p2 + x/(p1*p2)
	}
	if x < primes[s-1]*primes[s-1] {
		return primePi(x) - uint64(s) + 1
	}
	key := phiKey{x, s}
	if val, ok := phiMemo[key]; ok {
		return val
	}
	res := phi(x, s-1) - phi(x/primes[s-1], s-1)
	phiMemo[key] = res
	return res
}

func primePi(n uint64) uint64 {
	if n < uint64(len(piSmall)) {
		return piSmall[n]
	}
	if val, ok := piMemo[n]; ok {
		return val
	}
	a := primePi(uint64(math.Sqrt(math.Sqrt(float64(n)))))
	b := primePi(uint64(math.Sqrt(float64(n))))
	c := primePi(uint64(math.Cbrt(float64(n))))
	res := phi(n, int(a)) + (b+a-2)*(b-a+1)/2
	for i := a + 1; i <= b; i++ {
		p := primes[i-1]
		w := n / p
		res -= primePi(w)
		if i <= c {
			bi := primePi(uint64(math.Sqrt(float64(w))))
			for j := i; j <= bi; j++ {
				res -= primePi(w/primes[j-1]) - (j - 1)
			}
		}
	}
	piMemo[n] = res
	return res
}

func countFourDivisors(n uint64) uint64 {
	cbrt := icbrt(n)
	res := primePi(cbrt)
	for i, p := range primes {
		if p*p > n {
			break
		}
		limit := n / p
		cnt := primePi(limit) - uint64(i+1)
		if cnt > 0 {
			res += cnt
		}
	}
	return res
}

func expected(n uint64) string {
	limit := int(isqrt(n)) + 10
	sieve(limit)
	phiMemo = make(map[phiKey]uint64)
	piMemo = make(map[uint64]uint64)
	ans := countFourDivisors(n)
	return fmt.Sprintf("%d", ans)
}

func runCase(exe, input, exp string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp = strings.TrimSpace(exp)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesF.txt")
	if err != nil {
		fmt.Println("could not read testcasesF.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		val, _ := strconv.ParseUint(scan.Text(), 10, 64)
		input := fmt.Sprintf("%d\n", val)
		exp := expected(val) + "\n"
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", caseIdx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

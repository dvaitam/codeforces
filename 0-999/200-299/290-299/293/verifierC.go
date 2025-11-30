package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Embedded testcases (previously in testcasesC.txt) to keep verifier self contained.
const rawTestcasesC = `
49387
31425
39819
78524
95487
86205
92509
62739
31376
92311
2317
80223
27133
5613
7600
62421
63667
92973
46295
21671
21874
92511
79403
1061
99877
38395
22395
88214
28997
3673
99811
38898
99399
50458
42872
30876
40751
2246
57259
13746
99737
29747
2178
14915
47675
72731
85656
91623
97285
47520
6077
88998
68768
70755
24022
86258
23016
295
41002
30848
26965
21662
1245
14283
32309
26938
89742
12964
14080
84262
29911
78715
82351
34260
43533
8119
27445
36643
11250
50027
71111
11591
17233
25975
75769
14456
40371
87837
48836
75927
98641
97100
26918
59381
11838
57971
32754
92310
12146
67547
94969
`

// Factorization helpers from 293C.go
func modMul(a, b, mod int64) int64 {
	t := new(big.Int).Mul(big.NewInt(a), big.NewInt(b))
	t.Mod(t, big.NewInt(mod))
	return t.Int64()
}

func modPow(a, d, mod int64) int64 {
	result := int64(1)
	base := a % mod
	for d > 0 {
		if d&1 == 1 {
			result = modMul(result, base, mod)
		}
		base = modMul(base, base, mod)
		d >>= 1
	}
	return result
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func isPrime(n int64) bool {
	if n < 2 {
		return false
	}
	smallPrimes := []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
	for _, p := range smallPrimes {
		if n == p {
			return true
		}
		if n%p == 0 {
			return false
		}
	}
	d := n - 1
	s := 0
	for d&1 == 0 {
		d >>= 1
		s++
	}
	bases := []int64{2, 325, 9375, 28178, 450775, 9780504, 1795265022}
	for _, a := range bases {
		if a%n == 0 {
			continue
		}
		x := modPow(a, d, n)
		if x == 1 || x == n-1 {
			continue
		}
		skip := false
		for r := 1; r < s; r++ {
			x = modMul(x, x, n)
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

func rho(n int64) int64 {
	if n%2 == 0 {
		return 2
	}
	for {
		x := rand.Int63n(n-2) + 2
		y := x
		c := rand.Int63n(n-1) + 1
		d := int64(1)
		for d == 1 {
			x = (modMul(x, x, n) + c) % n
			y = (modMul(y, y, n) + c) % n
			y = (modMul(y, y, n) + c) % n
			d = gcd(abs(x-y), n)
			if d == n {
				break
			}
		}
		if d > 1 && d < n {
			return d
		}
	}
}

func factor(n int64, fs map[int64]int) {
	if n <= 1 {
		return
	}
	if isPrime(n) {
		fs[n]++
	} else {
		d := rho(n)
		factor(d, fs)
		factor(n/d, fs)
	}
}

func divisorsFromFactors(fs map[int64]int) []int64 {
	divs := []int64{1}
	for p, e := range fs {
		sz := len(divs)
		mul := int64(1)
		for i := 1; i <= e; i++ {
			mul *= p
			for j := 0; j < sz; j++ {
				divs = append(divs, divs[j]*mul)
			}
		}
	}
	return divs
}

// solve293C mirrors 293C.go to compute expected count for n.
func solve293C(n int64) int64 {
	if n%3 != 0 {
		return 0
	}
	m := n / 3
	fs := make(map[int64]int)
	factor(m, fs)
	divs := divisorsFromFactors(fs)
	var cnt int64
	for _, x := range divs {
		if x > m {
			continue
		}
		t1 := m / x
		for _, y := range divs {
			if y > t1 {
				continue
			}
			if t1%y != 0 {
				continue
			}
			z := t1 / y
			if (x+y+z)&1 != 0 {
				continue
			}
			a := (x + z - y) / 2
			b := (x + y - z) / 2
			c := (y + z - x) / 2
			if a > 0 && b > 0 && c > 0 {
				cnt++
			}
		}
	}
	return cnt
}

func loadTestcases() ([]int64, error) {
	fields := strings.Fields(rawTestcasesC)
	cases := make([]int64, 0, len(fields))
	for idx, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %d (%q): %w", idx+1, f, err)
		}
		cases = append(cases, v)
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	rand.Seed(time.Now().UnixNano())
	bin := os.Args[1]

	testcases, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, n := range testcases {
		expect := solve293C(n)
		input := fmt.Sprintf("%d\n", n)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("case %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(out.String())
		got, err := strconv.ParseInt(gotStr, 10, 64)
		if err != nil {
			fmt.Printf("case %d: failed to parse output %q\n", idx+1, gotStr)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %d got %d\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesG = `100
19
3
16
9
2
1
5
19
16
12
11
1
9
16
7
14
18
18
4
7
19
18
9
20
3
14
11
3
12
9
4
6
3
15
2
8
5
4
12
20
15
1
11
1
5
18
14
16
2
19
1
8
1
11
1
9
8
10
11
7
11
16
5
1
8
20
12
13
20
13
14
13
12
17
14
9
1
16
9
5
17
10
12
13
7
6
2
15
4
7
5
14
9
14
10
7
14
13
4
12
12`

const mod int64 = 1e9 + 7

func modPow(a, b, m int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % m
		}
		a = a * a % m
		b >>= 1
	}
	return res
}

func prepareFact(n int) ([]int64, []int64) {
	fact := make([]int64, n+1)
	inv := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	inv[n] = modPow(fact[n], mod-2, mod)
	for i := n; i > 0; i-- {
		inv[i-1] = inv[i] * int64(i) % mod
	}
	return fact, inv
}

func C(n, r int64, fact, inv []int64) int64 {
	if r < 0 || r > n {
		return 0
	}
	return fact[n] * inv[r] % mod * inv[n-r] % mod
}

func solve(n int, fact, inv []int64) int64 {
	res := int64(0)
	half := (n - 1) / 2
	for s := 0; s <= half; s++ {
		for x := s + 1; x <= 2*s+1 && x <= n; x++ {
			ways := C(int64(x-1), int64(s), fact, inv) * C(int64(n-x), int64(n-2*s-1), fact, inv) % mod
			res = (res + int64(x)*ways) % mod
		}
	}
	for s := (n + 1) / 2; s <= n; s++ {
		ways := C(int64(n), int64(s), fact, inv)
		mex := int64(2*s + 1)
		res = (res + mex*ways) % mod
	}
	return res % mod
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierG /path/to/binary")
		os.Exit(1)
	}
	fields := strings.Fields(testcasesG)
	if len(fields) == 0 {
		fmt.Println("no testcases")
		os.Exit(1)
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		fmt.Println("parse error:", err)
		os.Exit(1)
	}
	if len(fields) < t+1 {
		fmt.Printf("expected %d cases, found %d\n", t, len(fields)-1)
		os.Exit(1)
	}
	maxN := 5000
	fact, inv := prepareFact(maxN)
	expected := make([]int64, t)
	for i := 0; i < t; i++ {
		n, err := strconv.Atoi(fields[i+1])
		if err != nil {
			fmt.Println("parse error:", err)
			os.Exit(1)
		}
		if n > maxN {
			fmt.Printf("n too large: %d\n", n)
			os.Exit(1)
		}
		expected[i] = solve(n, fact, inv)
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(t))
	sb.WriteByte('\n')
	for i := 0; i < t; i++ {
		sb.WriteString(fields[i+1])
		sb.WriteByte('\n')
	}
	input := sb.String()
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outFields := strings.Fields(string(out))
	if len(outFields) != t {
		fmt.Printf("expected %d outputs, got %d\n", t, len(outFields))
		os.Exit(1)
	}
	for i := 0; i < t; i++ {
		got, _ := strconv.ParseInt(outFields[i], 10, 64)
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %d got %d\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}

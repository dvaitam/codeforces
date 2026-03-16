package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Matrix [2][2]uint64

func mul(a, b Matrix, mod uint64) Matrix {
	return Matrix{{(a[0][0]*b[0][0] + a[0][1]*b[1][0]) % mod, (a[0][0]*b[0][1] + a[0][1]*b[1][1]) % mod},
		{(a[1][0]*b[0][0] + a[1][1]*b[1][0]) % mod, (a[1][0]*b[0][1] + a[1][1]*b[1][1]) % mod}}
}

func matrixPow(mat Matrix, exp, mod uint64) Matrix {
	res := Matrix{{1, 0}, {0, 1}}
	for exp > 0 {
		if exp&1 == 1 {
			res = mul(res, mat, mod)
		}
		mat = mul(mat, mat, mod)
		exp >>= 1
	}
	return res
}

func fib(n, mod uint64) uint64 {
	if n == 0 {
		return 0
	}
	M := Matrix{{1, 1}, {1, 0}}
	P := matrixPow(M, n-1, mod)
	return P[0][0]
}

func powMod(base, exp, mod uint64) uint64 {
	res := uint64(1) % mod
	base %= mod
	for exp > 0 {
		if exp&1 == 1 {
			res = (res * base) % mod
		}
		base = (base * base) % mod
		exp >>= 1
	}
	return res
}

func expected(n, k, l, m uint64) uint64 {
	if l < 64 && k>>l != 0 {
		return 0
	}
	c1 := uint64(0)
	kk := k
	for kk > 0 {
		c1 += kk & 1
		kk >>= 1
	}
	c0 := l - c1
	A := fib(n+2, m)
	T := powMod(2, n, m)
	var B uint64
	if T >= A {
		B = T - A
	} else {
		B = T + m - A
	}
	r0 := powMod(A, c0, m)
	r1 := powMod(B, c1, m)
	return (r0 * r1) % m
}

const testcasesDRaw = `7 13 1 267
8 12 3 993
8 11 5 915
4 4 3 145
2 8 5 724
3 9 1 749
2 10 4 575
2 11 4 325
4 15 4 888
5 1 5 939
1 2 6 862
7 0 5 507
6 7 6 335
2 6 5 229
4 4 5 460
2 2 3 898
8 3 3 566
5 3 5 342
4 9 4 95
7 10 5 249
5 5 2 843
3 1 5 674
5 15 1 93
3 4 1 864
2 12 6 539
5 7 2 918
7 8 4 506
6 2 3 629
2 15 5 647
6 6 2 18
5 3 6 227
6 5 3 438
1 3 2 877
4 1 5 651
2 0 1 652
4 3 4 95
6 3 1 622
1 6 2 737
2 15 2 746
1 0 5 437
2 8 1 228
2 9 3 448
3 1 5 480
1 3 6 402
4 8 3 928
8 5 6 690
4 1 6 164
3 10 5 258
2 14 6 181
1 15 6 421
5 11 4 859
5 4 5 709
1 14 6 82
6 1 5 289
3 7 4 362
5 11 5 971
3 9 4 768
7 2 1 610
4 10 2 247
4 14 4 729
7 1 4 894
7 1 2 458
2 8 6 163
8 15 5 620
1 1 4 335
5 14 1 830
7 6 5 988
2 4 1 413
7 10 1 220
1 0 6 543
2 6 1 624
4 9 3 707
3 3 4 876
7 2 1 283
8 3 3 138
6 3 2 287
1 1 1 212
5 10 3 962
1 15 6 661
8 13 3 894
3 6 4 603
5 0 2 156
5 10 3 810
6 2 3 800
1 1 3 169
3 9 3 406
3 9 1 491
4 1 3 185
2 9 4 857
6 9 4 113
2 15 4 347
6 3 4 120
8 13 1 311
6 4 2 643
7 2 1 829
2 6 6 228
1 12 1 102
7 9 4 944
8 6 4 87
6 7 3 601
`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesDRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) != 4 {
			fmt.Printf("case %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.ParseUint(parts[0], 10, 64)
		k, _ := strconv.ParseUint(parts[1], 10, 64)
		l, _ := strconv.ParseUint(parts[2], 10, 64)
		m, _ := strconv.ParseUint(parts[3], 10, 64)
		exp := expected(n, k, l, m)
		input := fmt.Sprintf("%d %d %d %d\n", n, k, l, m)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(out.String())
		got, err := strconv.ParseUint(gotStr, 10, 64)
		if err != nil || got != exp {
			fmt.Printf("case %d failed: expected %d got %s\n", idx, exp, gotStr)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

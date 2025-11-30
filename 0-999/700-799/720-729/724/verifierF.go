package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesF.txt.
const embeddedTestcasesF = `2 2 1000000007
5 3 1000000007
4 3 1000000007
2 2 1000000007
10 5 1000000007
1 2 1000000007
2 3 1000000007
4 2 1000000007
9 3 1000000007
9 5 1000000007
4 5 1000000007
10 4 1000000007
1 3 1000000007
7 4 1000000007
5 3 1000000007
4 4 1000000007
2 2 1000000007
7 2 1000000007
6 4 1000000007
10 4 1000000007
1 5 1000000007
9 2 1000000007
7 2 1000000007
9 4 1000000007
10 4 1000000007
10 3 1000000007
2 2 1000000007
4 4 1000000007
2 3 1000000007
2 5 1000000007
5 5 1000000007
6 3 1000000007
6 4 1000000007
4 4 1000000007
2 3 1000000007
9 3 1000000007
3 5 1000000007
7 4 1000000007
9 3 1000000007
6 2 1000000007
4 2 1000000007
6 5 1000000007
5 2 1000000007
4 4 1000000007
4 5 1000000007
7 5 1000000007
3 4 1000000007
3 3 1000000007
9 4 1000000007
10 5 1000000007
10 5 1000000007
6 3 1000000007
3 5 1000000007
2 2 1000000007
2 3 1000000007
3 5 1000000007
10 2 1000000007
7 5 1000000007
10 5 1000000007
9 4 1000000007
9 2 1000000007
2 4 1000000007
6 2 1000000007
5 5 1000000007
3 5 1000000007
1 4 1000000007
9 3 1000000007
9 2 1000000007
5 3 1000000007
3 4 1000000007
3 2 1000000007
10 4 1000000007
8 2 1000000007
2 4 1000000007
5 3 1000000007
1 3 1000000007
10 2 1000000007
2 5 1000000007
2 3 1000000007
3 5 1000000007
9 3 1000000007
5 5 1000000007
4 3 1000000007
5 5 1000000007
6 5 1000000007
9 5 1000000007
2 3 1000000007
4 2 1000000007
6 2 1000000007
10 3 1000000007
10 3 1000000007
1 2 1000000007
1 3 1000000007
2 2 1000000007
6 2 1000000007
9 3 1000000007
5 5 1000000007
4 3 1000000007
10 5 1000000007
4 5 1000000007`

func modPow(a, e, m int64) int64 {
	res := int64(1)
	a %= m
	for e > 0 {
		if e&1 == 1 {
			res = res * a % m
		}
		a = a * a % m
		e >>= 1
	}
	return res
}

func modInv(a, m int64) int64 {
	return modPow(a, m-2, m)
}

func combSeq(x int64, t int, invFact []int64, MOD int64) int64 {
	num := int64(1)
	for i := 0; i < t; i++ {
		num = num * ((x + int64(i)) % MOD) % MOD
	}
	return num * invFact[t] % MOD
}

func solve724F(n, d int, mod int64) int64 {
	MOD := mod
	maxC := d
	invFact := make([]int64, maxC+1)
	fact := make([]int64, maxC+1)
	fact[0] = 1
	for i := 1; i <= maxC; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[maxC] = modInv(fact[maxC], MOD)
	for i := maxC; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}
	inv2 := (MOD + 1) / 2

	N := n
	H := make([]int64, N+1)
	R := make([]int64, N+1)
	cH := d - 1
	dpH := make([][]int64, cH+1)
	for i := range dpH {
		dpH[i] = make([]int64, N+1)
	}
	cR := d
	dpR := make([][]int64, cR+1)
	for i := range dpR {
		dpR[i] = make([]int64, N+1)
	}
	H[1] = 1
	dpH[0][0] = 1
	dpR[0][0] = 1
	for t := 1; t <= cH; t++ {
		for k := cH - t; k >= 0; k-- {
			for s := 0; s+t <= N; s++ {
				dpH[k+t][s+t] = (dpH[k+t][s+t] + dpH[k][s]) % MOD
			}
		}
	}
	for t := 1; t <= cR; t++ {
		for k := cR - t; k >= 0; k-- {
			for s := 0; s+t <= N; s++ {
				dpR[k+t][s+t] = (dpR[k+t][s+t] + dpR[k][s]) % MOD
			}
		}
	}
	R[1] = (H[1] + dpR[cR][0]) % MOD

	for m := 2; m <= N; m++ {
		if m-1 >= 0 {
			H[m] = dpH[cH][m-1]
		}
		var cntR int64
		if m-1 >= 0 {
			cntR = dpR[cR][m-1]
		}
		R[m] = (H[m] + cntR) % MOD

		waysH := make([]int64, cH+1)
		for t := 1; t <= cH; t++ {
			waysH[t] = combSeq(H[m], t, invFact, MOD)
		}
		waysR := make([]int64, cR+1)
		for t := 1; t <= cR; t++ {
			waysR[t] = combSeq(H[m], t, invFact, MOD)
		}
		for t := 1; t <= cH; t++ {
			w := waysH[t]
			if w == 0 {
				continue
			}
			for k := cH - t; k >= 0; k-- {
				for s := 0; s+m*t <= N; s++ {
					dpH[k+t][s+m*t] = (dpH[k+t][s+m*t] + dpH[k][s]*w) % MOD
				}
			}
		}
		for t := 1; t <= cR; t++ {
			w := waysR[t]
			if w == 0 {
				continue
			}
			for k := cR - t; k >= 0; k-- {
				for s := 0; s+m*t <= N; s++ {
					dpR[k+t][s+m*t] = (dpR[k+t][s+m*t] + dpR[k][s]*w) % MOD
				}
			}
		}
	}

	var conv int64
	for i := 1; i < n; i++ {
		conv = (conv + R[i]*R[n-i]) % MOD
	}
	var mid int64
	if n%2 == 0 {
		mid = R[n/2]
	}
	ans := (R[n] - (conv+mid)%MOD*inv2%MOD) % MOD
	if ans < 0 {
		ans += MOD
	}
	return ans
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcasesF), "\n")
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 3 {
			fmt.Fprintf(os.Stderr, "case %d malformed\n", idx+1)
			os.Exit(1)
		}
		n, err1 := strconv.Atoi(fields[0])
		d, err2 := strconv.Atoi(fields[1])
		modVal, err3 := strconv.ParseInt(fields[2], 10, 64)
		if err1 != nil || err2 != nil || err3 != nil {
			fmt.Fprintf(os.Stderr, "case %d: parse error\n", idx+1)
			os.Exit(1)
		}

		want := strconv.FormatInt(solve724F(n, d, modVal), 10)
		input := fmt.Sprintf("%d %d %d\n", n, d, modVal)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}

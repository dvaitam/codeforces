package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesD.txt so the verifier is self-contained.
const testcasesD = `865 777
912 42
266 249
415 156
992 367
598 517
143 36
774 634
819 546
723 151
318 38
921 701
339 287
104 56
324 105
990 489
454 267
267 8
96 52
728 641
2 2
848 250
748 721
892 196
940 228
245 37
823 459
94 41
897 502
112 71
299 281
341 105
987 618
561 456
94 50
325 124
298 97
842 34
628 488
71 17
898 40
863 717
946 554
700 538
283 121
870 696
604 594
282 253
677 366
85 79
119 76
646 195
249 188
278 113
381 171
437 32
104 19
876 225
47 41
932 617
697 28
128 31
401 190
854 38
621 200
985 736
127 27
745 696
24 14
636 267
72 10
663 359
447 32
516 41
611 401
205 92
927 482
859 174
715 209
990 60
808 163
866 351
543 121
612 180
14 11
420 261
940 665
366 337
257 7
469 41
344 279
288 123
781 361
625 368
605 318
398 213
849 83
2 1
716 164
246 164`

const MOD = 1000000007
const MAXN = 1000
const Kmax = 44

func add(a, b int) int {
	a += b
	if a >= MOD {
		a -= MOD
	}
	return a
}

func mul(a, b int) int {
	return int((int64(a) * int64(b)) % MOD)
}

var f [MAXN + 1][Kmax + 1]int

func precompute() {
	var D [Kmax + 1][MAXN + 1]int
	D[0][0] = 1
	for d := 0; d <= MAXN; d++ {
		for k := Kmax; k >= 1; k-- {
			for s := d; s <= MAXN; s++ {
				if D[k-1][s-d] != 0 {
					D[k][s] = add(D[k][s], D[k-1][s-d])
				}
			}
		}
	}
	fact := make([]int, MAXN+1)
	invfact := make([]int, MAXN+1)
	fact[0] = 1
	for i := 1; i <= MAXN; i++ {
		fact[i] = mul(fact[i-1], i)
	}
	invfact[MAXN] = modInv(fact[MAXN])
	for i := MAXN; i > 0; i-- {
		invfact[i-1] = mul(invfact[i], i)
	}
	comb := func(n, k int) int {
		if n < 0 || k < 0 || n < k {
			return 0
		}
		return mul(fact[n], mul(invfact[k], invfact[n-k]))
	}
	factk := make([]int, Kmax+1)
	factk[0] = 1
	for i := 1; i <= Kmax; i++ {
		factk[i] = mul(factk[i-1], i)
	}
	for k := 1; k <= Kmax; k++ {
		minS := k * (k - 1) / 2
		for s := minS; s <= MAXN; s++ {
			cnt := D[k][s]
			if cnt == 0 {
				continue
			}
			coef := mul(factk[k], cnt)
			for n := k + s; n <= MAXN; n++ {
				f[n][k] = add(f[n][k], mul(coef, comb(n-s, k)))
			}
		}
	}
}

func modInv(a int) int { return modPow(a, MOD-2) }

func modPow(a, e int) int {
	res := 1
	base := a % MOD
	for e > 0 {
		if e&1 != 0 {
			res = mul(res, base)
		}
		base = mul(base, base)
		e >>= 1
	}
	return res
}

func solve(n, k int) int {
	if k < 0 || k > Kmax || k > n {
		return 0
	}
	return f[n][k]
}

type testCase struct {
	n int
	k int
}

func parseCases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesD), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d: expected 2 numbers", idx+1)
		}
		n, err1 := strconv.Atoi(fields[0])
		k, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("line %d: bad integers", idx+1)
		}
		cases = append(cases, testCase{n: n, k: k})
	}
	return cases, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	precompute()

	cases, err := parseCases()
	if err != nil {
		fmt.Println("failed to load testcases:", err)
		os.Exit(1)
	}

	// build bulk input with leading t
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
	}

	gotStr, err := runCandidate(bin, sb.String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	outFields := strings.Fields(gotStr)
	if len(outFields) != len(cases) {
		fmt.Printf("expected %d outputs, got %d\n", len(cases), len(outFields))
		os.Exit(1)
	}
	for idx, tc := range cases {
		exp := solve(tc.n, tc.k)
		got, err := strconv.Atoi(outFields[idx])
		if err != nil {
			fmt.Printf("case %d: bad output\n", idx+1)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed: expected %d got %d\n", idx+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

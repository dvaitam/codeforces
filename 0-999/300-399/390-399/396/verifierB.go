package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesB.txt.
const testcasesBData = `
100
139
584
869
823
784
66
263
122
509
781
462
485
669
390
809
216
98
501
31
916
857
401
445
624
782
787
4
714
458
274
740
823
236
607
969
106
925
327
33
24
28
667
556
11
963
904
392
704
223
994
434
745
31
542
229
784
450
963
509
568
240
355
238
695
226
781
472
977
298
950
24
428
859
940
571
946
659
104
192
646
743
882
305
125
762
342
919
740
998
730
514
960
992
434
521
851
934
688
196
312
`

// ---------- reference solution (from 396B.go) ----------

func modMul(a, b, mod uint64) uint64 { return (a * b) % mod }

func modPow(a, d, mod uint64) uint64 {
	res := uint64(1)
	a %= mod
	for d > 0 {
		if d&1 == 1 {
			res = modMul(res, a, mod)
		}
		a = modMul(a, a, mod)
		d >>= 1
	}
	return res
}

func isPrime(n uint64) bool {
	if n < 2 {
		return false
	}
	small := []uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
	for _, p := range small {
		if n == p {
			return true
		}
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
	bases := []uint64{2, 7, 61}
	for _, a := range bases {
		if a >= n {
			continue
		}
		x := modPow(a, d, n)
		if x == 1 || x == n-1 {
			continue
		}
		composite := true
		for r := 1; r < s; r++ {
			x = modMul(x, x, n)
			if x == n-1 {
				composite = false
				break
			}
		}
		if composite {
			return false
		}
	}
	return true
}

func gcd(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solve(n uint64) string {
	pk := n
	for pk >= 2 && !isPrime(pk) {
		pk--
	}
	pk1 := n + 1
	for !isPrime(pk1) {
		pk1++
	}
	R := int64(n-pk+1) - int64(pk1)
	D := pk * pk1
	pnum := int64(D) + 2*R
	qden := uint64(2) * D
	qnum := uint64(pnum)
	g := gcd(qnum, qden)
	qnum /= g
	qden /= g
	return fmt.Sprintf("%d/%d", qnum, qden)
}

// ---------- runner ----------

type testCase struct {
	n uint64
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesBData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no data")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("parse T: %w", err)
	}
	if len(fields)-1 != t {
		return nil, fmt.Errorf("expected %d cases, got %d", t, len(fields)-1)
	}
	cases := make([]testCase, t)
	for i := 0; i < t; i++ {
		val, err := strconv.ParseUint(fields[1+i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("case %d: %w", i+1, err)
		}
		cases[i] = testCase{n: val}
	}
	return cases, nil
}

func runCandidate(bin string, tc testCase) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(fmt.Sprintf("1\n%d\n", tc.n))
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		expect := solve(tc.n)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

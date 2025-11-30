package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	c int
}

// Embedded testcases from testcasesG.txt.
const testcaseData = `
45
811
961
732
932
144
328
959
684
955
764
368
463
766
196
915
410
475
862
823
88
880
404
375
11
292
217
599
365
384
855
984
641
8
98
457
816
863
778
416
315
717
178
307
753
990
961
668
225
763
338
934
341
217
991
42
53
983
20
725
187
616
474
333
990
731
38
273
817
608
539
925
912
279
988
68
736
656
256
904
933
886
13
134
770
424
356
277
828
905
872
511
40
396
114
313
709
430
770
255
`

func parseTestcases() ([]testCase, int, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	res := make([]testCase, 0, len(lines))
	maxC := 0
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		v, err := strconv.Atoi(line)
		if err != nil {
			return nil, 0, fmt.Errorf("case %d bad integer: %v", i+1, err)
		}
		if v > maxC {
			maxC = v
		}
		res = append(res, testCase{c: v})
	}
	if len(res) == 0 {
		return nil, 0, fmt.Errorf("no test data")
	}
	return res, maxC, nil
}

// precompute mirrors 1512G.go but with configurable limit.
func precompute(limit int) []int32 {
	ans := make([]int32, limit+1)
	lp := make([]int32, limit+1)
	pw := make([]int32, limit+1)
	sumP := make([]int32, limit+1)
	sigma := make([]int32, limit+1)
	primes := make([]int32, 0)

	pw[1] = 1
	sumP[1] = 1
	sigma[1] = 1
	ans[1] = 1

	for i := 2; i <= limit; i++ {
		if lp[i] == 0 {
			lp[i] = int32(i)
			primes = append(primes, int32(i))
			pw[i] = int32(i)
			sumP[i] = 1 + int32(i)
			sigma[i] = sumP[i]
		}
		li := lp[i]
		for _, p32 := range primes {
			p := int(p32)
			if p > int(li) || i*p > limit {
				break
			}
			idx := i * p
			lp[idx] = p32
			if p32 == li {
				pw[idx] = pw[i] * p32
				sumP[idx] = sumP[i] + pw[idx]
				sigma[idx] = sigma[i/int(pw[i])] * sumP[idx]
			} else {
				pw[idx] = p32
				sumP[idx] = 1 + p32
				sigma[idx] = sigma[i] * sumP[idx]
			}
		}
		s := sigma[i]
		if int(s) <= limit && ans[s] == 0 {
			ans[s] = int32(i)
		}
	}
	return ans
}

func solve(ans []int32, c int) string {
	if c < len(ans) && ans[c] != 0 {
		return strconv.Itoa(int(ans[c]))
	}
	return "-1"
}

func runCandidate(bin string, tests []testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d\n", tc.c))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, maxC, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	limit := maxC
	if limit < 1 {
		limit = 1
	}
	ans := precompute(limit)

	got, err := runCandidate(bin, tests)
	if err != nil {
		fmt.Printf("failed: %v\n", err)
		os.Exit(1)
	}
	outputs := strings.Fields(got)
	if len(outputs) != len(tests) {
		fmt.Printf("expected %d outputs got %d\n", len(tests), len(outputs))
		os.Exit(1)
	}
	for i, tc := range tests {
		expect := solve(ans, tc.c)
		if outputs[i] != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, expect, outputs[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

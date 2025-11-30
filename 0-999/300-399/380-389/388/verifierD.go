package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const mod = 1000000007

func solve(k int64) int64 {
	var L int
	for i := 62; i >= 0; i-- {
		if (k>>i)&1 == 1 {
			L = i
			break
		}
	}
	maxn := L + 2
	pow2 := make([]int64, maxn)
	pow2[0] = 1
	for i := 1; i < maxn; i++ {
		pow2[i] = (pow2[i-1] * 2) % mod
	}

	dp := make([][2]int64, maxn)
	dp[0][1] = 1

	for pos := L; pos >= 0; pos-- {
		bit := (k >> pos) & 1
		next := make([][2]int64, maxn)
		for t := 0; t < maxn; t++ {
			for tight := 0; tight <= 1; tight++ {
				val := dp[t][tight]
				if val == 0 {
					continue
				}
				if !(tight == 1 && bit == 0) {
					nt := 0
					if tight == 1 && bit == 1 {
						nt = 1
					}
					if t+1 < maxn {
						next[t+1][nt] = (next[t+1][nt] + val) % mod
					}
				}
				mult0 := int64(1)
				if t > 0 {
					mult0 = pow2[t-1]
				}
				nt0 := 0
				if tight == 1 {
					if bit == 0 {
						nt0 = 1
					} else {
						nt0 = 0
					}
				}
				next[t][nt0] = (next[t][nt0] + val*mult0) % mod

				if t > 0 {
					mult1 := pow2[t-1]
					if !(tight == 1 && bit == 0) {
						nt1 := 0
						if tight == 1 && bit == 1 {
							nt1 = 1
						}
						next[t][nt1] = (next[t][nt1] + val*mult1) % mod
					}
				}
			}
		}
		dp = next
	}

	var ans int64
	for t := 0; t < maxn; t++ {
		ans = (ans + dp[t][0] + dp[t][1]) % mod
	}
	return ans
}

var testcasesRaw = `887063
172419
431143
475296
786131
610357
675706
934975
756089
1075
207685
48672
593975
650507
415778
77279
990454
739578
551909
421816
382967
214891
248942
597592
942707
106704
992938
527347
771474
936261
631620
141160
122784
106256
567150
858864
612456
771775
825038
237307
75827
296451
965313
987379
524309
206001
170203
66179
575075
567542
49288
180270
405136
907599
99285
366472
672614
684339
753915
683462
107144
587802
370366
10361
726339
462667
225445
768233
30401
291173
494218
576151
271521
823774
333790
592194
134926
350447
117025
816961
174457
493419
304386
99470
36132
976097
353073
159920
178185
306152
406829
478625
360866
241385
100288
42054
160550
372321
699392
111478
626385
865215
811720
334809`

func parseTestcases() ([]int64, error) {
	lines := strings.Fields(testcasesRaw)
	vals := make([]int64, 0, len(lines))
	for _, s := range lines {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		vals = append(vals, v)
	}
	return vals, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, k := range cases {
		expected := strconv.FormatInt(solve(k), 10)

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(strconv.FormatInt(k, 10) + "\n")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx+1, err, string(out))
			os.Exit(1)
		}

		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected: %s\n got: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(cases))
}

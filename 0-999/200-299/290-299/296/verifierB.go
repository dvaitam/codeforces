package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const mod = 1000000007

// Embedded testcases (previously in testcasesB.txt) to keep verifier self contained.
const rawTestcasesB = `
2 91 41
4 77?6 3170
4 6907 4391
3 000 ?80
4 ?360 8377
5 353?3 74068
6 12?415 868?34
3 978 690
4 366? 258?
6 517?81 286570
4 0499 96?2
2 83 03
5 83685 9574?
5 90682 88360
4 5983 8675
4 5088 9957
5 03?28 92184
1 ? 1
1 0 7
1 4 3
3 192 541
2 24 82
6 4?4757 710465
4 3414 8396
1 3 0
4 2027 8?68
2 ?8 73
5 ?06?9 5??60
6 423041 144269
3 208 093
5 72980 63513
5 ?6937 1?648
4 0596 4023
3 925 634
6 1685?8 783101
2 22 83
3 598 455
3 143 972
5 81506 16225
1 9 9
4 1983 9145
3 981 741
1 4 0
5 ?0161 03396
2 17 2?
2 21 66
5 48475 13?50
1 0 4
6 957656 115971
3 398 7?5
3 283 433
3 141 71?
5 ?5364 05259
3 351 899
5 13303 61481
6 10?045 772185
1 8 ?
2 22 25
3 189 423
2 80 59
6 832468 20?341
6 768487 870652
3 70? 690
1 5 9
2 92 24
3 696 291
2 70 28
3 8?7 ??3
2 57 ?7
2 65 89
6 ?4?301 8?5283
3 448 527
5 11989 62246
2 90 7?
4 ?568 2808
1 4 ?
1 4 1
2 9? ?1
4 3666 2572
5 73169 861?4
3 368 038
4 900? 9343
2 42 83
3 494 ?72
5 57613 96341
1 1 9
6 084??2 185946
5 ?5850 17754
5 65?97 1?663
5 04?98 37986
6 4279?8 3580?6
5 66599 173??
3 ?06 ?2?
4 4219 0546
6 842747 278048
1 9 6
1 5 1
6 702821 6?4943
5 33541 18?57
5 8024? 84593
`

func parseTestcases() ([]struct {
	n int
	s string
	t string
}, error) {
	lines := strings.Split(rawTestcasesB, "\n")
	var cases []struct {
		n int
		s string
		t string
	}
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 3 {
			return nil, fmt.Errorf("line %d: expected 3 fields got %d", idx+1, len(fields))
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", idx+1, err)
		}
		if len(fields[1]) != n || len(fields[2]) != n {
			return nil, fmt.Errorf("line %d: length mismatch", idx+1)
		}
		cases = append(cases, struct {
			n int
			s string
			t string
		}{n: n, s: fields[1], t: fields[2]})
	}
	return cases, nil
}

func solve296B(n int, s, t string) int64 {
	tot := int64(1)
	noGT := int64(1)
	noLT := int64(1)
	eqOnly := int64(1)
	for i := 0; i < n; i++ {
		var eq, lt, gt int64
		si := s[i]
		ti := t[i]
		if si != '?' && ti != '?' {
			if si == ti {
				eq = 1
			} else if si < ti {
				lt = 1
			} else {
				gt = 1
			}
		} else if si == '?' && ti == '?' {
			eq = 10
			lt = 45
			gt = 45
		} else if si == '?' {
			d := int64(ti - '0')
			eq = 1
			lt = d
			gt = 9 - d
		} else {
			d := int64(si - '0')
			eq = 1
			lt = 9 - d
			gt = d
		}
		total := (eq + lt + gt) % mod
		tot = tot * total % mod
		noGT = noGT * ((eq + lt) % mod) % mod
		noLT = noLT * ((eq + gt) % mod) % mod
		eqOnly = eqOnly * (eq % mod) % mod
	}
	ans := (tot - noGT - noLT + eqOnly) % mod
	if ans < 0 {
		ans += mod
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		expect := solve296B(tc.n, tc.s, tc.t)
		input := fmt.Sprintf("%d\n%s\n%s\n", tc.n, tc.s, tc.t)
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

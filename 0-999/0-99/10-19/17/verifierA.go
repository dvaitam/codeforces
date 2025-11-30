package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `100
866 6
778 6
43 4
990 8
499 6
942 4
993 7
368 9
915 3
518 2
290 2
775 1
635 4
933 8
724 9
925 2
319 1
749 1
922 10
340 7
575 1
364 6
325 9
657 3
991 8
490 7
888 8
268 0
826 8
939 0
97 6
729 10
642 0
628 7
849 5
251 5
722 1
197 9
229 3
824 2
824 8
460 1
84 5
898 8
957 7
113 4
566 4
725 1
562 5
836 8
210 9
562 9
296 7
95 9
819 6
326 9
249 4
190 3
843 2
35 9
674 4
489 1
93 10
777 2
899 2
947 0
864 1
921 8
701 6
859 8
284 8
833 3
871 3
918 10
605 6
595 4
463 7
678 10
719 5
86 5
629 1
500 9
647 5
867 3
250 0
751 4
121 3
382 2
342 6
837 0
105 2
877 3
48 9
651 8
618 10
77 0
129 10
195 9
852 9
124 6`

type testCase struct {
	input    string
	expected string
}

func solve(n, k int) string {
	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}
	primes := make([]int, 0)
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	count := 0
	for i := 0; i+1 < len(primes); i++ {
		s := primes[i] + primes[i+1] + 1
		if s <= n && isPrime[s] {
			count++
		}
	}
	if count >= k {
		return "YES"
	}
	return "NO"
}

func loadCases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}
	pos++
	cases := make([]testCase, 0, t)
	for caseNum := 0; caseNum < t; caseNum++ {
		if pos+1 >= len(fields) {
			return nil, fmt.Errorf("case %d: missing n/k", caseNum+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", caseNum+1, err)
		}
		pos++
		k, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad k: %w", caseNum+1, err)
		}
		pos++
		input := fmt.Sprintf("%d %d\n", n, k)
		cases = append(cases, testCase{
			input:    input,
			expected: solve(n, k),
		})
	}
	return cases, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
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
		fmt.Println("usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

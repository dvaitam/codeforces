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
5
5 13 2 3 18
3
12 19 2
7
7 2 3 14 14 3 8
3
18 14 2
7
4 8 19 2 19 19 13
3
8 2 18
4
10 14 5 18
3
19 10 18
8
6 4 19 19 7 12 4 18
8
3 19 2 20 7 16 18 14
5
15 19 15 12 10
4
6 8 3 19
5
17 16 11 15 10
7
3 4 17 14 6 11 5
6
14 2 3 18 19 11
5
12 20 16 19 15
3
3 9 16
8
3 2 10 19 15 10 13 12
3
15 12 6
7
4 16 2 7 10 5 8
6
13 16 3 6 15 13
7
9 5 14 18 9 14 12
8
13 8 5 3 6 5 8 8
3
16 19 6
5
10 1 5 14 18
5
20 19 11 5 17
7
2 15 18 13 13 13 13
3
16 13 2
4
3 7 15 6
3
11 20 2
3
1 19 5
7
4 12 20 1 3 7 20
6
5 9 12 20 12 16
3
4 16 15
6
16 10 3 5 4 11
8
9 16 6 17 1 7 17 12
4
18 1 17 10
8
3 9 17 12 6 12 8 18
7
17 11 8 20 7 8 13
8
8 7 17 16 12 1 1 9
6
9 7 20 12 15 12
5
3 8 4 8 16
4
11 7 16 20
7
1 16 12 3 4 13 7
6
6 14 11 3 13 15
6
3 6 6 5 1 5
7
15 5 20 20 16 12 5
7
18 5 1 1 4 17 5
6
7 7 1 9 7 10
7
8 19 11 9 18 14 5
3
12 15 19
7
14 17 5 18 5 17 17
3
15 6 20
3
5 6 5
6
20 4 18 2 11 17
7
18 16 4 18 2 8 7
5
2 4 17 15 18
3
3 15 11
7
17 20 17 7 9 15 17
7
16 17 8 17 9 18 7
6
5 14 4 13 15 11
3
8 14 3
4
10 4 5 12
4
9 5 15 8
8
4 13 16 6 8 6 14 17
6
11 14 7 12 11 3
8
12 1 11 18 15 15 1 13
5
17 20 10 17 3
3
8 4 3
5
9 2 6 9 5
6
9 13 5 18 17 19
6
11 3 9 2 6 14
3
9 1 3
5
3 20 8 3 9
3
15 1 11
7
14 9 20 5 2 17 8
3
6 9 2
4
7 10 10 17
4
10 15 17 6
5
12 1 9 2 1
3
17 18 7
7
16 8 15 4 14 16 18
6
17 10 7 8 11 7
8
5 13 12 2 5 1 3 9
6
6 2 3 13 17 10
7
8 10 2 15 6 6 9
6
1 9 12 11 18 11
4
2 10 7 12
4
1 11 13 3
6
9 17 7 8 17 1
3
9 3 5
6
19 2 13 1 10 10
8
8 3 19 17 5 20 13 11
8
16 5 10 20 5 2 17 14
8
17 5 17 17 19 1 19 8
3
1 2 5
8
12 4 13 15 18 2 1 18
8
8 16 9 1 15 3 17 18
3
17 3 16
5
3 9 8 7 8`

var primes []int

func initPrimes(limit int) {
	sieve := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		if !sieve[i] {
			primes = append(primes, i)
			for j := i * i; j <= limit; j += i {
				sieve[j] = true
			}
		}
	}
}

func divisors(n int) []int {
	res := []int{1}
	x := n
	for _, p := range primes {
		if p*p > x {
			break
		}
		if x%p == 0 {
			cnt := 0
			for x%p == 0 {
				x /= p
				cnt++
			}
			base := res
			res = make([]int, 0, len(base)*(cnt+1))
			pow := 1
			for i := 0; i <= cnt; i++ {
				for _, d := range base {
					res = append(res, d*pow)
				}
				pow *= p
			}
		}
	}
	if x > 1 {
		base := res
		res = make([]int, 0, len(base)*2)
		for _, d := range base {
			res = append(res, d)
			res = append(res, d*x)
		}
	}
	return res
}

func solveCase(arr []int) string {
	freq := make(map[int]int)
	for _, v := range arr {
		freq[v]++
	}
	var ans int64
	for _, c := range freq {
		if c >= 3 {
			ans += int64(c) * int64(c-1) * int64(c-2)
		}
	}
	for y, cy := range freq {
		ds := divisors(y)
		for _, i := range ds {
			if i == y {
				continue
			}
			k64 := int64(y) * int64(y) / int64(i)
			if k64 > 1000000000 {
				continue
			}
			k := int(k64)
			ci, ok1 := freq[i]
			ck, ok2 := freq[k]
			if ok1 && ok2 {
				ans += int64(ci) * int64(cy) * int64(ck)
			}
		}
	}
	return fmt.Sprintf("%d", ans)
}

func join(arr []int) string {
	var sb strings.Builder
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	return sb.String()
}

type testCase struct {
	input    string
	expected string
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
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing n", caseNum+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", caseNum+1, err)
		}
		pos++
		if pos+n > len(fields) {
			return nil, fmt.Errorf("case %d: missing array values", caseNum+1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[pos+i])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad value: %w", caseNum+1, err)
			}
			arr[i] = v
		}
		pos += n
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: solveCase(arr),
		})
	}
	return cases, nil
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	initPrimes(31623)
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierG2 /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := runExe(os.Args[1], tc.input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

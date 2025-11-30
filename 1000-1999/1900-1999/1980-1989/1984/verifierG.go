package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesG = `6 4 2 5 3 6 1
13 11 2 8 9 5 12 6 7 3 1 13 4 10
13 10 3 12 9 2 5 4 7 11 8 6 13 1
5 5 4 1 3 2
16 1 5 3 14 2 4 12 11 10 16 9 7 13 8 15 6
12 1 2 11 3 7 4 10 6 9 12 5 8
20 1 14 20 4 2 18 3 7 10 19 17 12 13 15 5 16 9 6 11 8
7 4 5 6 2 1 3 7
6 2 6 4 3 1 5
8 6 4 1 5 2 8 7 3
6 2 4 6 5 3 1
5 2 3 4 5 1
6 4 5 2 3 1 6
10 7 10 6 2 8 5 4 1 9 3
5 4 5 2 1 3
20 5 18 16 20 11 14 6 17 4 8 3 12 19 7 13 9 2 15 10 1
9 1 6 8 2 3 5 4 7 9
13 2 5 4 12 8 13 6 3 9 1 11 7 10
10 3 9 7 10 1 6 5 4 8 2
7 6 4 7 2 5 1 3
13 13 8 12 2 6 4 10 3 1 9 5 7 11
12 5 9 6 8 7 11 10 4 3 1 12 2
17 4 6 3 11 2 8 5 15 10 1 9 14 16 13 12 17 7
16 3 14 15 8 4 13 7 16 11 12 5 6 10 2 9 1
11 2 3 9 5 7 6 4 10 8 1 11
5 4 5 2 1 3
7 6 5 3 1 7 4 2
17 17 1 14 8 4 3 11 9 10 12 16 2 7 6 13 5 15
5 3 2 1 4 5
14 1 2 12 10 4 7 11 9 13 5 14 3 8 6
17 11 1 4 16 13 15 3 12 14 2 7 9 5 10 17 6 8
7 5 3 4 2 1 7 6
12 7 4 10 2 1 9 5 3 8 6 11 12
8 8 4 1 3 7 2 5 6
19 6 14 2 9 4 18 16 11 3 7 10 1 17 15 19 5 12 13 8
18 5 3 14 1 10 6 8 7 13 9 4 2 18 15 17 11 12 16
19 19 17 14 3 9 5 8 11 2 15 7 16 10 12 18 13 1 4 6
20 11 17 12 7 15 13 18 6 2 4 8 3 10 9 19 16 14 1 5 20
15 13 14 11 10 12 4 7 3 6 2 1 15 9 5 8
20 11 18 20 3 2 4 19 7 14 12 6 16 8 9 17 1 5 10 15 13
20 19 1 20 4 8 14 12 2 17 7 13 15 16 6 5 11 9 3 10 18
7 5 6 7 4 2 1 3
13 5 8 11 9 13 10 3 4 6 7 12 2 1
10 1 3 9 6 7 2 4 5 10 8
10 5 1 9 7 3 6 10 8 4 2
9 1 2 3 6 9 8 7 5 4
13 4 5 2 1 3 11 8 7 6 13 9 10 12
17 4 12 15 10 9 1 11 7 14 3 2 13 8 17 5 6 16
7 5 2 3 1 4 7 6
12 4 11 1 12 3 2 6 8 5 9 7 10
8 4 6 1 5 8 2 7 3
16 15 7 11 16 5 8 3 14 1 10 12 13 4 6 2 9
7 5 1 6 2 3 7 4
15 2 9 5 3 6 14 8 15 1 4 7 10 12 13 11
20 4 19 7 16 11 5 18 12 3 10 17 2 1 8 9 20 14 6 15 13
19 16 18 5 11 10 15 13 8 6 2 7 14 19 3 17 12 4 9 1
5 5 3 4 1 2
5 2 1 4 3 5
9 9 3 4 7 5 2 1 8 6
19 14 1 4 5 15 8 12 18 3 16 9 17 11 10 6 13 19 2 7
14 4 8 2 5 1 10 9 14 11 12 13 3 7 6
14 2 10 12 14 1 11 6 3 7 4 5 8 13 9
11 2 10 5 11 8 9 6 3 1 7 4
12 12 8 7 11 5 2 9 4 1 6 10 3
14 13 11 10 8 1 12 14 7 2 9 5 6 3 4
13 10 3 13 12 2 8 9 11 5 6 7 4 1
19 6 1 2 3 14 9 7 18 11 8 17 4 19 5 10 16 13 12 15
10 1 9 3 2 10 5 7 4 8 6
11 4 10 6 3 11 8 7 1 2 5 9
13 9 5 10 3 1 7 4 2 8 6 11 12 13
11 8 1 9 10 4 5 6 2 7 11 3
20 3 15 18 1 20 13 2 9 7 11 14 6 12 17 4 19 10 5 8 16
9 8 3 2 5 6 4 1 9 7
14 12 3 8 6 11 7 5 2 10 9 1 13 14 4
17 5 4 1 3 11 7 12 10 15 17 13 14 16 6 2 9 8
20 12 13 16 15 8 17 6 14 20 18 3 10 1 9 4 11 7 19 5 2
6 4 1 2 6 5 3
16 5 1 15 4 3 8 7 9 12 13 14 6 10 2 16 11
5 4 5 2 3 1
16 8 1 3 5 6 13 4 10 7 15 9 12 16 11 14 2
5 2 1 4 3 5
14 8 9 11 6 5 13 3 2 7 12 10 4 14 1
11 6 4 10 8 9 3 5 11 7 1 2
18 9 11 13 2 7 17 4 16 15 10 1 6 12 8 3 14 18 5
16 15 11 9 6 2 10 4 8 14 7 13 3 12 16 1 5
12 5 2 12 9 8 3 4 10 7 1 11 6
13 1 3 13 7 5 2 10 11 9 6 12 8 4
16 3 16 1 6 8 12 5 10 15 7 13 4 11 9 2 14
19 1 17 5 14 2 10 16 9 8 7 15 13 11 4 3 6 19 18 12
17 13 3 2 5 14 16 8 4 12 17 7 1 15 6 11 9 10
15 11 13 1 14 10 5 4 8 9 2 6 15 3 7 12
16 13 14 2 10 15 9 4 7 6 11 16 12 3 8 5 1
12 1 7 2 8 5 4 9 11 3 6 10 12
7 4 6 1 3 7 2 5
16 16 1 5 7 3 13 2 8 12 14 10 11 9 4 15 6
5 1 2 3 5 4
10 10 6 9 7 2 1 5 8 3 4
5 3 4 2 1 5
10 4 2 10 6 9 3 8 1 7 5
20 16 10 12 1 9 19 7 17 5 2 11 8 14 20 18 6 15 3 4 13
`

type pair struct{ l, r int }

// solveOne replicates the reference solution logic for a single test case.
func solveOne(arr []int) string {
	var a, b [1005]int
	for i, v := range arr {
		a[i+1] = v
	}
	n := len(arr)
	ans := make([]pair, 0, n*n)

	var sb strings.Builder

	printAns := func() string {
		fmt.Fprintln(&sb, len(ans))
		for _, p := range ans {
			fmt.Fprintf(&sb, "%d %d\n", p.l, p.r)
		}
		return strings.TrimSpace(sb.String())
	}

	move := func(l, r, x int) {
		length := r - l + 1
		if x < 0 {
			x += length
		}
		for i := l; i <= r; i++ {
			idx := (i+x-l)%length + l
			b[idx] = a[i]
		}
		for i := l; i <= r; i++ {
			a[i] = b[i]
		}
	}

	// Already sorted
	flag := true
	for i := 1; i <= n; i++ {
		if a[i] != i {
			flag = false
			break
		}
	}
	if flag {
		fmt.Fprintln(&sb, n)
		fmt.Fprintln(&sb, 0)
		return strings.TrimSpace(sb.String())
	}

	// Simple rotation case
	flag = true
	for i := 1; i < n; i++ {
		diff := a[i+1] - a[i]
		if diff != 1 && diff != 1-n {
			flag = false
			break
		}
	}
	if flag {
		fmt.Fprintln(&sb, n-1)
		p := 0
		for i := 1; i <= n; i++ {
			if a[i] == 1 {
				p = i
				break
			}
		}
		fmt.Fprintln(&sb, p-1)
		for i := 1; i < p; i++ {
			fmt.Fprintln(&sb, 2, 1)
		}
		return strings.TrimSpace(sb.String())
	}

	inv := 0
	for i := 1; i < n; i++ {
		for j := i + 1; j <= n; j++ {
			if a[i] > a[j] {
				inv++
			}
		}
	}

	if inv&1 == 1 && (n-1)&1 == 1 {
		fmt.Fprintln(&sb, n-3)
		posn := 0
		for i := 1; i <= n; i++ {
			if a[i] == n {
				posn = i
				break
			}
		}
		if posn == 1 {
			ans = append(ans, pair{2, 1})
			move(1, n-2, -1)
			posn = n - 2
		} else if posn == 2 {
			ans = append(ans, pair{1, 2})
			move(1, n-2, 1)
			posn = 3
		}
		for i := posn; i < n; i++ {
			ans = append(ans, pair{3, 4})
		}
		move(3, n, n-posn)
		n--
		if a[1] == 1 {
			ans = append(ans, pair{1, 2})
			tmp := a[n-1]
			for i := n - 1; i >= 2; i-- {
				a[i] = a[i-1]
			}
			a[1] = tmp
		} else if a[n] == 1 {
			ans = append(ans, pair{2, 3})
			for i := n; i >= 3; i-- {
				a[i] = a[i-1]
			}
			a[2] = 1
		}
		for i := 2; i <= n-2; i++ {
			pos1, posi := 0, 0
			for j := 1; j <= n; j++ {
				if a[j] == 1 {
					pos1 = j
				}
				if a[j] == i {
					posi = j
				}
			}
			if posi == 1 {
				if pos1 == 2 {
					ans = append(ans, pair{1, 2})
					tmp := a[n-1]
					for j := n - 1; j >= 2; j-- {
						a[j] = a[j-1]
					}
					a[1] = tmp
					pos1++
					posi++
				} else {
					ans = append(ans, pair{2, 1})
					for j := 1; j <= n-2; j++ {
						a[j] = a[j+1]
					}
					a[n-1] = i
					pos1--
					posi = n - 1
				}
			}
			if posi != n {
				if posi < pos1 {
					for j := 2; j <= posi; j++ {
						ans = append(ans, pair{3, 2})
					}
					move(2, n, -(posi - 1))
					pos1 -= posi - 1
					posi = n
				} else {
					for j := posi; j < n; j++ {
						ans = append(ans, pair{2, 3})
					}
					move(2, n, n-posi)
					pos1 += n - posi
					posi = n
				}
			}
			for j := pos1 + i - 2; j < n-1; j++ {
				ans = append(ans, pair{1, 2})
			}
			move(1, n-1, n-1-(pos1+i-2))
			pos1 += n - 1 - (pos1 + i - 2)
			ans = append(ans, pair{3, 2})
			move(2, n, -1)
		}
		ans = append(ans, pair{2, 1})
		move(1, n-1, -1)
		if a[n-1] == n-1 {
			return printAns()
		}
		for i := 1; i <= (n-3)/2; i++ {
			ans = append(ans, pair{1, 2})
			ans = append(ans, pair{2, 3})
			ans = append(ans, pair{2, 3})
			ans = append(ans, pair{2, 1})
			ans = append(ans, pair{3, 2})
			ans = append(ans, pair{2, 1})
			ans = append(ans, pair{2, 3})
			ans = append(ans, pair{2, 3})
		}
		ans = append(ans, pair{2, 3})
		return printAns()
	}

	fmt.Fprintln(&sb, n-2)
	if a[1] == 1 {
		ans = append(ans, pair{1, 2})
		tmp := a[n-1]
		for i := n - 1; i >= 2; i-- {
			a[i] = a[i-1]
		}
		a[1] = tmp
	} else if a[n] == 1 {
		ans = append(ans, pair{2, 3})
		for i := n; i >= 3; i-- {
			a[i] = a[i-1]
		}
		a[2] = 1
	}
	for i := 2; i <= n-2; i++ {
		pos1, posi := 0, 0
		for j := 1; j <= n; j++ {
			if a[j] == 1 {
				pos1 = j
			}
			if a[j] == i {
				posi = j
			}
		}
		if posi == 1 {
			if pos1 == 2 {
				ans = append(ans, pair{1, 2})
				tmp := a[n-1]
				for j := n - 1; j >= 2; j-- {
					a[j] = a[j-1]
				}
				a[1] = tmp
				pos1++
				posi++
			} else {
				ans = append(ans, pair{2, 1})
				for j := 1; j <= n-2; j++ {
					a[j] = a[j+1]
				}
				a[n-1] = i
				pos1--
				posi = n - 1
			}
		}
		if posi != n {
			if posi < pos1 {
				for j := 2; j <= posi; j++ {
					ans = append(ans, pair{3, 2})
				}
				move(2, n, -(posi - 1))
				pos1 -= posi - 1
				posi = n
			} else {
				for j := posi; j < n; j++ {
					ans = append(ans, pair{2, 3})
				}
				move(2, n, n-posi)
				pos1 += n - posi
				posi = n
			}
		}
		for j := pos1 + i - 2; j < n-1; j++ {
			ans = append(ans, pair{1, 2})
		}
		move(1, n-1, n-1-(pos1+i-2))
		pos1 += n - 1 - (pos1 + i - 2)
		ans = append(ans, pair{3, 2})
		move(2, n, -1)
	}
	ans = append(ans, pair{2, 1})
	move(1, n-1, -1)
	if a[n-1] == n-1 {
		return printAns()
	}
	for i := 1; i <= (n-3)/2; i++ {
		ans = append(ans, pair{1, 2})
		ans = append(ans, pair{2, 3})
		ans = append(ans, pair{2, 3})
		ans = append(ans, pair{2, 1})
		ans = append(ans, pair{3, 2})
		ans = append(ans, pair{2, 1})
		ans = append(ans, pair{2, 3})
		ans = append(ans, pair{2, 3})
	}
	ans = append(ans, pair{2, 3})
	return printAns()
}

type testCase struct {
	arr []int
}

func parseTests() ([]testCase, error) {
	fields := strings.Fields(testcasesG)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no data")
	}
	pos := 0
	readInt := func() (int, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected EOF")
		}
		v, err := strconv.Atoi(fields[pos])
		pos++
		return v, err
	}
	var tests []testCase
	for pos < len(fields) {
		n, err := readInt()
		if err != nil {
			return nil, err
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := readInt()
			if err != nil {
				return nil, err
			}
			arr[i] = val
		}
		tests = append(tests, testCase{arr: arr})
	}
	return tests, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(len(tc.arr)))
		for _, v := range tc.arr {
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierG /path/to/binary")
		os.Exit(1)
	}

	tests, err := parseTests()
	if err != nil {
		fmt.Println("parse error:", err)
		os.Exit(1)
	}

	input := buildInput(tests)
	output, err := runCandidate(os.Args[1], input)
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}

	gotFields := strings.Fields(output)
	pos := 0
	for i, tc := range tests {
		if pos >= len(gotFields) {
			fmt.Printf("missing output for case %d\n", i+1)
			os.Exit(1)
		}
		// Candidate outputs variable lines; use reference to know expected token count.
		expTokens := strings.Fields(solveOne(tc.arr))
		if pos+len(expTokens) > len(gotFields) {
			fmt.Printf("not enough output for case %d\n", i+1)
			os.Exit(1)
		}
		for j := 0; j < len(expTokens); j++ {
			if gotFields[pos+j] != expTokens[j] {
				fmt.Printf("case %d failed at token %d\nexpected: %s\ngot: %s\n", i+1, j+1, expTokens[j], gotFields[pos+j])
				os.Exit(1)
			}
		}
		pos += len(expTokens)
	}
	if pos != len(gotFields) {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

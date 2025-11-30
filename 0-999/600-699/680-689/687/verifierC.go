package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesC.txt.
const embeddedTestcasesC = `7 19 11 12 17 17 1 9 1
2 22 5 17
3 29 15 14 8
6 44 3 13 13 5 4 14
2 18 16 2
7 95 11 17 20 20 12 15 2
1 18 18
7 57 20 18 5 3 9 12 19
6 45 4 16 11 11 7 2
1 18 18
3 23 19 4 16
7 32 4 9 14 19 17 14 3
4 4 4 4 9 5
1 3 3
8 64 17 19 1 9 18 9 16 15
2 24 14 10
2 19 9 10
8 29 4 14 12 13 1 14 15 3
2 10 10 12
2 23 9 14
2 16 16 2
8 58 9 16 4 19 16 9 7 8
2 28 16 12
5 8 4 1 1 1 1
7 23 12 8 2 1 5 4 1
6 13 8 4 1 5 1 1
3 33 9 6 18
6 62 17 8 17 20 4 7
8 8 1 1 1 1 1 1 1 1
3 21 3 14 4
6 67 5 12 9 15 16 10
8 49 11 3 7 7 10 10 1 3
2 21 14 7
7 41 10 13 10 1 6 1 10
3 44 1 13 11
4 14 8 1 3 2
8 31 10 9 5 3 2 1 1 8
5 9 2 2 1 2 2
1 4 4
8 45 1 6 1 11 12 8 3 3
6 17 6 2 1 1 1 6
8 39 8 2 3 10 3 3 7 3
2 6 5 1
2 24 1 23
4 13 1 7 1 4
5 10 1 2 4 2 1
8 30 8 10 1 9 1 1 6 7
3 16 10 1 5
6 11 1 3 1 3 3 1
7 14 4 2 2 2 4 1 1
6 34 3 7 5 9 1 9
2 3 1 2
8 44 7 2 4 6 3 5 3 14
2 19 16 3
2 14 6 8
8 27 1 10 2 7 3 1 1 2
3 22 17 1 4
4 20 4 5 7 4
4 43 14 2 13 14
3 6 5 1 1
5 12 3 4 1 2 2
4 21 5 5 5 6
4 20 5 4 4 7
2 8 2 6
4 15 2 1 1 11
5 27 10 2 10 5 5
4 8 2 5 1 4
7 29 1 3 1 8 8 2 6
8 43 5 1 3 7 8 7 1 11
4 9 1 4 1 3
7 36 8 2 4 2 5 1 14
4 36 11 5 10 10
6 32 13 8 5 1 3 2
4 12 4 1 6 1
7 15 1 1 1 4 5 1 2
8 22 2 9 1 7 1 1 1 1
7 20 5 1 3 4 2 3 2
1 3 3
7 11 1 3 3 3 1 1 1
1 7 7
8 21 5 6 2 1 1 2 2 2
4 12 5 1 3 3
8 33 11 6 1 12 1 1 1 12
5 25 1 8 1 1 3
7 26 5 2 6 4 3 1 5
2 30 1 29
2 10 1 9
1 13 13
6 46 16 11 1 1 9 8
1 16 16
1 1 1
4 13 10 1 1 1
4 3 1 1 1 1
3 8 7 1 1
5 21 10 2 1 1 7
5 14 4 1 4 4 1
3 20 10 10 1
6 37 5 10 7 6 4 5
8 50 8 7 1 2 2 1 10 19
7 25 1 1 1 1 5 1 15
8 54 18 4 6 3 7 7 1 8
1 7 7
7 32 2 3 5 1 11 9 1
6 52 1 7 1 20 1 22
7 30 9 10 8 1 1 1 1
3 11 2 2 7
2 8 3 5
6 40 15 3 14 2 2 4
8 41 4 7 8 1 6 4 1 10
8 57 1 14 11 9 6 1 3 12
8 20 4 4 2 1 1 3 1 3
1 13 13
1 19 19
7 77 12 9 1 15 10 12 18
5 25 1 1 1 8 14
2 17 6 11
4 25 4 1 7 13
8 33 1 1 1 1 1 1 1 26
7 14 2 2 2 2 2 2 2
6 30 10 10 1 2 5 2
7 68 1 7 1 15 1 5 38
7 41 3 9 9 3 13 1 3
5 16 1 1 1 1 12
8 51 22 1 5 3 5 7 1 7
5 24 7 13 1 2 1
8 76 1 1 1 37 1 1 1 33
6 33 3 10 14 2 4 0
1 0 0
7 50 13 18 4 10 2 0 3
4 1 1 0 0 0
6 63 17 2 6 13 13 12
5 44 15 1 5 14 9
1 37 37
3 49 24 1 24
2 0 0 0
4 10 10 0 0 0
2 13 1 12
5 35 14 1 19 1 0
1 40 40
7 77 12 1 1 12 0 1 50
2 3 3 0
2 10 5 5
6 29 4 8 1 12 3 1
8 60 10 10 1 9 3 6 15 6
5 27 3 10 1 10 3
2 10 7 3
7 0 0 0 0 0 0 0 0
8 70 12 12 1 1 3 9 15 17
1 25 25
8 62 13 6 1 3 6 16 14 3
8 36 7 9 1 8 3 6 2 0
4 15 4 1 10 0
8 55 1 8 9 1 8 1 3 24
7 58 2 5 10 2 8 31 0
7 23 4 3 3 0 1 8 4
4 12 1 1 1 9
6 34 16 6 1 5 6 0
4 6 1 1 1 3
2 5 1 4
8 34 5 17 3 1 6 1 0 1
4 20 1 1 16 2
5 57 16 3 5 11 22`

func solve687C(n, k int, coins []int) []int {
	units := (k + 64) / 64
	dp := make([][]uint64, k+1)
	for i := range dp {
		dp[i] = make([]uint64, units)
	}
	dp[0][0] = 1

	for _, c := range coins {
		for j := k; j >= c; j-- {
			src := dp[j-c]
			dst := dp[j]
			for idx := 0; idx < units; idx++ {
				dst[idx] |= src[idx]
			}
			w := c / 64
			b := uint(c % 64)
			if b == 0 {
				for idx := units - 1; idx >= w; idx-- {
					dst[idx] |= src[idx-w]
				}
			} else {
				for idx := units - 1; idx > w; idx-- {
					dst[idx] |= src[idx-w]<<b | src[idx-w-1]>>(64-b)
				}
				if w < units {
					dst[w] |= src[0] << b
				}
			}
		}
	}

	var ans []int
	for x := 0; x <= k; x++ {
		if (dp[k][x/64]>>uint(x%64))&1 == 1 {
			ans = append(ans, x)
		}
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcasesC), "\n")
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Fprintf(os.Stderr, "case %d malformed\n", idx+1)
			os.Exit(1)
		}
		n, err1 := strconv.Atoi(fields[0])
		k, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid n or k\n", idx+1)
			os.Exit(1)
		}
		if len(fields) != 2+n {
			fmt.Fprintf(os.Stderr, "case %d: expected %d numbers, got %d\n", idx+1, 2+n, len(fields))
			os.Exit(1)
		}
		coins := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(fields[2+i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d: invalid coin\n", idx+1)
				os.Exit(1)
			}
			coins[i] = val
		}

		ans := solve687C(n, k, coins)
		var want strings.Builder
		want.WriteString(strconv.Itoa(len(ans)))
		if len(ans) > 0 {
			want.WriteByte('\n')
			for i, v := range ans {
				if i > 0 {
					want.WriteByte(' ')
				}
				want.WriteString(strconv.Itoa(v))
			}
		}

		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", n, k)
		for i, c := range coins {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(c))
		}
		input.WriteByte('\n')

		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want.String()) {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected:\n%s\ngot:\n%s\n", idx+1, want.String(), got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}

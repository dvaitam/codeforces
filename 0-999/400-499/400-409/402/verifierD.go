package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesD.txt.
const embeddedTestcasesD = `2 4
11 9
2 11 17 19
5 3
18 16 2 18 19
2 11 23
4 4
4 13 12 16
2 11 29 31
3 5
10 7 17
3 11 13 17 23
1 5
11
3 7 13 23 29
2 3
1 19
2 5 29
3 3
10 10 11
17 19 29
2 1
5 19
2
4 2
11 1 16 9
7 29
1 5
14
2 5 11 23 29
2 5
4 17
5 17 23 29 31
3 3
1 14 9
11 23 29
5 3
11 7 14 5 1
5 17 23
3 4
2 18 14
2 7 29 31
3 5
6 7 12
2 7 11 19 31
2 3
6 14
3 19 29
2 4
17 4
3 5 7 19
4 4
9 9 14 12
3 11 13 29
1 4
1
7 11 17 23
4 4
2 19 15 12
5 11 13 29
1 4
16
2 3 5 23
5 3
12 1 3 7 4
19 23 31
1 3
1
5 13 17
3 4
5 20 5
2 11 17 23
2 2
5 16
2 31
5 1
18 13 6 12 19
3
1 5
6
7 11 13 19 31
3 3
17 15 5
5 19 23
1 5
6
2 13 17 23 31
1 2
15
7 29
4 5
6 11 5 16
2 3 17 23 31
5 3
1 3 4 14 20
13 19 29
3 4
17 12 4
2 5 13 31
2 1
11 20
7
1 4
2
2 5 11 17
3 1
14 2 15
13
5 5
9 10 19 15 14
2 3 5 11 19
4 1
12 4 4 1
13
1 2
13
2 29
3 4
18 16 16
2 3 17 23
3 1
17 4 3
13
3 1
16 2 5
23
3 1
1 13 11
5
5 2
6 6 6 8 20
2 13
4 4
2 8 8 10
5 7 13 31
2 2
14 15
13 29
2 4
19 1
2 5 17 29
2 2
1 1
13 23
1 1
2
3
5 5
5 5 13 1 14
3 7 13 17 29
3 5
7 18 13
2 3 5 13 17
4 4
8 16 13 8
7 17 19 31
5 1
9 9 17 12 18
2
5 5
16 8 9 2 20
2 3 13 17 23
4 1
14 13 14 4
19
5 4
6 6 11 16 14
3 5 11 29
3 3
5 12 16
2 23 31
2 3
6 19
11 13 17
1 3
18
2 17 31
3 4
7 12 5
3 5 13 31
1 4
19
3 17 19 23
4 5
18 5 6 5
2 5 7 11 31
2 4
12 20
3 11 13 17
3 2
11 17 11
5 23
4 5
10 8 13 12
7 11 17 19 23
4 1
5 5 1 19
29
5 1
7 20 20 17 4
11
5 2
13 3 2 1 4
13 19
3 1
15 12 19
11
4 2
6 18 18 3
5 23
1 2
17
17 29
2 4
13 9
2 5 17 29
2 4
19 2
3 13 17 31
2 5
15 2
3 11 17 19 29
5 4
18 13 16 9 6
5 7 13 23
3 5
5 15 3
3 11 17 19 23
1 3
16
3 7 11
2 3
4 5
2 5 31
5 5
7 1 2 13 18
3 11 19 29 31
3 3
4 1 8
7 11 19
3 2
1 16 12
13 23
1 1
10
29
4 2
12 13 5 8
7 11
4 3
10 13 20 5
3 13 17
5 4
8 12 12 14 9
3 11 13 17
4 3
4 15 5 12
5 7 13
4 2
4 13 13 15
19 23
5 5
8 13 17 10 16
2 7 13 19 23
4 3
13 8 14 2
2 17 29
1 3
7
3 5 13
2 3
1 8
2 17 31
5 1
5 4 20 20 7
3
4 2
1 17 20 14
3 23
2 2
8 14
17 19
1 4
7
2 11 17 29
5 3
12 11 15 5 20
3 11 23
1 1
9
2
2 5
5 13
7 13 19 29 31
5 5
4 18 4 16 4
5 7 19 23 29
5 3
5 14 9 13 3
13 23 29`

func sieve(limit int) []int {
	isPrime := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= limit; i++ {
		if isPrime[i] {
			for j := i * i; j <= limit; j += i {
				isPrime[j] = false
			}
		}
	}
	var primes []int
	for i := 2; i <= limit; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	return primes
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solve402D(n, m int, arr []int64, badList []int64) int64 {
	primes := sieve(32000)
	bad := make(map[int64]bool, m)
	for _, v := range badList {
		bad[v] = true
	}
	fcache := make(map[int64]int)
	var f func(int64) int
	f = func(x int64) int {
		if v, ok := fcache[x]; ok {
			return v
		}
		orig := x
		score := 0
		for _, p := range primes {
			pp := int64(p)
			if pp*pp > x {
				break
			}
			for x%pp == 0 {
				if bad[pp] {
					score--
				} else {
					score++
				}
				x /= pp
			}
		}
		if x > 1 {
			if bad[x] {
				score--
			} else {
				score++
			}
		}
		fcache[orig] = score
		return score
	}

	var ans int64
	for _, v := range arr {
		ans += int64(f(v))
	}

	prefix := make([]int64, n)
	prefix[0] = arr[0]
	for i := 1; i < n; i++ {
		prefix[i] = gcd(prefix[i-1], arr[i])
	}

	var rem int64 = 1
	for i := n - 1; i >= 0; i-- {
		g := prefix[i] / rem
		s := f(g)
		if s < 0 {
			ans -= int64(s) * int64(i+1)
			rem = prefix[i]
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

func mustAtoi(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}
	return v
}

func parseLineToIntSlice(line string) []int {
	fields := strings.Fields(line)
	res := make([]int, len(fields))
	for i, f := range fields {
		res[i] = mustAtoi(f)
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcasesD), "\n")
	if len(lines)%3 != 0 {
		fmt.Fprintln(os.Stderr, "embedded testcases malformed: expected groups of 3 lines")
		os.Exit(1)
	}
	for i := 0; i < len(lines); i += 3 {
		caseID := i/3 + 1
		header := parseLineToIntSlice(lines[i])
		if len(header) != 2 {
			fmt.Fprintf(os.Stderr, "case %d: header should have 2 numbers\n", caseID)
			os.Exit(1)
		}
		n, m := header[0], header[1]
		aVals := parseLineToIntSlice(lines[i+1])
		if len(aVals) != n {
			fmt.Fprintf(os.Stderr, "case %d: expected %d a-values, got %d\n", caseID, n, len(aVals))
			os.Exit(1)
		}
		bVals := parseLineToIntSlice(lines[i+2])
		if len(bVals) != m {
			fmt.Fprintf(os.Stderr, "case %d: expected %d bad primes, got %d\n", caseID, m, len(bVals))
			os.Exit(1)
		}

		arr64 := make([]int64, n)
		for idx, v := range aVals {
			arr64[idx] = int64(v)
		}
		bad64 := make([]int64, m)
		for idx, v := range bVals {
			bad64[idx] = int64(v)
		}
		want := strconv.FormatInt(solve402D(n, m, arr64, bad64), 10)
		input := lines[i] + "\n" + lines[i+1] + "\n" + lines[i+2] + "\n"
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseID, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", caseID, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines)/3)
}

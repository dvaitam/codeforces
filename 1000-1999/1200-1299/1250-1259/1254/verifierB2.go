package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const solution1254B2Source = `package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct {
	pos   int64
	count int64
}

func primeFactors(x int64) []int64 {
	factors := make([]int64, 0)
	for i := int64(2); i*i <= x; i++ {
		if x%i == 0 {
			factors = append(factors, i)
			for x%i == 0 {
				x /= i
			}
		}
	}
	if x > 1 {
		factors = append(factors, x)
	}
	return factors
}

func groupCost(group []pair, k int64) int64 {
	mid := (k + 1) / 2
	var median int64
	var acc int64
	for _, p := range group {
		acc += p.count
		if acc >= mid {
			median = p.pos
			break
		}
	}
	var cost int64
	for _, p := range group {
		if p.pos > median {
			cost += (p.pos - median) * p.count
		} else {
			cost += (median - p.pos) * p.count
		}
	}
	return cost
}

func costForK(a []int64, k int64) int64 {
	var cost int64
	var group []pair
	var sum int64
	for i, v := range a {
		val := v % k
		for val > 0 {
			need := k - sum
			take := val
			if take > need {
				take = need
			}
			group = append(group, pair{pos: int64(i), count: take})
			sum += take
			val -= take
			if sum == k {
				cost += groupCost(group, k)
				group = group[:0]
				sum = 0
			}
		}
	}
	return cost
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int64, n)
	var total int64
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		total += a[i]
	}
	if total <= 1 {
		fmt.Fprintln(writer, -1)
		return
	}
	factors := primeFactors(total)
	best := int64(1<<63 - 1)
	for _, k := range factors {
		c := costForK(a, k)
		if c < best {
			best = c
		}
	}
	fmt.Fprintln(writer, best)
}
`

// Keep the embedded reference solution reachable.
var _ = solution1254B2Source

type testCase struct {
	a []int64
}

const testcasesRaw = `100
1
2
2
11 5
5
8 19 6 19 1
10
5 13 20 12 16 11 17 14 16 8
1
0
6
14 10 12 13 16 5
9
5 7 7 0 5 10 5 4 16
9
11 16 17 5 14 13 16 11 18
6
11 14 5 12 14 20
9
7 15 8 15 16 16 11 14 14
6
18 17 14 15 7 10
3
19 8 15
5
9 16 17 16 16
10
18 13 9 6 15 16 11 19 2 10
1
6
2
1 18
1
8
10
7 3 16 4 8 7 6 1 13 1
1
11
6
5 7 0 2 3 2
1
1
1
11
5
4 5 5 16 0
7
18 1 7 4 1 0 11
10
20 3 9 10 15 0 9 14 17 19
1
8
7
19 4 15 7 2 10 3
1
14
3
16 18 12
8
16 10 4 10 8 8 19 13
1
17
3
1 8 1
3
5 5 3
8
20 7 16 1 7 7 14 2
5
2 18 7 19 19
6
8 13 8 16 0 4
1
12
7
5 3 16 2 7 3 3
1
5
4
3 6 0 16
8
14 9 17 20 12 6 6 13
7
16 0 18 18 1 13 16
10
5 3 15 11 0 16 3 19 11 9
6
9 0 13 3 3 9
4
0 14 1 13
8
14 6 18 19 2 0 9 0
6
9 2 7 15 6 3
10
11 12 14 4 11 12 3 8 3 3
2
19 10
7
6 3 0 19 15 1 15
5
11 14 4 11 8
8
16 15 13 15 9 12 7 5
8
19 8 17 13 2 18 18 3
2
11 5
9
4 13 2 2 20 1 4 9 12
4
10 14 5 16
5
3 4 17 13 3
6
16 7 16 8 5 5
8
7 12 11 18 4 14 14 0
10
12 5 12 16 1 15 8 12 8 13
8
11 17 10 2 7 17 19 6
7
12 20 0 10 14 16 14
3
3 0 12
4
18 19 12 6
2
12 17
4
8 18 18 6
8
19 4 0 19 13 15 8 16
10
5 14 6 2 11 0 15 17 2 18
8
10 14 8 16 14 0 2 19
6
5 12 8 20 4 1
3
15 12 14
5
4 0 9 17 14
1
11
1
17
7
18 14 6 9 15 20 4
8
17 9 2 8 10 9 10 20
5
20 20 12 16 2
9
20 6 12 19 16 4 16 20 2
5
1 7 14 17 7
9
8 1 3 3 12 11 6 10 11
2
10 14
6
5 15 14 9 14 4
8
20 6 8 10 5 3 7 15
4
11 5 11 4
3
7 8 17
7
12 10 8 19 16 18 10
7
20 9 17 19 20 2 11
5
12 15 5 8 11
8
15 2 5 10 12 4 0 3
6
5 11 2 20 13 0
9
10 7 19 12 17 9 15 20 4
6
10 6 15 3 4 6
6
8 4 13 11 8 2
6
6 7 7 19 1 10
6
20 19 1 4 5 2
7
14 8 4 10 16 18 3
6
20 19 12 7 1 12
8
15 19 10 17 19 19 2 18
9
17 15 12 14 5 13 12 16 14
1
3
`

func parseTestcases() []testCase {
	in := bufio.NewReader(strings.NewReader(testcasesRaw))
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		panic(err)
	}
	cases := make([]testCase, t)
	for i := 0; i < t; i++ {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			panic(err)
		}
		a := make([]int64, n)
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscan(in, &a[j]); err != nil {
				panic(err)
			}
		}
		cases[i] = testCase{a: a}
	}
	return cases
}

type pair struct {
	pos   int64
	count int64
}

func primeFactors(x int64) []int64 {
	factors := make([]int64, 0)
	for i := int64(2); i*i <= x; i++ {
		if x%i == 0 {
			factors = append(factors, i)
			for x%i == 0 {
				x /= i
			}
		}
	}
	if x > 1 {
		factors = append(factors, x)
	}
	return factors
}

func groupCost(group []pair, k int64) int64 {
	mid := (k + 1) / 2
	var median int64
	var acc int64
	for _, p := range group {
		acc += p.count
		if acc >= mid {
			median = p.pos
			break
		}
	}
	var cost int64
	for _, p := range group {
		if p.pos > median {
			cost += (p.pos - median) * p.count
		} else {
			cost += (median - p.pos) * p.count
		}
	}
	return cost
}

func costForK(a []int64, k int64) int64 {
	var cost int64
	var group []pair
	var sum int64
	for i, v := range a {
		val := v % k
		for val > 0 {
			need := k - sum
			take := val
			if take > need {
				take = need
			}
			group = append(group, pair{pos: int64(i), count: take})
			sum += take
			val -= take
			if sum == k {
				cost += groupCost(group, k)
				group = group[:0]
				sum = 0
			}
		}
	}
	return cost
}

func solveExpected(tc testCase) int64 {
	n := len(tc.a)
	a := tc.a
	var total int64
	for i := 0; i < n; i++ {
		total += a[i]
	}
	if total <= 1 {
		return -1
	}
	factors := primeFactors(total)
	best := int64(1<<63 - 1)
	for _, k := range factors {
		c := costForK(a, k)
		if c < best {
			best = c
		}
	}
	return best
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tc.a))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin string, idx int, tc testCase) error {
	input := buildInput(tc)
	expect := solveExpected(tc)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("case %d failed: %v\nstderr: %s\ninput:\n%s", idx, err, string(out), input)
	}
	gotStr := strings.TrimSpace(string(out))
	got, err := strconv.ParseInt(gotStr, 10, 64)
	if err != nil {
		return fmt.Errorf("case %d failed: invalid output %q\ninput:\n%s", idx, gotStr, input)
	}
	if got != expect {
		return fmt.Errorf("case %d failed: expected %d got %d\ninput:\n%s", idx, expect, got, input)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tcs := parseTestcases()
	for i, tc := range tcs {
		if err := runCase(bin, i+1, tc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tcs))
}

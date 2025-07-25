package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const maxTemp = 200000

type TestCase struct {
	input    string
	expected []int
}

func solveCase(n, k, q int, ranges [][2]int, queries [][2]int) []int {
	diff := make([]int, maxTemp+2)
	for _, r := range ranges {
		l, rr := r[0], r[1]
		diff[l]++
		if rr+1 <= maxTemp {
			diff[rr+1]--
		}
	}
	freq := make([]int, maxTemp+1)
	for i := 1; i <= maxTemp; i++ {
		freq[i] = freq[i-1] + diff[i]
	}
	pref := make([]int, maxTemp+1)
	for i := 1; i <= maxTemp; i++ {
		pref[i] = pref[i-1]
		if freq[i] >= k {
			pref[i]++
		}
	}
	res := make([]int, q)
	for i, qu := range queries {
		a, b := qu[0], qu[1]
		if a < 1 {
			a = 1
		}
		if b > maxTemp {
			b = maxTemp
		}
		if a > b {
			res[i] = 0
		} else {
			res[i] = pref[b] - pref[a-1]
		}
	}
	return res
}

func generateTest() TestCase {
	n := rand.Intn(5) + 1
	k := rand.Intn(n) + 1
	q := rand.Intn(5) + 1

	ranges := make([][2]int, n)
	for i := 0; i < n; i++ {
		l := rand.Intn(30) + 1
		r := l + rand.Intn(30-l+1)
		ranges[i] = [2]int{l, r}
	}
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		a := rand.Intn(30) + 1
		b := a + rand.Intn(30-a+1)
		queries[i] = [2]int{a, b}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, k, q))
	for _, r := range ranges {
		sb.WriteString(fmt.Sprintf("%d %d\n", r[0], r[1]))
	}
	for _, qu := range queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", qu[0], qu[1]))
	}
	expected := solveCase(n, k, q, ranges, queries)
	return TestCase{input: sb.String(), expected: expected}
}

func runBinary(bin string, input string) ([]int, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("%v: %s", err, errBuf.String())
	}
	fields := strings.Fields(out.String())
	res := make([]int, len(fields))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		res[i] = v
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	const tests = 100
	for t := 0; t < tests; t++ {
		tc := generateTest()
		got, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s", t+1, err, tc.input)
			os.Exit(1)
		}
		if len(got) != len(tc.expected) {
			fmt.Fprintf(os.Stderr, "wrong number of lines on test %d: expected %d got %d\ninput:\n%s", t+1, len(tc.expected), len(got), tc.input)
			os.Exit(1)
		}
		for i := range tc.expected {
			if got[i] != tc.expected[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d line %d: expected %d got %d\ninput:\n%s", t+1, i+1, tc.expected[i], got[i], tc.input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed!")
}

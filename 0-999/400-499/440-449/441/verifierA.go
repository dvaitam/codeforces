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

type testCase struct {
	n      int
	v      int
	prices [][]int
}

func (tc testCase) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.v))
	for i := 0; i < tc.n; i++ {
		sb.WriteString(fmt.Sprintf("%d", len(tc.prices[i])))
		for _, p := range tc.prices[i] {
			sb.WriteString(fmt.Sprintf(" %d", p))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	v := rng.Intn(100) + 1
	prices := make([][]int, n)
	for i := range prices {
		k := rng.Intn(10) + 1
		prices[i] = make([]int, k)
		for j := 0; j < k; j++ {
			prices[i][j] = rng.Intn(200) + 1
		}
	}
	return testCase{n, v, prices}
}

func expected(tc testCase) []int {
	var res []int
	for i := 0; i < tc.n; i++ {
		ok := false
		for _, p := range tc.prices[i] {
			if p < tc.v {
				ok = true
				break
			}
		}
		if ok {
			res = append(res, i+1)
		}
	}
	return res
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	input := tc.Input()
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	gotCnt, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("bad count: %v", err)
	}
	exp := expected(tc)
	if gotCnt != len(exp) {
		return fmt.Errorf("expected count %d got %d", len(exp), gotCnt)
	}
	if gotCnt > 0 {
		if len(fields[1:]) != gotCnt {
			return fmt.Errorf("expected %d indices got %d", gotCnt, len(fields[1:]))
		}
		for i, idx := range exp {
			val, err := strconv.Atoi(fields[1+i])
			if err != nil {
				return fmt.Errorf("bad index: %v", err)
			}
			if val != idx {
				return fmt.Errorf("index %d: expected %d got %d", i+1, idx, val)
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []testCase{
		{n: 1, v: 5, prices: [][]int{{3, 5}}},
		{n: 3, v: 7, prices: [][]int{{8, 9}, {1}, {7, 4, 9}}},
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.Input())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

type testCaseB struct {
	arrays [][]int64
}

func solveB(arrs [][]int64) int64 {
	const INF int64 = 1 << 62
	minFirst := INF
	minSecond := INF
	var sumSecond int64
	for _, arr := range arrs {
		first, second := INF, INF
		for _, x := range arr {
			if x < first {
				second = first
				first = x
			} else if x < second {
				second = x
			}
		}
		if first < minFirst {
			minFirst = first
		}
		if second < minSecond {
			minSecond = second
		}
		sumSecond += second
	}
	return minFirst + sumSecond - minSecond
}

func generateCasesB() []testCaseB {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCaseB, 0, 100)
	for len(cases) < 100 {
		n := rng.Intn(4) + 1
		arrs := make([][]int64, n)
		for i := 0; i < n; i++ {
			m := rng.Intn(4) + 2
			a := make([]int64, m)
			for j := 0; j < m; j++ {
				a[j] = int64(rng.Intn(20) + 1)
			}
			arrs[i] = a
		}
		cases = append(cases, testCaseB{arrays: arrs})
	}
	return cases
}

func runCase(bin string, tc testCaseB) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.arrays)))
	for _, arr := range tc.arrays {
		sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(gotStr, 10, 64)
	if err != nil {
		return fmt.Errorf("non-integer output %q", gotStr)
	}
	exp := solveB(tc.arrays)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCasesB()
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

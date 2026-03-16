package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesBRaw = `100
3 1 3 12 6 3
7 3 20 7 20 2 19 6 14 6 4 5
8 8 17 9 2 1 12 15 11 13 7 5 2 8 6 1 3 4
5 3 6 5 17 17 12 5 2 4
9 9 12 19 12 12 15 6 13 15 17 4 8 3 9 5 7 6 2 1
8 8 16 8 11 6 20 9 16 10 5 7 6 8 4 2 1 3
8 2 11 1 7 4 2 19 2 9 4 6
4 2 9 8 7 2 4 3
3 1 12 12 6 1
3 1 4 3 1 1
3 2 9 5 6 3 1
3 2 19 2 8 1 3
3 2 20 4 10 2 3
3 2 15 18 20 3 1
7 7 13 20 5 16 8 3 11 7 1 6 4 5 2 3
8 3 11 9 9 20 14 1 18 5 1 3 8
5 2 6 4 15 8 17 1 2
6 6 15 3 9 3 19 8 5 6 3 2 4 1
3 1 2 13 14 1
4 1 8 4 4 1 2
6 1 7 1 17 15 15 10 5
9 4 7 14 14 17 1 19 19 2 14 9 3 1 6
10 6 1 17 4 20 12 10 12 10 1 14 2 10 5 9 6 1
10 1 14 16 15 7 19 20 3 1 10 1 6
7 6 3 8 16 7 4 19 12 4 6 7 2 5 3
4 3 4 4 3 20 3 4 2
6 6 4 1 20 16 2 16 3 6 4 1 2 5
10 9 16 14 16 10 13 8 6 16 20 9 9 7 2 5 6 10 1 4 8
5 5 5 14 3 3 2 2 3 5 1 4
10 3 17 10 4 5 18 14 4 11 17 8 9 5 3
5 4 8 13 12 19 5 4 5 3 1
9 3 13 17 2 16 9 13 9 14 16 6 9 8
4 2 18 20 7 13 4 3
3 2 15 17 15 3 1
4 1 13 7 19 20 4
6 1 13 18 7 9 19 19 2
10 10 5 1 20 14 16 9 17 19 6 15 4 2 6 1 10 5 7 3 9 8
10 5 17 15 1 3 20 12 6 13 9 5 1 3 8 4 7
7 2 1 10 18 15 1 12 2 5 4
10 4 10 16 5 16 18 10 3 9 11 10 6 5 7 9
4 2 13 20 17 5 1 2
3 1 15 18 8 3
7 1 4 4 13 12 7 11 12 1
8 8 12 6 16 15 10 15 5 15 4 3 7 2 1 8 5 6
8 3 12 5 5 8 9 18 13 13 6 3 8
8 7 10 18 20 3 12 10 13 16 3 8 7 4 5 1 6
8 7 5 1 4 12 6 12 3 14 1 5 3 2 4 6 7
10 3 12 11 7 16 4 5 7 11 9 5 7 6 5
4 3 7 8 8 20 1 2 3
3 1 6 3 14 2
7 2 11 17 19 4 11 20 13 2 1
9 8 16 20 11 18 20 20 3 19 17 9 8 4 6 7 2 5 3
10 1 4 15 19 5 4 17 6 3 13 10 8
3 2 4 12 8 1 3
5 4 3 11 15 2 16 2 1 5 4
3 1 17 18 2 1
6 5 1 17 11 17 8 5 3 4 1 5 6
6 1 15 7 2 20 7 13 3
9 9 17 6 17 4 5 7 6 13 7 5 6 4 2 7 8 9 3 1
4 1 16 9 10 17 4
7 2 14 5 18 4 1 20 18 7 2
6 2 13 19 2 5 1 9 6 4
3 3 8 5 20 2 1 3
4 2 18 6 3 15 3 1
5 3 9 17 19 3 14 4 1 2
7 6 4 9 1 7 14 11 9 5 6 4 2 7 1
5 4 15 12 13 16 20 3 2 5 4
10 4 16 19 11 10 3 6 12 20 16 8 10 3 5 7
6 5 10 4 1 1 7 11 1 3 5 4 2
10 2 14 16 1 10 19 19 5 7 5 6 10 7
4 4 9 3 16 16 2 1 4 3
6 5 11 19 20 13 17 14 2 6 1 5 4
5 5 13 14 4 15 13 4 5 2 3 1
6 2 2 18 17 3 20 18 6 1
3 2 14 13 8 3 2
4 3 17 12 17 16 1 3 2
6 3 1 1 16 2 5 5 2 3 6
3 3 5 10 4 3 1 2
9 3 2 10 17 9 16 2 18 12 11 2 6 1
8 6 9 16 10 17 20 5 1 2 6 4 8 1 3 7
3 3 3 18 17 3 2 1
9 4 6 6 20 2 1 19 12 6 10 1 9 2 5
6 4 3 12 4 20 3 8 2 5 6 1
3 3 13 3 17 2 1 3
9 7 5 5 14 5 15 12 2 19 6 9 8 4 5 7 2 6
5 3 5 1 9 6 5 4 3 2
10 8 7 14 14 9 8 12 2 13 20 1 7 5 1 9 4 10 3 6
7 2 15 15 12 17 20 15 8 5 7
5 4 10 12 14 4 17 2 4 1 5
10 10 17 15 3 13 15 20 12 18 12 6 3 4 10 9 7 6 2 5 8 1
9 5 9 11 1 13 19 2 7 15 4 2 1 7 3 6
8 7 6 11 3 17 20 16 13 20 4 8 1 5 6 2 7
7 4 5 2 1 11 13 16 20 1 4 5 2
6 6 20 11 6 11 10 13 5 6 4 2 3 1
6 5 12 8 12 12 6 8 5 6 2 1 4
6 4 9 19 7 17 6 1 4 6 3 5
5 2 20 6 16 20 1 2 5
6 1 20 3 15 7 6 9 4
3 1 12 4 10 1
8 8 13 4 3 15 6 5 15 3 5 2 4 8 1 3 7 6
`

type testcaseB struct {
	n      int
	k      int
	beauty []int64
	caps   []int
}

func expectedB(tc testcaseB) int64 {
	n := tc.n
	c := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		c[i] = tc.beauty[i-1]
	}
	caps := tc.caps
	isCap := make([]bool, n+1)
	for _, x := range caps {
		isCap[x] = true
	}
	var totalBeauty int64
	for i := 1; i <= n; i++ {
		totalBeauty += c[i]
	}
	var sumCap, sumCapSq int64
	for _, x := range caps {
		sumCap += c[x]
		sumCapSq += c[x] * c[x]
	}
	var cycleSum int64
	for i := 1; i < n; i++ {
		cycleSum += c[i] * c[i+1]
	}
	cycleSum += c[n] * c[1]
	T1 := totalBeauty*sumCap - sumCapSq
	var T2 int64
	for _, x := range caps {
		prev := x - 1
		if prev < 1 {
			prev = n
		}
		next := x + 1
		if next > n {
			next = 1
		}
		T2 += c[x] * (c[prev] + c[next])
	}
	capTotal := T1 - T2
	capPairsSum := (sumCap*sumCap - sumCapSq) / 2
	var capNeighborSum int64
	for i := 1; i < n; i++ {
		if isCap[i] && isCap[i+1] {
			capNeighborSum += c[i] * c[i+1]
		}
	}
	if isCap[n] && isCap[1] {
		capNeighborSum += c[n] * c[1]
	}
	capContribution := capTotal - capPairsSum + capNeighborSum
	return cycleSum + capContribution
}

func runCase(exe string, tc testcaseB) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	for i := 0; i < tc.n; i++ {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprint(tc.beauty[i]))
	}
	input.WriteByte('\n')
	for i := 0; i < tc.k; i++ {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprint(tc.caps[i]))
	}
	input.WriteByte('\n')
	exp := fmt.Sprint(expectedB(tc))
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	sc := bufio.NewScanner(strings.NewReader(testcasesBRaw))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(sc.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		sc.Scan()
		n, _ := strconv.Atoi(sc.Text())
		sc.Scan()
		k, _ := strconv.Atoi(sc.Text())
		beauty := make([]int64, n)
		for i := 0; i < n; i++ {
			sc.Scan()
			v, _ := strconv.ParseInt(sc.Text(), 10, 64)
			beauty[i] = v
		}
		caps := make([]int, k)
		for i := 0; i < k; i++ {
			sc.Scan()
			v, _ := strconv.Atoi(sc.Text())
			caps[i] = v
		}
		tc := testcaseB{n, k, beauty, caps}
		if err := runCase(exe, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

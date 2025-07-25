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
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	sc := bufio.NewScanner(bytes.NewReader(data))
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

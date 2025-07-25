package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func expectedSales(n, f int, kVals, lVals []int64) int64 {
	base := int64(0)
	diffs := make([]int64, n)
	for i := 0; i < n; i++ {
		sold := kVals[i]
		if lVals[i] < kVals[i] {
			sold = lVals[i]
		}
		base += sold
		doubled := kVals[i] * 2
		soldDouble := doubled
		if lVals[i] < doubled {
			soldDouble = lVals[i]
		}
		diffs[i] = soldDouble - sold
	}
	sort.Slice(diffs, func(i, j int) bool { return diffs[i] > diffs[j] })
	if f > n {
		f = n
	}
	add := int64(0)
	for i := 0; i < f; i++ {
		if diffs[i] > 0 {
			add += diffs[i]
		}
	}
	return base + add
}

func runCase(bin string, n, f int, kVals, lVals []int64) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, f))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", kVals[i], lVals[i]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	outStr := strings.TrimSpace(string(out))
	expected := fmt.Sprintf("%d", expectedSales(n, f, kVals, lVals))
	if outStr != expected {
		return fmt.Errorf("expected %s got %s", expected, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Println("could not open testcasesB.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseNum := 0; caseNum < t; caseNum++ {
		if !scan.Scan() {
			fmt.Printf("missing n on case %d\n", caseNum+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Printf("missing f on case %d\n", caseNum+1)
			os.Exit(1)
		}
		fval, _ := strconv.Atoi(scan.Text())
		kVals := make([]int64, n)
		lVals := make([]int64, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Printf("missing k on case %d index %d\n", caseNum+1, i+1)
				os.Exit(1)
			}
			kv, _ := strconv.ParseInt(scan.Text(), 10, 64)
			if !scan.Scan() {
				fmt.Printf("missing l on case %d index %d\n", caseNum+1, i+1)
				os.Exit(1)
			}
			lv, _ := strconv.ParseInt(scan.Text(), 10, 64)
			kVals[i] = kv
			lVals[i] = lv
		}
		if err := runCase(bin, n, fval, kVals, lVals); err != nil {
			fmt.Printf("case %d failed: %v\n", caseNum+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

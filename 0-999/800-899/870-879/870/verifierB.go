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

func solveB(n, k int, arr []int) string {
	minVal, maxVal := arr[0], arr[0]
	for _, v := range arr[1:] {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}
	var ans int
	switch {
	case k == 1:
		ans = minVal
	case k >= 3:
		ans = maxVal
	default:
		if arr[0] > arr[n-1] {
			ans = arr[0]
		} else {
			ans = arr[n-1]
		}
	}
	return fmt.Sprintf("%d\n", ans)
}

func runCase(exe string, input, expect string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expect)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
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
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		k, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			arr[i], _ = strconv.Atoi(scan.Text())
		}
		inputBuilder := &strings.Builder{}
		fmt.Fprintf(inputBuilder, "%d %d\n", n, k)
		for i, v := range arr {
			if i > 0 {
				inputBuilder.WriteByte(' ')
			}
			fmt.Fprintf(inputBuilder, "%d", v)
		}
		inputBuilder.WriteByte('\n')
		expect := solveB(n, k, arr)
		if err := runCase(exe, inputBuilder.String(), expect); err != nil {
			fmt.Printf("case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

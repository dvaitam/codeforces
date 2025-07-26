package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func expected(n, k int, arr []int) float64 {
	sum := make([]int, n)
	if n > 0 {
		sum[0] = arr[0]
		for i := 1; i < n; i++ {
			sum[i] = sum[i-1] + arr[i]
		}
	}
	best := math.Inf(-1)
	for L := k; L <= n; L++ {
		maxSum := sum[L-1]
		for j := L; j < n; j++ {
			s := sum[j] - sum[j-L]
			if s > maxSum {
				maxSum = s
			}
		}
		cur := float64(maxSum) / float64(L)
		if cur > best {
			best = cur
		}
	}
	return best
}

func runCase(bin string, n, k int, arr []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseFloat(gotStr, 64)
	if err != nil {
		return fmt.Errorf("invalid output %q", gotStr)
	}
	want := expected(n, k, arr)
	if math.Abs(got-want) > 1e-6 {
		return fmt.Errorf("expected %.6f got %.6f", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Println("could not open testcasesC.txt:", err)
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
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		k, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			v, _ := strconv.Atoi(scan.Text())
			arr[j] = v
		}
		if err := runCase(bin, n, k, arr); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

func nextPerm(a []int) bool {
	n := len(a)
	i := n - 2
	for i >= 0 && a[i] >= a[i+1] {
		i--
	}
	if i < 0 {
		return false
	}
	j := n - 1
	for a[j] <= a[i] {
		j--
	}
	a[i], a[j] = a[j], a[i]
	for l, r := i+1, n-1; l < r; l, r = l+1, r-1 {
		a[l], a[r] = a[r], a[l]
	}
	return true
}

func solveCase(nums []string, k int) int {
	perm := make([]int, k)
	for i := 0; i < k; i++ {
		perm[i] = i
	}
	const inf = int64(1<<63 - 1)
	best := inf
	for {
		var mn, mx int64
		for idx, s := range nums {
			var v int64
			for _, p := range perm {
				v = v*10 + int64(s[p]-'0')
			}
			if idx == 0 || v < mn {
				mn = v
			}
			if idx == 0 || v > mx {
				mx = v
			}
		}
		if diff := mx - mn; diff < best {
			best = diff
		}
		if !nextPerm(perm) {
			break
		}
	}
	return int(best)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseNum := 0; caseNum < t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		k, _ := strconv.Atoi(scan.Text())
		nums := make([]string, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			nums[j] = scan.Text()
		}
		expected := solveCase(nums, k)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for _, s := range nums {
			input.WriteString(s)
			input.WriteByte('\n')
		}
		cmd := exec.Command(os.Args[1])
		cmd.Stdin = strings.NewReader(input.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", caseNum+1, err)
			os.Exit(1)
		}
		outScan := bufio.NewScanner(bytes.NewReader(out))
		outScan.Split(bufio.ScanWords)
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", caseNum+1)
			os.Exit(1)
		}
		got, err := strconv.Atoi(outScan.Text())
		if err != nil {
			fmt.Printf("bad output for test %d\n", caseNum+1)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed: expected %d got %d\n", caseNum+1, expected, got)
			os.Exit(1)
		}
		if outScan.Scan() {
			fmt.Printf("extra output on test %d\n", caseNum+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}

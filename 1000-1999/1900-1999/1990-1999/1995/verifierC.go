package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/big"
	"os"
	"os/exec"
	"strings"
)

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func isGE(a int64, k int, pb int64, pk int) bool {
	if k == pk {
		return a >= pb
	}
	if k > pk {
		d := k - pk
		exp := int64(1 << uint(d))
		powa := new(big.Int).Exp(big.NewInt(a), big.NewInt(exp), nil)
		return powa.Cmp(big.NewInt(pb)) >= 0
	}
	d := pk - k
	exp := int64(1 << uint(d))
	powpb := new(big.Int).Exp(big.NewInt(pb), big.NewInt(exp), nil)
	return big.NewInt(a).Cmp(powpb) >= 0
}

func expectedC(arr []int64) int64 {
	n := len(arr)
	total := int64(0)
	prevBase := arr[0]
	prevK := 0
	ok := true
	for i := 1; i < n && ok; i++ {
		currA := arr[i]
		var currK int
		if currA == 1 {
			if prevBase != 1 {
				ok = false
				break
			}
			currK = 0
		} else {
			minK := math.MaxInt
			found := false
			start := maxInt(0, prevK-6)
			end := prevK + 6
			for ck := start; ck <= end; ck++ {
				if isGE(currA, ck, prevBase, prevK) {
					if ck < minK {
						minK = ck
					}
					found = true
				}
			}
			if !found {
				ok = false
				break
			}
			currK = minK
		}
		total += int64(currK)
		prevBase = currA
		prevK = currK
	}
	if !ok {
		return -1
	}
	return total
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 1 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		var n int
		fmt.Sscan(fields[0], &n)
		if len(fields) != 1+n {
			fmt.Printf("test %d: invalid number of values\n", idx)
			os.Exit(1)
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Sscan(fields[1+i], &arr[i])
		}
		expect := expectedC(arr)
		input := fmt.Sprintf("1\n%d\n", n)
		for i := 0; i < n; i++ {
			input += fmt.Sprintf("%d ", arr[i])
		}
		input += "\n"
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		res := strings.TrimSpace(out.String())
		var got int64
		if _, err := fmt.Sscan(res, &got); err != nil {
			fmt.Printf("test %d: failed to parse output %q\n", idx, res)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

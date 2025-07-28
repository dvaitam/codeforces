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

func expectedSum(arr []int) int {
	cnt := make([]int, 101)
	for _, v := range arr {
		if v >= 0 && v <= 100 {
			cnt[v]++
		}
	}
	seq := make([]int, 0, len(arr))
	for i := 0; i <= 100; i++ {
		if cnt[i] > 0 {
			seq = append(seq, i)
			cnt[i]--
		}
	}
	for i := 0; i <= 100; i++ {
		for cnt[i] > 0 {
			seq = append(seq, i)
			cnt[i]--
		}
	}
	seen := make([]bool, 102)
	mex := 0
	sum := 0
	for _, v := range seq {
		if v >= 0 && v < len(seen) {
			seen[v] = true
		}
		for mex < len(seen) && seen[mex] {
			mex++
		}
		sum += mex
	}
	return sum
}

func sumMex(arr []int) int {
	seen := make(map[int]bool)
	mex := 0
	sum := 0
	for _, v := range arr {
		seen[v] = true
		for {
			if !seen[mex] {
				break
			}
			mex++
		}
		sum += mex
	}
	return sum
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Println("failed to open testcasesA.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	caseNum := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		caseNum++
		fields := strings.Fields(line)
		if len(fields) < 1 {
			fmt.Printf("case %d invalid format\n", caseNum)
			os.Exit(1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil || len(fields)-1 != n {
			fmt.Printf("case %d invalid data\n", caseNum)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(fields[i+1])
			if err != nil {
				fmt.Printf("case %d invalid integer %q\n", caseNum, fields[i+1])
				os.Exit(1)
			}
			arr[i] = val
		}
		exp := expectedSum(arr)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(arr[i]))
		}
		input.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		err = cmd.Run()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\nstderr: %s\n", caseNum, err, errBuf.String())
			os.Exit(1)
		}
		outFields := strings.Fields(strings.TrimSpace(out.String()))
		if len(outFields) != n {
			fmt.Printf("case %d: expected %d numbers, got %d\n", caseNum, n, len(outFields))
			os.Exit(1)
		}
		outArr := make([]int, n)
		freq := map[int]int{}
		for i, s := range outFields {
			val, err := strconv.Atoi(s)
			if err != nil {
				fmt.Printf("case %d: invalid output integer %q\n", caseNum, s)
				os.Exit(1)
			}
			outArr[i] = val
			freq[val]++
		}
		inFreq := map[int]int{}
		for _, v := range arr {
			inFreq[v]++
		}
		for k, c := range inFreq {
			if freq[k] != c {
				fmt.Printf("case %d: output is not a permutation of input\n", caseNum)
				os.Exit(1)
			}
		}
		gotSum := sumMex(outArr)
		if gotSum != exp {
			fmt.Printf("case %d failed: expected sum %d got %d\n", caseNum, exp, gotSum)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", caseNum)
}

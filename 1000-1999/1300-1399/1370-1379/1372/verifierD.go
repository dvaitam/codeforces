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

func runCandidate(bin string, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(a []int64) int64 {
	n := len(a)
	if n == 1 {
		return a[0]
	}
	even := make([]int64, 2*n+1)
	odd := make([]int64, 2*n+1)
	for i := 0; i < 2*n; i++ {
		v := a[i%n]
		if i%2 == 0 {
			even[i+1] = even[i] + v
			odd[i+1] = odd[i]
		} else {
			odd[i+1] = odd[i] + v
			even[i+1] = even[i]
		}
	}
	k := n / 2
	var best int64
	for i := 0; i < n; i++ {
		var s int64
		if i%2 == 0 {
			s = even[i+2*k+1] - even[i]
		} else {
			s = odd[i+2*k+1] - odd[i]
		}
		if s > best {
			best = s
		}
	}
	return best
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
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
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			v, _ := strconv.ParseInt(scan.Text(), 10, 64)
			arr[i] = v
		}
		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", n)
		for i, v := range arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", v)
		}
		input.WriteByte('\n')
		expect := solveCase(arr)
		out, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Printf("case %d: invalid output\n", caseIdx+1)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %d got %d\n", caseIdx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

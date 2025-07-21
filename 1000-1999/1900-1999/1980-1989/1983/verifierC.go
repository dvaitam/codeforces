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

func expected(n int, arr [3][]int64) string {
	need := int64(0)
	for _, v := range arr[0] {
		need += v
	}
	need = (need + 2) / 3
	perms := [6][3]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}
	seg := [3][2]int{}
	for _, p := range perms {
		ptr := 0
		ok := true
		for i := range seg {
			seg[i][0] = 0
			seg[i][1] = -1
		}
		for _, x := range p {
			start := ptr
			cur := int64(0)
			for ptr < n && cur < need {
				cur += arr[x][ptr]
				ptr++
			}
			if cur < need {
				ok = false
				break
			}
			seg[x][0] = start
			seg[x][1] = ptr - 1
		}
		if ok {
			return fmt.Sprintf("%d %d %d %d %d %d",
				seg[0][0]+1, seg[0][1]+1,
				seg[1][0]+1, seg[1][1]+1,
				seg[2][0]+1, seg[2][1]+1)
		}
	}
	return "-1"
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
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
	scan := bufio.NewScanner(file)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid testcases file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for idx := 1; idx <= t; idx++ {
		if !scan.Scan() {
			fmt.Printf("missing n for case %d\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		var arr [3][]int64
		for i := 0; i < 3; i++ {
			arr[i] = make([]int64, n)
			for j := 0; j < n; j++ {
				if !scan.Scan() {
					fmt.Printf("missing values in case %d\n", idx)
					os.Exit(1)
				}
				v, _ := strconv.ParseInt(scan.Text(), 10, 64)
				arr[i][j] = v
			}
		}
		expect := expected(n, arr)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d\n", n))
		for k := 0; k < 3; k++ {
			for j := 0; j < n; j++ {
				if j > 0 {
					input.WriteByte(' ')
				}
				input.WriteString(fmt.Sprintf("%d", arr[k][j]))
			}
			input.WriteByte('\n')
		}
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed:\nexpected: %s\n got: %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}

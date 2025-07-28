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

const mod int64 = 998244353

func comb(n, k int64) int64 {
	if k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	res := int64(1)
	for i := int64(1); i <= k; i++ {
		res = res * (n - i + 1) / i
	}
	return res % mod
}

func solve(n, m int64) int64 {
	if n == 1 || m == 1 {
		return 0
	}
	if n == 2 {
		return 2 * comb(m+2, 4) % mod
	}
	if m == 2 {
		return 2 * comb(n+2, 4) % mod
	}
	return 0
}

func parseTests(path string) ([][2]int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	in := bufio.NewReader(f)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return nil, err
	}
	cases := make([][2]int64, t)
	for i := 0; i < t; i++ {
		fmt.Fscan(in, &cases[i][0], &cases[i][1])
	}
	return cases, nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	cases, err := parseTests("testcasesE.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for idx, tc := range cases {
		exp := solve(tc[0], tc[1])
		input := fmt.Sprintf("%d %d\n", tc[0], tc[1])
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", idx+1, err, out)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil || val != exp {
			fmt.Printf("case %d failed: expected %d got %s\n", idx+1, exp, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

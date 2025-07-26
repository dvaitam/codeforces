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

const inf int64 = 1 << 60

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func expected(n int, a, b int64, s string) int64 {
	dp0 := make([]int64, n+1)
	dp1 := make([]int64, n+1)
	for i := 0; i <= n; i++ {
		dp0[i] = inf
		dp1[i] = inf
	}
	dp0[0] = b
	for i := 1; i <= n; i++ {
		if s[i-1] == '1' {
			dp0[i] = inf
			dp1[i] = dp1[i-1] + a + 2*b
		} else {
			dp0[i] = min64(dp0[i-1]+a+b, dp1[i-1]+2*a+b)
			dp1[i] = min64(dp1[i-1]+a+2*b, dp0[i-1]+2*a+2*b)
		}
	}
	return dp0[n]
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	cand, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("bad test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	types := make([]struct {
		n    int
		a, b int64
		s    string
	}, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		a64, _ := strconv.ParseInt(scan.Text(), 10, 64)
		scan.Scan()
		b64, _ := strconv.ParseInt(scan.Text(), 10, 64)
		scan.Scan()
		s := scan.Text()
		types[i] = struct {
			n    int
			a, b int64
			s    string
		}{n, a64, b64, s}
	}
	cmd := exec.Command(cand)
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i, c := range types {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		got, _ := strconv.ParseInt(outScan.Text(), 10, 64)
		exp := expected(c.n, c.a, c.b, c.s)
		if got != exp {
			fmt.Printf("test %d failed: expected %d got %d\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", t)
}

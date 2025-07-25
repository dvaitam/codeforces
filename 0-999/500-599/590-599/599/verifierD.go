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

type pair struct{ n, m int64 }

func expected(x int64) string {
	res := make([]pair, 0)
	for n := int64(1); ; n++ {
		minSq := n * (n + 1) * (2*n + 1) / 6
		if minSq > x {
			break
		}
		denom := n * (n + 1)
		sixx := 6 * x
		if sixx%denom != 0 {
			continue
		}
		t := sixx / denom
		if (t+n-1)%3 != 0 {
			continue
		}
		m := (t + n - 1) / 3
		if m < n {
			continue
		}
		if n*(n+1)*(3*m-n+1)/6 == x {
			res = append(res, pair{int64(n), m})
			if m != n {
				res = append(res, pair{m, int64(n)})
			}
		}
	}
	// sort
	for i := 0; i < len(res); i++ {
		for j := i + 1; j < len(res); j++ {
			if res[j].n < res[i].n || (res[j].n == res[i].n && res[j].m < res[i].m) {
				res[i], res[j] = res[j], res[i]
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(res)))
	for _, p := range res {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.n, p.m))
	}
	return sb.String()
}

func runCase(exe, input, exp string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp = strings.TrimSpace(exp)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		x, _ := strconv.ParseInt(scan.Text(), 10, 64)
		input := fmt.Sprintf("%d\n", x)
		exp := expected(x)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

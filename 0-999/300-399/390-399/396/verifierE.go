package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type pair struct{ p, e int }

func run(bin, input string) (string, error) {
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

func factor(x int) map[int]int {
	m := make(map[int]int)
	for p := 2; p*p <= x; p++ {
		for x%p == 0 {
			m[p]++
			x /= p
		}
	}
	if x > 1 {
		m[x]++
	}
	return m
}

func solveCase(m int, pairs []pair, k int) string {
	curr := make(map[int]int)
	for _, pr := range pairs {
		curr[pr.p] = pr.e
	}
	ans := make(map[int]int)
	for h := 0; len(curr) > 0 && h <= k; h++ {
		kh := k - h
		next := make(map[int]int)
		for p, e0 := range curr {
			if e0 > kh {
				ans[p] += e0 - kh
			}
			A := e0
			if A > kh {
				A = kh
			}
			if A <= 0 {
				continue
			}
			fac := factor(p - 1)
			for f, cnt := range fac {
				next[f] += cnt * A
			}
		}
		curr = next
	}
	primes := make([]int, 0, len(ans))
	for p := range ans {
		primes = append(primes, p)
	}
	sort.Ints(primes)
	var buf bytes.Buffer
	fmt.Fprintln(&buf, len(primes))
	for _, p := range primes {
		fmt.Fprintf(&buf, "%d %d\n", p, ans[p])
	}
	return strings.TrimSpace(buf.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	cases := make([]struct {
		m     int
		pairs []pair
		k     int
	}, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		mval, _ := strconv.Atoi(scan.Text())
		pairs := make([]pair, mval)
		for j := 0; j < mval; j++ {
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			p, _ := strconv.Atoi(scan.Text())
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			e, _ := strconv.Atoi(scan.Text())
			pairs[j] = pair{p, e}
		}
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		kval, _ := strconv.Atoi(scan.Text())
		cases[i] = struct {
			m     int
			pairs []pair
			k     int
		}{mval, pairs, kval}
	}
	for idx, tc := range cases {
		var input bytes.Buffer
		fmt.Fprintln(&input, tc.m)
		for _, pr := range tc.pairs {
			fmt.Fprintf(&input, "%d %d\n", pr.p, pr.e)
		}
		fmt.Fprintln(&input, tc.k)
		expected := solveCase(tc.m, tc.pairs, tc.k)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, input.String())
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected:\n%s\ngot:\n%s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

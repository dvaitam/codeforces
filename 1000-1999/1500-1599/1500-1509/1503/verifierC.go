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

type testCaseC struct {
	n      int
	cities []struct{ a, c int }
}

func parseTests(path string) ([]testCaseC, error) {
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
	cases := make([]testCaseC, t)
	for i := 0; i < t; i++ {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return nil, err
		}
		cities := make([]struct{ a, c int }, n)
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscan(in, &cities[j].a, &cities[j].c); err != nil {
				return nil, err
			}
		}
		cases[i] = testCaseC{n: n, cities: cities}
	}
	return cases, nil
}

func solve(tc testCaseC) int64 {
	cities := make([]struct{ a, c int }, len(tc.cities))
	copy(cities, tc.cities)
	sort.Slice(cities, func(i, j int) bool { return cities[i].a < cities[j].a })
	ans := int64(cities[0].c)
	reach := int64(cities[0].a + cities[0].c)
	for i := 1; i < len(cities); i++ {
		if int64(cities[i].a) > reach {
			ans += int64(cities[i].a) - reach
			reach = int64(cities[i].a)
		}
		ans += int64(cities[i].c)
		if int64(cities[i].a+cities[i].c) > reach {
			reach = int64(cities[i].a + cities[i].c)
		}
	}
	return ans
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
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	cases, err := parseTests("testcasesC.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for idx, tc := range cases {
		exp := solve(tc)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for _, c := range tc.cities {
			sb.WriteString(fmt.Sprintf("%d %d\n", c.a, c.c))
		}
		out, err := runBinary(bin, sb.String())
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

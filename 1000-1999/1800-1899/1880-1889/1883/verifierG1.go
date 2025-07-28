package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCaseG1 struct {
	n int
	m int
	a []int
	b []int
}

func solveCaseG1(tc testCaseG1) string {
	n := tc.n
	a := append([]int{1}, tc.a...)
	b := append([]int(nil), tc.b...)
	sort.Ints(a)
	sort.Ints(b)
	j := 0
	match := 0
	for _, x := range a {
		for j < n && b[j] <= x {
			j++
		}
		if j == n {
			break
		}
		match++
		j++
	}
	return fmt.Sprint(n - match)
}

func runCaseG1(bin string, tc testCaseG1) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := solveCaseG1(tc)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func randomCaseG1(rng *rand.Rand) testCaseG1 {
	n := rng.Intn(6) + 2
	m := 1
	a := make([]int, n-1)
	for i := range a {
		a[i] = rng.Intn(10) + 1
	}
	b := make([]int, n)
	for i := range b {
		b[i] = rng.Intn(15) + 1
	}
	return testCaseG1{n: n, m: m, a: a, b: b}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCaseG1{{n: 2, m: 1, a: []int{2}, b: []int{3, 4}}}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCaseG1(rng))
	}
	for idx, tc := range cases {
		if err := runCaseG1(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

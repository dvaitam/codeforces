package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCaseF struct {
	n int
	a []int
}

func solveCaseF(tc testCaseF) string {
	n := tc.n
	a := tc.a
	first := make([]bool, n)
	seen := make(map[int]struct{})
	for i, v := range a {
		if _, ok := seen[v]; !ok {
			first[i] = true
			seen[v] = struct{}{}
		}
	}
	last := make([]bool, n)
	seen = make(map[int]struct{})
	for i := n - 1; i >= 0; i-- {
		v := a[i]
		if _, ok := seen[v]; !ok {
			last[i] = true
			seen[v] = struct{}{}
		}
	}
	suf := make([]int, n+1)
	for i := n - 1; i >= 0; i-- {
		if last[i] {
			suf[i] = suf[i+1] + 1
		} else {
			suf[i] = suf[i+1]
		}
	}
	var ans int64
	for i := 0; i < n; i++ {
		if first[i] {
			ans += int64(suf[i])
		}
	}
	return fmt.Sprint(ans)
}

func runCaseF(bin string, tc testCaseF) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.a {
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
	exp := solveCaseF(tc)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func randomCaseF(rng *rand.Rand) testCaseF {
	n := rng.Intn(8) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(10) + 1
	}
	return testCaseF{n: n, a: arr}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCaseF{{n: 1, a: []int{1}}, {n: 3, a: []int{1, 2, 1}}}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCaseF(rng))
	}
	for idx, tc := range cases {
		if err := runCaseF(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

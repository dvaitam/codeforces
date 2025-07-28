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

type testCaseE struct {
	n   int
	arr []int64
}

func solveCaseE(tc testCaseE) string {
	ops := 0
	prev := tc.arr[0]
	for i := 1; i < tc.n; i++ {
		cur := tc.arr[i]
		for cur < prev {
			cur <<= 1
			ops++
		}
		prev = cur
	}
	return fmt.Sprint(ops)
}

func runCaseE(bin string, tc testCaseE) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.arr {
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
	exp := solveCaseE(tc)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func randomCaseE(rng *rand.Rand) testCaseE {
	n := rng.Intn(8) + 1
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = int64(rng.Intn(20) + 1)
	}
	return testCaseE{n: n, arr: arr}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCaseE{{n: 1, arr: []int64{1}}, {n: 2, arr: []int64{5, 1}}}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCaseE(rng))
	}
	for idx, tc := range cases {
		if err := runCaseE(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

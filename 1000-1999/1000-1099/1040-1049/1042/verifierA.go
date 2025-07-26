package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCaseA struct {
	n int64
	m int64
	a []int64
}

func solveA(tc testCaseA) (int64, int64) {
	maxx := int64(0)
	sum := int64(0)
	for _, v := range tc.a {
		if v > maxx {
			maxx = v
		}
		sum += v
	}
	avg := (sum + tc.m) / int64(tc.n)
	if (sum+tc.m)%int64(tc.n) != 0 {
		avg++
	}
	if maxx > avg {
		avg = maxx
	}
	return avg, maxx + tc.m
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() []testCaseA {
	rng := rand.New(rand.NewSource(42))
	tests := make([]testCaseA, 100)
	for i := range tests {
		n := int64(rng.Intn(10) + 1)
		m := int64(rng.Intn(50))
		a := make([]int64, n)
		for j := range a {
			a[j] = int64(rng.Intn(20))
		}
		tests[i] = testCaseA{n, m, a}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
		for j, v := range tc.a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		exp1, exp2 := solveA(tc)
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d: execution error: %v\n", i+1, err)
			return
		}
		var got1, got2 int64
		fmt.Sscan(out, &got1, &got2)
		if got1 != exp1 || got2 != exp2 {
			fmt.Printf("test %d failed:\ninput:%sexpected %d %d got %s\n", i+1, sb.String(), exp1, exp2, out)
			return
		}
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}

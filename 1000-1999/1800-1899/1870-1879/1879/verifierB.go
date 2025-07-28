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

type testCase struct {
	a []int64
	b []int64
}

func expected(tc testCase) string {
	n := int64(len(tc.a))
	sumA, sumB := int64(0), int64(0)
	minA, minB := tc.a[0], tc.b[0]
	for i := range tc.a {
		sumA += tc.a[i]
		sumB += tc.b[i]
		if tc.a[i] < minA {
			minA = tc.a[i]
		}
		if tc.b[i] < minB {
			minB = tc.b[i]
		}
	}
	costRows := sumA + n*minB
	costCols := sumB + n*minA
	if costCols < costRows {
		return fmt.Sprintf("%d", costCols)
	}
	return fmt.Sprintf("%d", costRows)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 102)
	cases = append(cases, testCase{a: []int64{1}, b: []int64{1}})
	cases = append(cases, testCase{a: []int64{5, 2}, b: []int64{3, 4}})
	for len(cases) < 102 {
		n := rng.Intn(10) + 1
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Int63n(100) + 1
			b[i] = rng.Int63n(100) + 1
		}
		cases = append(cases, testCase{a: a, b: b})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genCases()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", len(tc.a)))
		for j := range tc.a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", tc.a[j]))
		}
		sb.WriteByte('\n')
		for j := range tc.b {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", tc.b[j]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		want := expected(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

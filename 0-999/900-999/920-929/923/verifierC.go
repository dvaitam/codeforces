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

type testC struct {
	n int
	A []int
	P []int
}

func solveC(tc testC) []int {
	n := tc.n
	used := make([]bool, n)
	res := make([]int, n)
	for i, a := range tc.A {
		bestVal := int(^uint(0) >> 1)
		bestIdx := -1
		for j, p := range tc.P {
			if used[j] {
				continue
			}
			v := a ^ p
			if v < bestVal {
				bestVal = v
				bestIdx = j
			}
		}
		res[i] = bestVal
		used[bestIdx] = true
	}
	return res
}

func genTest(rng *rand.Rand) testC {
	n := rng.Intn(8) + 1
	A := make([]int, n)
	P := make([]int, n)
	for i := 0; i < n; i++ {
		A[i] = rng.Intn(256)
	}
	for i := 0; i < n; i++ {
		P[i] = rng.Intn(256)
	}
	return testC{n, A, P}
}

func formatInput(tc testC) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.n)
	for i, v := range tc.A {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i, v := range tc.P {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genTest(rng)
		input := formatInput(tc)
		expVals := solveC(tc)
		var exp strings.Builder
		for j, v := range expVals {
			if j > 0 {
				exp.WriteByte(' ')
			}
			fmt.Fprintf(&exp, "%d", v)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp.String() {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp.String(), got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

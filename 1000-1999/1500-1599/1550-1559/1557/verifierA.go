package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testA struct {
	n   int
	arr []float64
}

func solveA(tc testA) float64 {
	maxv := tc.arr[0]
	sum := 0.0
	for _, v := range tc.arr {
		if v > maxv {
			maxv = v
		}
		sum += v
	}
	sum -= maxv
	return sum/float64(tc.n-1) + maxv
}

func generateA(rng *rand.Rand) testA {
	n := rng.Intn(9) + 2
	arr := make([]float64, n)
	for i := 0; i < n; i++ {
		arr[i] = float64(rng.Intn(200) - 100)
	}
	return testA{n, arr}
}

func runCase(bin string, tc testA) (string, error) {
	var input strings.Builder
	fmt.Fprintf(&input, "1\n%d\n", tc.n)
	for i, v := range tc.arr {
		if i > 0 {
			input.WriteByte(' ')
		}
		fmt.Fprintf(&input, "%d", int(v))
	}
	input.WriteByte('\n')

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testA, 0, 102)
	// two edge cases
	cases = append(cases, testA{2, []float64{1, 2}})
	cases = append(cases, testA{3, []float64{-7, -6, -6}})
	for i := 0; i < 100; i++ {
		cases = append(cases, generateA(rng))
	}

	for i, tc := range cases {
		expect := solveA(tc)
		out, err := runCase(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		var got float64
		fmt.Sscan(out, &got)
		if math.Abs(got-expect) > 1e-6*math.Max(1, math.Abs(expect)) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %.8f got %s\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

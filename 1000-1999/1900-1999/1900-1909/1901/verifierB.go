package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

// expected computes the minimum number of teleports for a single test case.
// Formula: sum of positive rises in c (treating c[-1]=0 as baseline), minus 1.
func expected(c []int64) int64 {
	var segments, prev int64
	for _, v := range c {
		if v > prev {
			segments += v - prev
		}
		prev = v
	}
	return segments - 1
}

func genTest(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = rng.Int63n(1000)
	}
	if arr[0] == 0 {
		arr[0] = 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", n)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("%v: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rng := rand.New(rand.NewSource(2))
	for i := 0; i < 1000; i++ {
		tc := genTest(rng)

		// Parse c from the generated test case.
		lines := strings.Split(strings.TrimSpace(tc), "\n")
		var n int
		fmt.Sscan(lines[1], &n)
		c := make([]int64, n)
		fields := strings.Fields(lines[2])
		for j := 0; j < n; j++ {
			fmt.Sscan(fields[j], &c[j])
		}
		exp := expected(c)

		got, err := runBinary(binary, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "binary failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}

		var gotVal int64
		if _, err := fmt.Sscan(got, &gotVal); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: could not parse output %q: %v\n", i+1, got, err)
			os.Exit(1)
		}
		if gotVal != exp {
			fmt.Printf("mismatch on test %d\ninput:\n%s\nexpected: %d\nactual:   %d\n", i+1, tc, exp, gotVal)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

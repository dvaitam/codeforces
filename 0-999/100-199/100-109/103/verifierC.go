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
	input  string
	expect string
}

func cal(n, k, x int64) byte {
	if n&1 == 1 {
		if x == n {
			if k > 0 {
				return 'X'
			}
			return '.'
		}
		n--
		k--
	}
	if x%2 == 0 {
		if x <= n-k*2 {
			return '.'
		}
		return 'X'
	}
	if x < n-(k-n/2)*2 {
		return '.'
	}
	return 'X'
}

func solve(n, k int64, qs []int64) string {
	res := make([]byte, len(qs))
	for i, x := range qs {
		res[i] = cal(n, k, x)
	}
	return string(res) + "\n"
}

func buildCase(n, k int64, qs []int64) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, k, len(qs))
	for _, x := range qs {
		fmt.Fprintf(&sb, "%d ", x)
	}
	sb.WriteByte('\n')
	return testCase{input: sb.String(), expect: solve(n, k, qs)}
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Int63n(40) + 1
	k := rng.Int63n(n + 1)
	p := rng.Intn(20) + 1
	qs := make([]int64, p)
	for i := range qs {
		qs[i] = rng.Int63n(n) + 1
	}
	return buildCase(n, k, qs)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := out.String()
	if strings.TrimSpace(got) != strings.TrimSpace(tc.expect) {
		return fmt.Errorf("expected %q got %q", tc.expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	cases = append(cases, buildCase(1, 1, []int64{1}))
	cases = append(cases, buildCase(6, 3, []int64{1, 2, 3, 4, 5, 6}))

	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

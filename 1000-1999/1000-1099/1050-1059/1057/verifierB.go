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
	arr []int
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

func expected(arr []int) int {
	n := len(arr)
	pref := make([]int, n+1)
	for i, v := range arr {
		pref[i+1] = pref[i] + v
	}
	best := 0
	for i := 1; i <= n; i++ {
		for j := 0; j < i; j++ {
			sum := pref[i] - pref[j]
			if sum > (i-j)*100 {
				if i-j > best {
					best = i - j
				}
			}
		}
	}
	return best
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(40) + 1 // 1..40
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(201) // 0..200
	}
	return testCase{arr}
}

func runCase(bin string, tc testCase) error {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tc.arr))
	for i, v := range tc.arr {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", v)
	}
	b.WriteByte('\n')

	got, err := run(bin, b.String())
	if err != nil {
		return err
	}
	exp := fmt.Sprintf("%d", expected(tc.arr))
	if strings.TrimSpace(got) != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := make([]testCase, 0, 100)
	// deterministic cases
	cases = append(cases, testCase{arr: []int{50}})
	cases = append(cases, testCase{arr: []int{200, 1}})
	cases = append(cases, testCase{arr: []int{0, 0, 0}})
	for len(cases) < 100 {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

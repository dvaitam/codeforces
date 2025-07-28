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
	input    string
	expected string
}

func solveCase(n, a, b int64) string {
	if a == 1 {
		if (n-1)%b == 0 {
			return "Yes"
		}
		return "No"
	}
	x := int64(1)
	for x <= n {
		if (n-x)%b == 0 {
			return "Yes"
		}
		if x > n/a {
			break
		}
		x *= a
	}
	return "No"
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Int63n(1_000_000_000) + 1
	a := rng.Int63n(10) + 1
	b := rng.Int63n(10) + 1
	input := fmt.Sprintf("1\n%d %d %d\n", n, a, b)
	return testCase{input: input, expected: solveCase(n, a, b)}
}

func runCase(bin string, tc testCase) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	result := strings.TrimSpace(out.String())
	if !strings.EqualFold(result, tc.expected) {
		return fmt.Errorf("expected %s got %s", tc.expected, result)
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

	cases := []testCase{
		{input: "1\n1 1 1\n", expected: solveCase(1, 1, 1)},
		{input: "1\n10 2 3\n", expected: solveCase(10, 2, 3)},
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}

	for idx, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

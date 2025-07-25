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

func runSolution(bin, input string) (string, error) {
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

func costOf(s string) int {
	cnt := make([]int, 26)
	for i := 0; i < len(s); i++ {
		c := s[i] - 'a'
		if c < 26 {
			cnt[c]++
		}
	}
	total := 0
	for _, x := range cnt {
		total += x * (x - 1) / 2
	}
	return total
}

func verifyA(input, output string) error {
	input = strings.TrimSpace(input)
	var k int
	if _, err := fmt.Sscan(input, &k); err != nil {
		return fmt.Errorf("bad input: %v", err)
	}
	output = strings.TrimSpace(output)
	if len(output) == 0 {
		return fmt.Errorf("empty output")
	}
	if len(output) > 100000 {
		return fmt.Errorf("output too long: %d", len(output))
	}
	if costOf(output) != k {
		return fmt.Errorf("expected cost %d got %d", k, costOf(output))
	}
	return nil
}

func generateCase(rng *rand.Rand) string {
	k := rng.Intn(100001)
	return fmt.Sprintf("%d\n", k)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{"0\n", "1\n", "2\n", "3\n", "10\n", "100\n", "99999\n", "100000\n"}
	for len(cases) < 100 {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		out, err := runSolution(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if err := verifyA(tc, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, tc, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

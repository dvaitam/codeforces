package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func solve(d, k, a, b, t int64) int64 {
	if d <= k {
		return d * a
	}
	ans := k * a
	d -= k
	full := d / k
	costDrive := t + k*a
	costWalk := k * b
	if costDrive < costWalk {
		ans += full * costDrive
	} else {
		ans += full * costWalk
		return ans + (d%k)*b
	}
	rem := d % k
	if rem > 0 {
		ans += min(t+rem*a, rem*b)
	}
	return ans
}

func buildCase(d, k, a, b, t int64) testCase {
	input := fmt.Sprintf("%d %d %d %d %d\n", d, k, a, b, t)
	expect := strconv.FormatInt(solve(d, k, a, b, t), 10)
	return testCase{input: input, expected: expect}
}

func genRandomCase(rng *rand.Rand) testCase {
	d := int64(rng.Intn(1000) + 1)
	k := int64(rng.Intn(20) + 1)
	a := int64(rng.Intn(5) + 1)
	b := a + int64(rng.Intn(5)+1)
	t := int64(rng.Intn(10) + 1)
	return buildCase(d, k, a, b, t)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	cases = append(cases, buildCase(5, 2, 1, 3, 5))
	cases = append(cases, buildCase(10, 3, 1, 4, 2))
	for i := 0; i < 100; i++ {
		cases = append(cases, genRandomCase(rng))
	}

	for idx, tc := range cases {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

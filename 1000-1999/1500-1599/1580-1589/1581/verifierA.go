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
	n int
}

const mod int64 = 1000000007

func expectedA(n int) string {
	fac := int64(1)
	for i := 2; i <= 2*n; i++ {
		fac = fac * int64(i) % mod
	}
	ans := fac * ((mod + 1) / 2) % mod
	return fmt.Sprint(ans)
}

func buildInputA(tc testCaseA) string {
	return fmt.Sprintf("1\n%d\n", tc.n)
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd = exec.CommandContext(ctx, cmd.Path, cmd.Args[1:]...)
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func randomCaseA(rng *rand.Rand) testCaseA {
	return testCaseA{n: rng.Intn(100000) + 1}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCaseA
	// deterministic cases
	cases = append(cases, testCaseA{n: 1})
	cases = append(cases, testCaseA{n: 2})
	cases = append(cases, testCaseA{n: 3})
	cases = append(cases, testCaseA{n: 5})
	cases = append(cases, testCaseA{n: 10})
	cases = append(cases, testCaseA{n: 100000})
	cases = append(cases, testCaseA{n: 99999})
	cases = append(cases, testCaseA{n: 54321})

	for len(cases) < 110 {
		cases = append(cases, randomCaseA(rng))
	}

	for i, tc := range cases {
		input := buildInputA(tc)
		expect := expectedA(tc.n)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

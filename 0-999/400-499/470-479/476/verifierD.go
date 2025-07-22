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

func runCandidate(bin, input string) (string, error) {
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

type testCase struct {
	n int
	k int
}

func expected(tc testCase) string {
	var sb strings.Builder
	m := (6*tc.n - 1) * tc.k
	sb.WriteString(fmt.Sprintf("%d\n", m))
	for i := 0; i < tc.n; i++ {
		x := 1 + 6*i
		a := x * tc.k
		b := (x + 2) * tc.k
		c := (x + 4) * tc.k
		d := (x + 1) * tc.k
		if (x+1)%3 == 0 {
			d = (x + 3) * tc.k
		}
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", a, b, c, d))
	}
	return strings.TrimSpace(sb.String())
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 1
	k := rng.Intn(20) + 1
	return testCase{n, k}
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%d %d\n", tc.n, tc.k)
	exp := expected(tc)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) != exp {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{{1, 1}, {2, 3}}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

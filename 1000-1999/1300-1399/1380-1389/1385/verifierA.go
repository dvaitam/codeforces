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
	x, y, z  int64
	input    string
	expected string
}

func solveCase(x, y, z int64) (bool, int64, int64, int64) {
	if x == y && y == z {
		return true, x, y, z
	}
	if x == y {
		if x >= z {
			return true, x, z, 1
		}
		return false, 0, 0, 0
	}
	if x == z {
		if x >= y {
			return true, y, 1, x
		}
		return false, 0, 0, 0
	}
	if y == z {
		if y >= x {
			return true, 1, x, y
		}
		return false, 0, 0, 0
	}
	return false, 0, 0, 0
}

func buildCase(x, y, z int64) testCase {
	ok, a, b, c := solveCase(x, y, z)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d %d\n", x, y, z))
	var exp strings.Builder
	if ok {
		exp.WriteString("YES\n")
		exp.WriteString(fmt.Sprintf("%d %d %d\n", a, b, c))
	} else {
		exp.WriteString("NO\n")
	}
	return testCase{x: x, y: y, z: z, input: sb.String(), expected: exp.String()}
}

func randomCase(rng *rand.Rand) testCase {
	x := rng.Int63n(100) + 1
	y := rng.Int63n(100) + 1
	z := rng.Int63n(100) + 1
	return buildCase(x, y, z)
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
	if strings.TrimSpace(got) != strings.TrimSpace(tc.expected) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(tc.expected), strings.TrimSpace(got))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		buildCase(1, 1, 1),
		buildCase(3, 3, 2),
		buildCase(10, 20, 30),
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

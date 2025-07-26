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

type testCaseC struct {
	input    string
	expected string
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solveAngle(angle int) int {
	pgcd := gcd(360, angle*2)
	if pgcd == 1 {
		return 1
	}
	r := 360 / pgcd
	if angle < 90 {
		return r
	}
	m := 2 * angle
	if 360-angle*2 < m {
		m = 360 - angle*2
	}
	if m == 360/r {
		r *= 2
	}
	return r
}

func generateCaseC(rng *rand.Rand) testCaseC {
	t := rng.Intn(5) + 1
	var in strings.Builder
	var out strings.Builder
	in.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		ang := rng.Intn(179) + 1
		in.WriteString(fmt.Sprintf("%d\n", ang))
		out.WriteString(fmt.Sprintf("%d\n", solveAngle(ang)))
	}
	return testCaseC{input: in.String(), expected: out.String()}
}

func runCaseC(bin string, tc testCaseC) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(tc.expected)
	if got != exp {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	cases := []testCaseC{{input: "1\n1\n", expected: fmt.Sprintf("%d\n", solveAngle(1))}}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCaseC(rng))
	}
	for i, tc := range cases {
		if err := runCaseC(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

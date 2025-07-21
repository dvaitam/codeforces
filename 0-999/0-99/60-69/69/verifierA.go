package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseA struct {
	n        int
	vecs     [][3]int
	expected string
}

func generateTests() []testCaseA {
	rng := rand.New(rand.NewSource(1))
	cases := make([]testCaseA, 100)
	for i := range cases {
		n := rng.Intn(100) + 1
		vecs := make([][3]int, n)
		sx, sy, sz := 0, 0, 0
		for j := 0; j < n; j++ {
			x := rng.Intn(201) - 100
			y := rng.Intn(201) - 100
			z := rng.Intn(201) - 100
			vecs[j] = [3]int{x, y, z}
			sx += x
			sy += y
			sz += z
		}
		exp := "NO"
		if sx == 0 && sy == 0 && sz == 0 {
			exp = "YES"
		}
		cases[i] = testCaseA{n: n, vecs: vecs, expected: exp}
	}
	return cases
}

func run(bin string, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for _, v := range tc.vecs {
			fmt.Fprintf(&sb, "%d %d %d\n", v[0], v[1], v[2])
		}
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

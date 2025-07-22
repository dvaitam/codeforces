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

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveCase(x1, y1, x2, y2 int, lines [][3]int) string {
	cnt := 0
	for _, line := range lines {
		a, b, c := line[0], line[1], line[2]
		f1 := a*x1 + b*y1 + c
		f2 := a*x2 + b*y2 + c
		if (f1 < 0 && f2 > 0) || (f1 > 0 && f2 < 0) {
			cnt++
		}
	}
	return fmt.Sprintf("%d", cnt)
}

func generateCase(rng *rand.Rand) testCase {
	x1 := rng.Intn(21) - 10
	y1 := rng.Intn(21) - 10
	x2 := rng.Intn(21) - 10
	y2 := rng.Intn(21) - 10
	n := rng.Intn(9) + 1 // 1..9
	lines := make([][3]int, n)
	for i := 0; i < n; i++ {
		a := rng.Intn(21) - 10
		b := rng.Intn(21) - 10
		if a == 0 && b == 0 {
			a = 1
		}
		c := rng.Intn(21) - 10
		// ensure points are not on line
		for a*x1+b*y1+c == 0 || a*x2+b*y2+c == 0 {
			a = rng.Intn(21) - 10
			b = rng.Intn(21) - 10
			if a == 0 && b == 0 {
				a = 1
			}
			c = rng.Intn(21) - 10
		}
		lines[i] = [3]int{a, b, c}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", x1, y1))
	sb.WriteString(fmt.Sprintf("%d %d\n", x2, y2))
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, line := range lines {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", line[0], line[1], line[2]))
	}
	return testCase{input: sb.String(), expected: solveCase(x1, y1, x2, y2, lines)}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		{input: "0 0\n1 1\n1\n1 0 -1\n", expected: "1"},
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected: %s\nfound: %s\n", i+1, tc.input, tc.expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

type testCase struct {
	input    string
	expected string
}

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func solve(a, b, c, d int, sizes []int) string {
	width := b
	if d > width {
		width = d
	}
	height := a + c
	grid := make([][]byte, width)
	for i := 0; i < width; i++ {
		row := make([]byte, height)
		for j := range row {
			row[j] = '.'
		}
		grid[i] = row
	}
	x, y, dx := 0, 0, 1
	if a%2 == 1 {
		x = b - 1
		dx = -1
	}
	next := func() {
		nx := x + dx
		if nx < 0 {
			dx = 1
			x = 0
			y++
		} else if y < a && nx >= b {
			dx = -1
			x = b - 1
			y++
		} else if y >= a && nx >= d {
			dx = -1
			x = d - 1
			y++
		} else {
			x = nx
		}
	}
	for i := 0; i < len(sizes); i++ {
		for sizes[i] > 0 {
			grid[x][y] = byte('a' + i)
			sizes[i]--
			next()
		}
	}
	var out bytes.Buffer
	out.WriteString("YES\n")
	for i := 0; i < width; i++ {
		out.Write(grid[i])
		out.WriteByte('\n')
	}
	return strings.TrimSpace(out.String())
}

func generateCases() []testCase {
	rand.Seed(4)
	cases := make([]testCase, 100)
	for t := 0; t < 100; t++ {
		a := rand.Intn(4) + 1
		b := rand.Intn(4) + 1
		c := rand.Intn(4) + 1
		d := rand.Intn(4) + 1
		n := rand.Intn(5) + 1
		area := a*b + c*d
		if n > 26 {
			n = 26
		}
		sizes := make([]int, n)
		rem := area
		for i := 0; i < n; i++ {
			if i == n-1 {
				sizes[i] = rem
			} else {
				maxv := rem - (n - i - 1)
				val := rand.Intn(maxv) + 1
				sizes[i] = val
				rem -= val
			}
		}
		buf := bytes.Buffer{}
		fmt.Fprintf(&buf, "%d %d %d %d %d\n", a, b, c, d, n)
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(&buf, " ")
			}
			fmt.Fprint(&buf, sizes[i])
		}
		buf.WriteByte('\n')
		expected := solve(a, b, c, d, append([]int(nil), sizes...))
		cases[t] = testCase{input: buf.String(), expected: expected}
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierD.go <binary>")
		os.Exit(1)
	}
	cases := generateCases()
	for i, tc := range cases {
		out, err := runBinary(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", i+1, tc.input, tc.expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

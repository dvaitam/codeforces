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

func solve(x, y, x0, y0 int, s string) []int {
	visited := make([][]bool, x)
	for i := range visited {
		visited[i] = make([]bool, y)
	}
	cx, cy := x0-1, y0-1
	visited[cx][cy] = true
	visitedCount := 1
	counts := make([]int, len(s)+1)
	counts[0] = 1
	for i, ch := range s {
		switch ch {
		case 'L':
			if cy > 0 {
				cy--
			}
		case 'R':
			if cy+1 < y {
				cy++
			}
		case 'U':
			if cx > 0 {
				cx--
			}
		case 'D':
			if cx+1 < x {
				cx++
			}
		}
		if !visited[cx][cy] {
			visited[cx][cy] = true
			counts[i+1]++
			visitedCount++
		}
	}
	counts[len(s)] += x*y - visitedCount
	return counts
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateCase(rng *rand.Rand) (string, string) {
	x := rng.Intn(10) + 1
	y := rng.Intn(10) + 1
	x0 := rng.Intn(x) + 1
	y0 := rng.Intn(y) + 1
	length := rng.Intn(50) + 1
	moves := []byte{'L', 'R', 'U', 'D'}
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(moves[rng.Intn(4)])
	}
	s := sb.String()
	counts := solve(x, y, x0, y0, s)
	outParts := make([]string, len(counts))
	for i, v := range counts {
		outParts[i] = fmt.Sprint(v)
	}
	expected := strings.Join(outParts, " ")
	input := fmt.Sprintf("%d %d %d %d\n%s\n", x, y, x0, y0, s)
	return input, expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	special := []struct {
		x, y, x0, y0 int
		s            string
	}{
		{1, 1, 1, 1, "L"},
		{2, 2, 1, 1, "URDL"},
		{3, 3, 2, 2, "UUURRRDDDLLL"},
		{7, 8, 4, 5, "LURDLURD"},
		{10, 10, 5, 5, strings.Repeat("U", 20)},
		{5, 5, 5, 5, "LLLLRRRRUUUUDDDD"},
		{10, 10, 1, 1, strings.Repeat("L", 20)},
		{10, 10, 10, 10, strings.Repeat("R", 20)},
		{2, 3, 1, 2, "UDLR"},
		{4, 4, 2, 3, "LRUDLRUD"},
		{3, 3, 3, 3, "DDDLLLUUU"},
	}
	for i, tc := range special {
		input := fmt.Sprintf("%d %d %d %d\n%s\n", tc.x, tc.y, tc.x0, tc.y0, tc.s)
		expectedCounts := solve(tc.x, tc.y, tc.x0, tc.y0, tc.s)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("special case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		expectedParts := make([]string, len(expectedCounts))
		for j, v := range expectedCounts {
			expectedParts[j] = fmt.Sprint(v)
		}
		expected := strings.Join(expectedParts, " ")
		if out != expected {
			fmt.Printf("special case %d failed: expected %s got %s\n", i+1, expected, out)
			os.Exit(1)
		}
	}

	for i := 0; i < 90; i++ {
		input, expected := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != expected {
			fmt.Printf("case %d failed:\ninput:%sexpected %s got %s\n", i+1, input, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

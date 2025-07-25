package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	input    string
	expected string
}

func solve(points [][2]int64) int {
	x := [3]int64{points[0][0], points[1][0], points[2][0]}
	y := [3]int64{points[0][1], points[1][1], points[2][1]}
	if (x[0] == x[1] && x[1] == x[2]) || (y[0] == y[1] && y[1] == y[2]) {
		return 1
	}
	for k := 0; k < 3; k++ {
		i := (k + 1) % 3
		j := (k + 2) % 3
		if x[i] == x[k] && y[j] == y[k] {
			return 2
		}
		if y[i] == y[k] && x[j] == x[k] {
			return 2
		}
	}
	if x[0] == x[1] || x[0] == x[2] || x[1] == x[2] || y[0] == y[1] || y[0] == y[2] || y[1] == y[2] {
		return 3
	}
	return 4
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(45))
	tests := make([]test, 100)
	for i := 0; i < 100; i++ {
		pts := make([][2]int64, 3)
		var sb strings.Builder
		for j := 0; j < 3; j++ {
			x := int64(rng.Intn(21) - 10)
			y := int64(rng.Intn(21) - 10)
			pts[j] = [2]int64{x, y}
			fmt.Fprintf(&sb, "%d %d\n", x, y)
		}
		ans := solve(pts)
		tests[i] = test{sb.String(), fmt.Sprintf("%d", ans)}
	}
	return tests
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := run(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

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

func solve(n int, x1, y1, x2, y2 int, points [][2]int) int64 {
	d1 := make([]int64, n)
	d2 := make([]int64, n)
	for i := 0; i < n; i++ {
		dx1 := int64(points[i][0] - x1)
		dy1 := int64(points[i][1] - y1)
		d1[i] = dx1*dx1 + dy1*dy1
		dx2 := int64(points[i][0] - x2)
		dy2 := int64(points[i][1] - y2)
		d2[i] = dx2*dx2 + dy2*dy2
	}
	const inf int64 = 1 << 62
	ans := inf
	for i := 0; i < n; i++ {
		r1 := d1[i]
		var r2 int64
		for j := 0; j < n; j++ {
			if d1[j] > r1 {
				if d2[j] > r2 {
					r2 = d2[j]
				}
			}
		}
		if r1+r2 < ans {
			ans = r1 + r2
		}
	}
	var r2 int64
	for i := 0; i < n; i++ {
		if d2[i] > r2 {
			r2 = d2[i]
		}
	}
	if r2 < ans {
		ans = r2
	}
	return ans
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(44))
	var tests []test
	for len(tests) < 100 {
		n := rng.Intn(10) + 1
		x1 := rng.Intn(21) - 10
		y1 := rng.Intn(21) - 10
		x2 := rng.Intn(21) - 10
		y2 := rng.Intn(21) - 10
		points := make([][2]int, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d %d %d\n", n, x1, y1, x2, y2)
		for i := 0; i < n; i++ {
			px := rng.Intn(21) - 10
			py := rng.Intn(21) - 10
			points[i] = [2]int{px, py}
			fmt.Fprintf(&sb, "%d %d\n", px, py)
		}
		ans := solve(n, x1, y1, x2, y2, points)
		tests = append(tests, test{sb.String(), fmt.Sprintf("%d", ans)})
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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

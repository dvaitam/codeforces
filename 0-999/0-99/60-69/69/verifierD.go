package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseD struct {
	x, y     int
	n, d     int
	dx, dy   []int
	expected string
}

var (
	d2     int
	offset int
	memo   map[int]bool
	vis    map[int]bool
)

func encode(x, y, rA, rB, turn int) int {
	xi := x + offset
	yi := y + offset
	return (((xi << 9) | yi) << 3) | (rA << 2) | (rB << 1) | turn
}

func dfs(x, y, rA, rB, turn int, moves [][2]int) bool {
	key := encode(x, y, rA, rB, turn)
	if vis[key] {
		return memo[key]
	}
	vis[key] = true
	win := false
	for _, mv := range moves {
		nx := x + mv[0]
		ny := y + mv[1]
		if nx*nx+ny*ny > d2 {
			continue
		}
		if !dfs(nx, ny, rA, rB, 1-turn, moves) {
			win = true
			break
		}
	}
	if !win {
		if turn == 0 && rA == 1 {
			if !dfs(y, x, 0, rB, 1-turn, moves) {
				win = true
			}
		} else if turn == 1 && rB == 1 {
			if !dfs(y, x, rA, 0, 1-turn, moves) {
				win = true
			}
		}
	}
	memo[key] = win
	return win
}

func solveCase(tc testCaseD) string {
	offset = tc.d
	d2 = tc.d * tc.d
	memo = make(map[int]bool)
	vis = make(map[int]bool)
	moves := make([][2]int, tc.n)
	for i := 0; i < tc.n; i++ {
		moves[i] = [2]int{tc.dx[i], tc.dy[i]}
	}
	if dfs(tc.x, tc.y, 1, 1, 0, moves) {
		return "Anton"
	}
	return "Dasha"
}

func generateTests() []testCaseD {
	rng := rand.New(rand.NewSource(4))
	cases := make([]testCaseD, 100)
	for i := range cases {
		d := rng.Intn(5) + 2
		x := rng.Intn(2*d) - d
		y := rng.Intn(2*d) - d
		n := rng.Intn(3) + 1
		dx := make([]int, n)
		dy := make([]int, n)
		for j := 0; j < n; j++ {
			dx[j] = rng.Intn(d)
			dy[j] = rng.Intn(d)
			if dx[j] == 0 && dy[j] == 0 {
				dx[j] = 1
			}
		}
		tc := testCaseD{x: x, y: y, n: n, d: d, dx: dx, dy: dy}
		tc.expected = solveCase(tc)
		cases[i] = tc
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d %d\n", tc.x, tc.y, tc.n, tc.d)
		for j := 0; j < tc.n; j++ {
			fmt.Fprintf(&sb, "%d %d\n", tc.dx[j], tc.dy[j])
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

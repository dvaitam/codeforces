package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesC.txt.
const testcasesCData = `
3 100 100 50
4 50 50 50 50
4 50 50 100 100
6 50 50 50 50 50 50
6 100 100 100 100 100 100
7 50 50 50 100 100 100 100
`

const MOD = 1000000007

type testCase struct {
	weights []int
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesCData), "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", idx+1, err)
		}
		if len(fields)-1 != n {
			return nil, fmt.Errorf("line %d: expected %d weights, got %d", idx+1, n, len(fields)-1)
		}
		w := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(fields[i+1])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse weight %d: %w", idx+1, i, err)
			}
			w[i] = val
		}
		cases = append(cases, testCase{weights: w})
	}
	return cases, nil
}

func precomputeC(maxN int) [][]int {
	C := make([][]int, maxN+1)
	for i := 0; i <= maxN; i++ {
		C[i] = make([]int, i+1)
		C[i][0] = 1
		for j := 1; j <= i; j++ {
			if j == i {
				C[i][j] = 1
			} else {
				C[i][j] = (C[i-1][j-1] + C[i-1][j]) % MOD
			}
		}
	}
	return C
}

// solve mirrors 295C.go for a given testcase and returns (moves, ways).
func solve(tc testCase) (int, int) {
	n := len(tc.weights)
	total50, total100 := 0, 0
	for _, w := range tc.weights {
		if w == 50 {
			total50++
		} else {
			total100++
		}
	}
	C := precomputeC(n)

	type move struct{ x, y int }
	moves := make([]move, 0)
	// maximum k is 200 (50*4), but from original problem, bag capacity k=200
	k := 200
	maxY := k / 100
	for y := 0; y <= maxY; y++ {
		rem := k - 100*y
		maxX := rem / 50
		for x := 0; x <= maxX; x++ {
			if x == 0 && y == 0 {
				continue
			}
			moves = append(moves, move{x, y})
		}
	}

	const INF = 1 << 30
	dist := make([][][]int, total50+1)
	ways := make([][][]int, total50+1)
	for a := 0; a <= total50; a++ {
		dist[a] = make([][]int, total100+1)
		ways[a] = make([][]int, total100+1)
		for b := 0; b <= total100; b++ {
			dist[a][b] = []int{INF, INF}
			ways[a][b] = []int{0, 0}
		}
	}
	type state struct{ a, b, side int }
	queue := make([]state, 0, (total50+1)*(total100+1)*2)
	dist[total50][total100][0] = 0
	ways[total50][total100][0] = 1
	queue = append(queue, state{total50, total100, 0})
	head := 0
	for head < len(queue) {
		cur := queue[head]
		head++
		d := dist[cur.a][cur.b][cur.side]
		wcur := ways[cur.a][cur.b][cur.side]
		for _, mv := range moves {
			x, y := mv.x, mv.y
			var na, nb, nside int
			var c50, c100 int
			if cur.side == 0 {
				if x > cur.a || y > cur.b {
					continue
				}
				na = cur.a - x
				nb = cur.b - y
				nside = 1
				c50 = cur.a
				c100 = cur.b
			} else {
				right50 := total50 - cur.a
				right100 := total100 - cur.b
				if x > right50 || y > right100 {
					continue
				}
				na = cur.a + x
				nb = cur.b + y
				nside = 0
				c50 = right50
				c100 = right100
			}
			nd := d + 1
			addWays := int64(C[c50][x]) * int64(C[c100][y]) % MOD
			if dist[na][nb][nside] > nd {
				dist[na][nb][nside] = nd
				ways[na][nb][nside] = int((int64(wcur) * addWays) % MOD)
				queue = append(queue, state{na, nb, nside})
			} else if dist[na][nb][nside] == nd {
				ways[na][nb][nside] = int((int64(ways[na][nb][nside]) + int64(wcur)*addWays) % MOD)
			}
		}
	}
	if dist[0][0][1] >= INF {
		return -1, 0
	}
	return dist[0][0][1], ways[0][0][1]
}

func runCandidate(bin string, tc testCase) (int, int, error) {
	var input strings.Builder
	fmt.Fprintf(&input, "%d %d\n", len(tc.weights), 200)
	for i, w := range tc.weights {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprint(w))
	}
	input.WriteByte('\n')

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return 0, 0, fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("expected 2 outputs, got %d", len(fields))
	}
	d, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, fmt.Errorf("parse distance: %v", err)
	}
	w, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, fmt.Errorf("parse ways: %v", err)
	}
	return d, w, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		expD, expW := solve(tc)
		gotD, gotW, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if gotD != expD || gotW != expW {
			fmt.Fprintf(os.Stderr, "case %d failed: expected (%d %d) got (%d %d)\n", i+1, expD, expW, gotD, gotW)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

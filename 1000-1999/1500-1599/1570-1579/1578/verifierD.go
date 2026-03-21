package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// dragonOracle builds a lookup table mapping (x,y) -> "curve pos" by simulating
// the four infinite-order dragon curves for the first maxSeg segments each.
func dragonOracle(maxSeg int) map[[2]int][2]int {
	// Directions: 0=NE(+1,+1), 1=NW(-1,+1), 2=SW(-1,-1), 3=SE(+1,-1)
	dx := [4]int{1, -1, -1, 1}
	dy := [4]int{1, 1, -1, -1}

	grid := make(map[[2]int][2]int) // (x,y) -> [curve, position]

	for curve := 0; curve < 4; curve++ {
		cx, cy := 0, 0
		dir := curve
		for seg := 1; seg <= maxSeg; seg++ {
			nx, ny := cx+dx[dir], cy+dy[dir]
			sx := cx
			if dx[dir] < 0 {
				sx = nx
			}
			sy := cy
			if dy[dir] < 0 {
				sy = ny
			}
			k := [2]int{sx, sy}
			if _, ok := grid[k]; !ok {
				grid[k] = [2]int{curve + 1, seg}
			}
			cx, cy = nx, ny
			if seg < maxSeg {
				tz := bits.TrailingZeros(uint(seg))
				bit := (seg >> (tz + 1)) & 1
				if bit == 0 {
					dir = (dir + 1) % 4 // left
				} else {
					dir = (dir + 3) % 4 // right
				}
			}
		}
	}

	return grid
}

func genCase(r *rand.Rand, coordRange int) string {
	n := r.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		x := r.Intn(2*coordRange+1) - coordRange
		y := r.Intn(2*coordRange+1) - coordRange
		sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
	}
	return sb.String()
}

func oracleAnswer(grid map[[2]int][2]int, input string) (string, error) {
	fields := strings.Fields(input)
	n, _ := strconv.Atoi(fields[0])
	var results []string
	idx := 1
	for i := 0; i < n; i++ {
		x, _ := strconv.Atoi(fields[idx])
		y, _ := strconv.Atoi(fields[idx+1])
		idx += 2
		v, ok := grid[[2]int{x, y}]
		if !ok {
			return "", fmt.Errorf("coordinate (%d,%d) not in oracle table", x, y)
		}
		results = append(results, fmt.Sprintf("%d %d", v[0], v[1]))
	}
	return strings.Join(results, "\n"), nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func normalizeOutput(s string) string {
	lines := strings.Split(s, "\n")
	var out []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			out = append(out, l)
		}
	}
	return strings.Join(out, "\n")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	// Coordinate range for random tests; simulate enough segments to cover it.
	const coordRange = 5
	const maxSeg = 512
	grid := dragonOracle(maxSeg)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 20; i++ {
		input := genCase(rng, coordRange)
		expect, err := oracleAnswer(grid, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if normalizeOutput(got) != normalizeOutput(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected:\n%s\ngot:\n%s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 20 tests passed")
}

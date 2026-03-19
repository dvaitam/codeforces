package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// For small n, m (<=10) and small k (<=5), we can directly compute pile sizes.
func computeGame(n, m, k int, rawCuts [][4]int) (int, map[int]int, map[int]int) {
	// Vertical lines: x = 1..n-1, each has m unit segments (from y=0 to y=m)
	// A vertical cut at (x, y1, x, y2) marks segments [min(y1,y2), max(y1,y2)] as cut on line x.
	// Horizontal lines: y = 1..m-1, each has n unit segments (from x=0 to x=n)
	// A horizontal cut at (x1, y, x2, y) marks segments [min(x1,x2), max(x1,x2)] as cut on line y.

	vSegs := make(map[int][][2]int) // x -> list of [y1,y2] intervals
	hSegs := make(map[int][][2]int) // y -> list of [x1,x2] intervals

	for _, c := range rawCuts {
		x1, y1, x2, y2 := c[0], c[1], c[2], c[3]
		if x1 == x2 {
			// vertical cut on line x=x1
			a, b := y1, y2
			if a > b {
				a, b = b, a
			}
			vSegs[x1] = append(vSegs[x1], [2]int{a, b})
		} else {
			// horizontal cut on line y=y1
			a, b := x1, x2
			if a > b {
				a, b = b, a
			}
			hSegs[y1] = append(hSegs[y1], [2]int{a, b})
		}
	}

	// Compute pile sizes
	vPiles := make(map[int]int) // x -> pile size
	hPiles := make(map[int]int) // y -> pile size

	xorSum := 0

	// Vertical lines
	for x := 1; x < n; x++ {
		segs := vSegs[x]
		pileSize := computePileSize(m, segs)
		vPiles[x] = pileSize
		xorSum ^= pileSize
	}

	// Horizontal lines
	for y := 1; y < m; y++ {
		segs := hSegs[y]
		pileSize := computePileSize(n, segs)
		hPiles[y] = pileSize
		xorSum ^= pileSize
	}

	return xorSum, vPiles, hPiles
}

func computePileSize(total int, segs [][2]int) int {
	if len(segs) == 0 {
		return total
	}
	// Merge intervals
	sort.Slice(segs, func(i, j int) bool { return segs[i][0] < segs[j][0] })
	merged := [][2]int{segs[0]}
	for _, s := range segs[1:] {
		last := &merged[len(merged)-1]
		if s[0] <= last[1] {
			if s[1] > last[1] {
				last[1] = s[1]
			}
		} else {
			merged = append(merged, s)
		}
	}
	cut := 0
	for _, s := range merged {
		cut += s[1] - s[0]
	}
	return total - cut
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(9) + 2
	m := rng.Intn(9) + 2
	k := rng.Intn(5)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 0; i < k; i++ {
		if rng.Intn(2) == 0 {
			// vertical cut
			x := rng.Intn(n-1) + 1
			y1 := rng.Intn(m-1) + 1
			y2 := rng.Intn(m-1) + 1
			if y1 > y2 {
				y1, y2 = y2, y1
			}
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", x, y1, x, y2))
		} else {
			// horizontal
			y := rng.Intn(m-1) + 1
			x1 := rng.Intn(n-1) + 1
			x2 := rng.Intn(n-1) + 1
			if x1 > x2 {
				x1, x2 = x2, x1
			}
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", x1, y, x2, y))
		}
	}
	return sb.String()
}

func parseInput(input string) (int, int, int, [][4]int) {
	r := strings.NewReader(input)
	var n, m, k int
	fmt.Fscan(r, &n, &m, &k)
	cuts := make([][4]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(r, &cuts[i][0], &cuts[i][1], &cuts[i][2], &cuts[i][3])
	}
	return n, m, k, cuts
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		n, m, _, cuts := parseInput(input)
		xorSum, _, _ := computeGame(n, m, 0, cuts)

		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}

		lines := strings.Fields(got)
		if xorSum == 0 {
			// Second player wins
			if len(lines) == 0 || lines[0] != "SECOND" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected SECOND, got: %s\ninput:%s", i+1, got, input)
				os.Exit(1)
			}
		} else {
			// First player wins
			if len(lines) == 0 || lines[0] != "FIRST" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected FIRST, got: %s\ninput:%s", i+1, got, input)
				os.Exit(1)
			}
			// Parse move: 4 integers
			if len(lines) < 5 {
				fmt.Fprintf(os.Stderr, "case %d failed: FIRST but missing move coordinates, got: %s\ninput:%s", i+1, got, input)
				os.Exit(1)
			}
			var x1, y1, x2, y2 int
			fmt.Sscan(lines[1], &x1)
			fmt.Sscan(lines[2], &y1)
			fmt.Sscan(lines[3], &x2)
			fmt.Sscan(lines[4], &y2)

			// Validate the move is a valid cut
			if x1 == x2 && y1 == y2 {
				fmt.Fprintf(os.Stderr, "case %d failed: zero-length cut\ninput:%s", i+1, input)
				os.Exit(1)
			}
			if x1 != x2 && y1 != y2 {
				fmt.Fprintf(os.Stderr, "case %d failed: diagonal cut\ninput:%s", i+1, input)
				os.Exit(1)
			}

			// Add the move to cuts and verify XOR becomes 0
			newCuts := append(cuts, [4]int{x1, y1, x2, y2})
			newXor, _, _ := computeGame(n, m, 0, newCuts)
			if newXor != 0 {
				fmt.Fprintf(os.Stderr, "case %d failed: move does not lead to losing position (xor=%d)\nmove: %d %d %d %d\ninput:%s", i+1, newXor, x1, y1, x2, y2, input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}

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

// computePileSize returns the total uncut length on a line of the given total
// length, after merging the cut intervals.
func computePileSize(total int, segs [][2]int) int {
	if len(segs) == 0 {
		return total
	}
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

// computeGame returns the Nim-XOR of all pile sizes across all grid lines.
func computeGame(n, m int, rawCuts [][4]int) int {
	vSegs := make(map[int][][2]int)
	hSegs := make(map[int][][2]int)

	for _, c := range rawCuts {
		x1, y1, x2, y2 := c[0], c[1], c[2], c[3]
		if x1 == x2 {
			a, b := y1, y2
			if a > b {
				a, b = b, a
			}
			if a == b {
				continue // zero-length, ignore
			}
			vSegs[x1] = append(vSegs[x1], [2]int{a, b})
		} else {
			a, b := x1, x2
			if a > b {
				a, b = b, a
			}
			if a == b {
				continue // zero-length, ignore
			}
			hSegs[y1] = append(hSegs[y1], [2]int{a, b})
		}
	}

	xorSum := 0
	for x := 1; x < n; x++ {
		xorSum ^= computePileSize(m, vSegs[x])
	}
	for y := 1; y < m; y++ {
		xorSum ^= computePileSize(n, hSegs[y])
	}
	return xorSum
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(9) + 2
	m := rng.Intn(9) + 2
	k := rng.Intn(6)
	var cuts [][4]int
	var sb strings.Builder

	// Generate valid cuts (non-zero length, not on border)
	for len(cuts) < k {
		if rng.Intn(2) == 0 {
			// vertical cut on line x (1..n-1), with y endpoints in [0..m]
			x := rng.Intn(n-1) + 1
			y1 := rng.Intn(m + 1)
			y2 := rng.Intn(m + 1)
			if y1 == y2 {
				continue // skip zero-length cuts
			}
			cuts = append(cuts, [4]int{x, y1, x, y2})
		} else {
			// horizontal cut on line y (1..m-1), with x endpoints in [0..n]
			y := rng.Intn(m-1) + 1
			x1 := rng.Intn(n + 1)
			x2 := rng.Intn(n + 1)
			if x1 == x2 {
				continue // skip zero-length cuts
			}
			cuts = append(cuts, [4]int{x1, y, x2, y})
		}
	}
	k = len(cuts)

	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for _, c := range cuts {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", c[0], c[1], c[2], c[3]))
	}
	return sb.String()
}

func parseInput(input string) (int, int, [][4]int) {
	r := strings.NewReader(input)
	var n, m, k int
	fmt.Fscan(r, &n, &m, &k)
	cuts := make([][4]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(r, &cuts[i][0], &cuts[i][1], &cuts[i][2], &cuts[i][3])
	}
	return n, m, cuts
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		input := generateCase(rng)
		n, m, cuts := parseInput(input)
		xorSum := computeGame(n, m, cuts)

		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}

		lines := strings.Fields(got)
		if xorSum == 0 {
			if len(lines) == 0 || lines[0] != "SECOND" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected SECOND, got: %s\ninput:%s", i+1, got, input)
				os.Exit(1)
			}
		} else {
			if len(lines) == 0 || lines[0] != "FIRST" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected FIRST, got: %s\ninput:%s", i+1, got, input)
				os.Exit(1)
			}
			if len(lines) < 5 {
				fmt.Fprintf(os.Stderr, "case %d failed: FIRST but missing move coordinates, got: %s\ninput:%s", i+1, got, input)
				os.Exit(1)
			}
			var x1, y1, x2, y2 int
			fmt.Sscan(lines[1], &x1)
			fmt.Sscan(lines[2], &y1)
			fmt.Sscan(lines[3], &x2)
			fmt.Sscan(lines[4], &y2)

			if x1 == x2 && y1 == y2 {
				fmt.Fprintf(os.Stderr, "case %d failed: zero-length cut\ninput:%s", i+1, input)
				os.Exit(1)
			}
			if x1 != x2 && y1 != y2 {
				fmt.Fprintf(os.Stderr, "case %d failed: diagonal cut\ninput:%s", i+1, input)
				os.Exit(1)
			}

			// Validate the cut is on a valid internal grid line
			if x1 == x2 {
				if x1 < 1 || x1 >= n {
					fmt.Fprintf(os.Stderr, "case %d failed: vertical cut on invalid line x=%d (n=%d)\ninput:%s", i+1, x1, n, input)
					os.Exit(1)
				}
			} else {
				if y1 < 1 || y1 >= m {
					fmt.Fprintf(os.Stderr, "case %d failed: horizontal cut on invalid line y=%d (m=%d)\ninput:%s", i+1, y1, m, input)
					os.Exit(1)
				}
			}

			// The move must cut at least one new cell (the computeGame check below
			// implicitly verifies this, because if no new cell is cut the XOR
			// would remain non-zero).

			newCuts := append(cuts, [4]int{x1, y1, x2, y2})
			newXor := computeGame(n, m, newCuts)
			if newXor != 0 {
				fmt.Fprintf(os.Stderr, "case %d failed: move does not lead to losing position (xor=%d)\nmove: %d %d %d %d\ninput:%s", i+1, newXor, x1, y1, x2, y2, input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}

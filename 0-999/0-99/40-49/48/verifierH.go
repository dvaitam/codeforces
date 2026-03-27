package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// Tile definitions matching the CF-accepted solution for 48H.
// Each tile is a 2x2 block with specific edge constraints.
// id 0: ".." / ".."  (white, all edges 0)
// id 1: "**" / "**"  (black, all edges 1)
// id 2: "\*" / ".\"  (mixed, top=1,right=1,bottom=0,left=0)
// id 3: "./" / "/*"  (mixed, top=0,right=1,bottom=1,left=0)
// id 4: "/." / "*/"  (mixed, top=0,right=0,bottom=1,left=1)
// id 5: "*\" / "\."  (mixed, top=1,right=0,bottom=0,left=1)

type Tile struct {
	id                       int
	top, right, bottom, left int
	chars                    [2]string
}

var tiles = []Tile{
	{0, 0, 0, 0, 0, [2]string{"..", ".."}},
	{1, 1, 1, 1, 1, [2]string{"**", "**"}},
	{2, 1, 1, 0, 0, [2]string{"\\*", ".\\"}},
	{3, 0, 1, 1, 0, [2]string{"./", "/*"}},
	{4, 0, 0, 1, 1, [2]string{"/.", "*/"}},
	{5, 1, 0, 0, 1, [2]string{"*\\", "\\."}},
}

// solve uses backtracking to produce a valid tiling, matching the CF-accepted approach.
func solve(n, m, a, b, c int) string {
	if a+b+c != n*m {
		return ""
	}

	bottomEdge := make([][]int, n)
	rightEdge := make([][]int, n)
	for i := 0; i < n; i++ {
		bottomEdge[i] = make([]int, m)
		rightEdge[i] = make([]int, m)
	}
	ans := make([]int, n*m)

	var bt func(idx, remB, remW, remM int) bool
	bt = func(idx, remB, remW, remM int) bool {
		if idx == n*m {
			return true
		}
		r := idx / m
		col := idx % m

		top := -1
		if r > 0 {
			top = bottomEdge[r-1][col]
		}
		left := -1
		if col > 0 {
			left = rightEdge[r][col-1]
		}

		order := []int{0, 1, 2, 3, 4, 5}

		score := func(id int) int {
			if id == 0 {
				if remW == 0 {
					return -1
				}
				return remW
			}
			if id == 1 {
				if remB == 0 {
					return -1
				}
				return remB
			}
			if remM == 0 {
				return -1
			}
			return remM
		}

		sort.SliceStable(order, func(i, j int) bool {
			return score(order[i]) > score(order[j])
		})

		for _, id := range order {
			if score(id) == -1 {
				continue
			}
			t := tiles[id]
			if top != -1 && t.top != top {
				continue
			}
			if left != -1 && t.left != left {
				continue
			}

			bottomEdge[r][col] = t.bottom
			rightEdge[r][col] = t.right
			ans[idx] = id

			nB, nW, nM := remB, remW, remM
			if id == 0 {
				nW--
			} else if id == 1 {
				nB--
			} else {
				nM--
			}

			if bt(idx+1, nB, nW, nM) {
				return true
			}
		}
		return false
	}

	if !bt(0, a, b, c) {
		return ""
	}

	var sb strings.Builder
	for i := 0; i < n; i++ {
		for row := 0; row < 2; row++ {
			for j := 0; j < m; j++ {
				sb.WriteString(tiles[ans[i*m+j]].chars[row])
			}
			sb.WriteByte('\n')
		}
	}
	return strings.TrimSpace(sb.String())
}

// verify checks that a candidate output is a valid tiling for the given input.
func verify(n, m, a, b, c int, output string) error {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 2*n {
		return fmt.Errorf("expected %d lines, got %d", 2*n, len(lines))
	}
	for i, line := range lines {
		if len(line) != 2*m {
			return fmt.Errorf("line %d: expected length %d, got %d", i, 2*m, len(line))
		}
	}

	gotA, gotB, gotC := 0, 0, 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			r0 := 2 * i
			c0 := 2 * j
			block := [2]string{
				lines[r0][c0 : c0+2],
				lines[r0+1][c0 : c0+2],
			}
			found := -1
			for _, t := range tiles {
				if t.chars == block {
					found = t.id
					break
				}
			}
			if found == -1 {
				return fmt.Errorf("invalid tile at (%d,%d): %q %q", i, j, block[0], block[1])
			}
			t := tiles[found]
			// Check edge consistency with neighbors
			if i > 0 {
				// Check top edge matches bottom edge of tile above
				aboveIdx := -1
				aboveR := 2*(i-1)
				aboveC := c0
				aboveBlock := [2]string{
					lines[aboveR][aboveC : aboveC+2],
					lines[aboveR+1][aboveC : aboveC+2],
				}
				for _, tt := range tiles {
					if tt.chars == aboveBlock {
						aboveIdx = tt.id
						break
					}
				}
				if aboveIdx >= 0 && tiles[aboveIdx].bottom != t.top {
					return fmt.Errorf("edge mismatch at (%d,%d) top", i, j)
				}
			}
			if j > 0 {
				leftC := 2*(j-1)
				leftBlock := [2]string{
					lines[r0][leftC : leftC+2],
					lines[r0+1][leftC : leftC+2],
				}
				leftIdx := -1
				for _, tt := range tiles {
					if tt.chars == leftBlock {
						leftIdx = tt.id
						break
					}
				}
				if leftIdx >= 0 && tiles[leftIdx].right != t.left {
					return fmt.Errorf("edge mismatch at (%d,%d) left", i, j)
				}
			}
			if found == 0 {
				gotB++
			} else if found == 1 {
				gotA++
			} else {
				gotC++
			}
		}
	}
	if gotA != a || gotB != b || gotC != c {
		return fmt.Errorf("tile counts wrong: got a=%d b=%d c=%d, expected a=%d b=%d c=%d", gotA, gotB, gotC, a, b, c)
	}
	return nil
}

func generateCase(rng *rand.Rand) (int, int, int, int, int) {
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 1
	total := n * m
	a := rng.Intn(total + 1)
	b := rng.Intn(total - a + 1)
	c := total - a - b
	return n, m, a, b, c
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		n, m, a, b, c := generateCase(rng)
		input := fmt.Sprintf("%d %d\n%d %d %d\n", n, m, a, b, c)

		// Check if our reference can solve it (to know if a solution exists)
		refOut := solve(n, m, a, b, c)

		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}

		if refOut == "" {
			// No solution exists; candidate should output nothing useful
			// (The problem guarantees a+b+c=n*m so this shouldn't happen with valid inputs)
			continue
		}

		// Verify the candidate output is a valid tiling
		if err := verify(n, m, a, b, c, got); err != nil {
			fmt.Printf("Test %d failed: %v\nInput:\n%sGot:\n%s\n", i+1, err, input, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

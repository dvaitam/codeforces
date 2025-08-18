package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type UnionFind struct {
	parent []int
}

func NewUnionFind(n int) *UnionFind {
	parent := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
	}
	return &UnionFind{parent}
}

func (uf *UnionFind) find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.find(uf.parent[x])
	}
	return uf.parent[x]
}

func (uf *UnionFind) union(x, y int) bool {
	px, py := uf.find(x), uf.find(y)
	if px == py {
		return false
	}
	uf.parent[py] = px
	return true
}

func expected(n int, edges [][2]int) (int, [][4]int) {
	uf := NewUnionFind(n)
	bad := make([][2]int, 0)

	for _, edge := range edges {
		x, y := edge[0], edge[1]
		if !uf.union(x, y) {
			bad = append(bad, [2]int{x, y})
		}
	}

	roots := make([]int, 0)
	for i := 1; i <= n; i++ {
		if uf.parent[i] == i {
			roots = append(roots, i)
		}
	}

	result := make([][4]int, len(bad))
	for i := 0; i < len(bad); i++ {
		result[i] = [4]int{bad[i][0], bad[i][1], roots[i], roots[i+1]}
	}

	return len(bad), result
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func isValidSolution(n int, originalEdges [][2]int, operations [][4]int) bool {
	uf := NewUnionFind(n)

	for _, edge := range originalEdges {
		uf.union(edge[0], edge[1])
	}

	for _, op := range operations {
		x1, y1, x2, y2 := op[0], op[1], op[2], op[3]

		if uf.find(x1) != uf.find(y1) {
			return false
		}

		if uf.find(x2) == uf.find(y2) {
			return false
		}

		uf.union(x2, y2)
	}

	components := 0
	for i := 1; i <= n; i++ {
		if uf.find(i) == i {
			components++
		}
	}

	return components == 1
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		// Generate a graph with exactly n-1 edges, composed of several components,
		// and exactly one cycle per component except one (so total cycles = components-1),
		// which matches the problem's invariant and avoids degenerate cases for candidates.
		n := rng.Intn(18) + 2 // 2..19
		// Decide number of components c (at least 1, at most n)
		c := rng.Intn(imin(n, 5)) + 1 // keep small number of components
		// Split n nodes into c positive sizes
		sizes := make([]int, c)
		rem := n
		for k := 0; k < c; k++ {
			if k == c-1 {
				sizes[k] = rem
				break
			}
			// at least 1 node per remaining component
			maxHere := rem - (c - 1 - k)
			// choose s in [1, maxHere-1] if possible; if maxHere==1, s must be 1
			var s int
			if maxHere <= 1 {
				s = 1
			} else {
				s = rng.Intn(maxHere-1) + 1
			}
			sizes[k] = s
			rem -= s
		}
		// Assign labels 1..n to components contiguously
		edges := make([][2]int, 0, n-1)
		base := 1
		for k := 0; k < c; k++ {
			s := sizes[k]
			// build a tree on nodes [base..base+s-1]
			for t := 1; t < s; t++ {
				u := base + t
				v := base + rng.Intn(t) // connect to some previous node
				edges = append(edges, [2]int{u, v})
			}
			base += s
		}
		// add (c-1) extra edges, one within each of the first (c-1) components, to create a single cycle per component
		base = 1
		for k := 0; k < c-1; k++ {
			s := sizes[k]
			if s >= 2 {
				u := base + rng.Intn(s)
				v := base + rng.Intn(s)
				for v == u {
					v = base + rng.Intn(s)
				}
				edges = append(edges, [2]int{u, v})
			} else {
				// If component size is 1, connect it to itself is invalid; instead, add extra edge in last component
				// fallback: add inside last component (which must have size>=2 unless n==c)
				lastBase := base
				for kk := k + 1; kk < c; kk++ {
					lastBase += sizes[kk-1]
				}
				s2 := sizes[c-1]
				if s2 < 2 { // degenerate, adjust by linking into previous component
					// connect two nodes from previous component (which had at least 2 because we built a tree on it if s>1)
					prevBase := base - sizes[k]
					sPrev := sizes[k]
					if sPrev < 2 {
						continue
					}
					u := prevBase + rng.Intn(sPrev)
					v := prevBase + rng.Intn(sPrev)
					for v == u {
						v = prevBase + rng.Intn(sPrev)
					}
					edges = append(edges, [2]int{u, v})
				} else {
					u := lastBase + rng.Intn(s2)
					v := lastBase + rng.Intn(s2)
					for v == u {
						v = lastBase + rng.Intn(s2)
					}
					edges = append(edges, [2]int{u, v})
				}
			}
			base += sizes[k]
		}
		// Now edges should be exactly n-1
		if len(edges) != n-1 {
			// Trim or pad by connecting within the last component
			for len(edges) > n-1 {
				edges = edges[:len(edges)-1]
			}
			for len(edges) < n-1 {
				// add within last component
				lastStart := 1
				for kk := 0; kk < c-1; kk++ {
					lastStart += sizes[kk]
				}
				s := sizes[c-1]
				if s < 2 {
					break
				}
				u := lastStart + rng.Intn(s)
				v := lastStart + rng.Intn(s)
				for v == u {
					v = lastStart + rng.Intn(s)
				}
				edges = append(edges, [2]int{u, v})
			}
		}

		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for _, edge := range edges {
			input.WriteString(fmt.Sprintf("%d %d\n", edge[0], edge[1]))
		}

		expectedDays, _ := expected(n, edges)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}

		lines := strings.Split(got, "\n")
		if len(lines) < 1 {
			fmt.Fprintf(os.Stderr, "case %d failed: empty output\ninput:\n%s", i+1, input.String())
			os.Exit(1)
		}

		gotDays, parseErr := strconv.Atoi(lines[0])
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: cannot parse number of days: %v\ninput:\n%s", i+1, parseErr, input.String())
			os.Exit(1)
		}

		if gotDays != expectedDays {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d days, got %d\ninput:\n%s", i+1, expectedDays, gotDays, input.String())
			os.Exit(1)
		}

		if len(lines) != gotDays+1 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d operation lines, got %d\ninput:\n%s", i+1, gotDays, len(lines)-1, input.String())
			os.Exit(1)
		}

		operations := make([][4]int, gotDays)
		for j := 1; j <= gotDays; j++ {
			parts := strings.Fields(lines[j])
			if len(parts) != 4 {
				fmt.Fprintf(os.Stderr, "case %d failed: operation line %d should have 4 values\ninput:\n%s", i+1, j, input.String())
				os.Exit(1)
			}
			for k := 0; k < 4; k++ {
				val, parseErr := strconv.Atoi(parts[k])
				if parseErr != nil {
					fmt.Fprintf(os.Stderr, "case %d failed: cannot parse operation value: %v\ninput:\n%s", i+1, parseErr, input.String())
					os.Exit(1)
				}
				operations[j-1][k] = val
			}
		}

		if !isValidSolution(n, edges, operations) {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid solution operations\ninput:\n%s", i+1, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

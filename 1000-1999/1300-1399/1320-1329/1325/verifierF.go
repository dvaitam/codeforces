package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// verifyCase checks that the output is a valid answer for 1325F.
// The problem: given a graph with n vertices and m edges, either find:
// 1) An independent set of size ceil(sqrt(n)), or
// 2) A simple cycle of length at least ceil(sqrt(n))
func verifyCase(bin string, n int, edges [][2]int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	input := sb.String()

	got, err := run(bin, input)
	if err != nil {
		return err
	}

	lines := strings.Split(got, "\n")
	if len(lines) < 2 {
		return fmt.Errorf("too few output lines")
	}

	var typ int
	fmt.Sscan(lines[0], &typ)

	sqrtN := int(math.Ceil(math.Sqrt(float64(n))))

	if typ == 1 {
		// Independent set
		tokens := strings.Fields(lines[1])
		if len(tokens) < sqrtN {
			return fmt.Errorf("independent set too small: %d < %d", len(tokens), sqrtN)
		}
		nodes := make(map[int]bool)
		for _, tok := range tokens {
			var v int
			fmt.Sscan(tok, &v)
			if v < 1 || v > n {
				return fmt.Errorf("node %d out of range", v)
			}
			if nodes[v] {
				return fmt.Errorf("duplicate node %d", v)
			}
			nodes[v] = true
		}
		// Check no edge between nodes in the set
		for _, e := range edges {
			if nodes[e[0]] && nodes[e[1]] {
				return fmt.Errorf("edge (%d,%d) connects two nodes in independent set", e[0], e[1])
			}
		}
	} else if typ == 2 {
		// Cycle
		var cycleLen int
		fmt.Sscan(lines[1], &cycleLen)
		if cycleLen < sqrtN {
			return fmt.Errorf("cycle too short: %d < %d", cycleLen, sqrtN)
		}
		if len(lines) < 3 {
			return fmt.Errorf("missing cycle nodes")
		}
		tokens := strings.Fields(lines[2])
		if len(tokens) < cycleLen {
			return fmt.Errorf("cycle has %d nodes but claimed %d", len(tokens), cycleLen)
		}
		cycle := make([]int, cycleLen)
		seen := make(map[int]bool)
		for i := 0; i < cycleLen; i++ {
			fmt.Sscan(tokens[i], &cycle[i])
			if cycle[i] < 1 || cycle[i] > n {
				return fmt.Errorf("node %d out of range", cycle[i])
			}
			if seen[cycle[i]] {
				return fmt.Errorf("duplicate node %d in cycle", cycle[i])
			}
			seen[cycle[i]] = true
		}
		// Build adjacency set for quick edge lookup
		edgeSet := make(map[[2]int]bool)
		for _, e := range edges {
			edgeSet[[2]int{e[0], e[1]}] = true
			edgeSet[[2]int{e[1], e[0]}] = true
		}
		// Check all consecutive edges in cycle
		for i := 0; i < cycleLen; i++ {
			u := cycle[i]
			v := cycle[(i+1)%cycleLen]
			if !edgeSet[[2]int{u, v}] {
				return fmt.Errorf("edge (%d,%d) not in graph", u, v)
			}
		}
	} else {
		return fmt.Errorf("invalid type %d", typ)
	}

	return nil
}

func genGraph(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, 0)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	mExtra := rng.Intn(n)
	for i := 0; i < mExtra; i++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			i--
			continue
		}
		edges = append(edges, [2]int{u, v})
	}
	return edges
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(15) + 2
		edges := genGraph(rng, n)
		if err := verifyCase(bin, n, edges); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

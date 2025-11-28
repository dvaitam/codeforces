package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const modE int = 1000000007

type CaseE struct {
	input    string
	expected int
}

// Reference solver based on DSU
func solveE(n int, pairs [][2]int) int {
	// Map coordinates to 1..limit to handle potentially sparse (though generator makes them dense)
	// Actually generator guarantees 1..2n.
	// Max coordinate
	maxCoord := 0
	for _, p := range pairs {
		if p[0] > maxCoord {
			maxCoord = p[0]
		}
		if p[1] > maxCoord {
			maxCoord = p[1]
		}
	}

	parent := make([]int, maxCoord+1)
	sz := make([]int, maxCoord+1)
	hasLoop := make([]bool, maxCoord+1)
	selfLoop := make([]bool, maxCoord+1)
	
	for i := 1; i <= maxCoord; i++ {
		parent[i] = i
		sz[i] = 1
	}

	var find func(int) int
	find = func(i int) int {
		if parent[i] == i {
			return i
		}
		parent[i] = find(parent[i])
		return parent[i]
	}

	union := func(i, j int) {
		rootI := find(i)
		rootJ := find(j)
		if rootI != rootJ {
			parent[rootI] = rootJ
			sz[rootJ] += sz[rootI]
			if hasLoop[rootI] { hasLoop[rootJ] = true }
			if hasLoop[rootJ] { hasLoop[rootJ] = true } // redundant but clear
			if selfLoop[rootI] { selfLoop[rootJ] = true }
			if selfLoop[rootJ] { selfLoop[rootJ] = true }
		} else {
			hasLoop[rootI] = true
		}
	}

	// Process edges
	for _, p := range pairs {
		u, v := p[0], p[1]
		if u == v {
			root := find(u)
			hasLoop[root] = true
			selfLoop[root] = true
		} else {
			union(u, v)
		}
	}

	// Calculate answer
	ans := 1
	
	// Use map to collect unique roots for iteration
	roots := make(map[int]bool)
	for _, p := range pairs {
		roots[find(p[0])] = true
	}
	
	for r := range roots {
		if hasLoop[r] {
			if !selfLoop[r] {
				ans = (ans * 2) % modE
			}
		} else {
			ans = (ans * sz[r]) % modE
		}
	}
	
	return ans
}

func generateCaseE(rng *rand.Rand) CaseE {
	n := rng.Intn(20) + 1
	
	// Generate distinct starting seats u from 1..2n
	seats := rng.Perm(2 * n)
	for i := range seats {
		seats[i]++ // 1-based
	}
	
	uList := seats[:n]
	
	pairs := make([][2]int, n)
	for i := 0; i < n; i++ {
		u := uList[i]
		v := rng.Intn(2*n) + 1
		pairs[i] = [2]int{u, v}
	}
	
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, p := range pairs {
		sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
	}
	return CaseE{sb.String(), solveE(n, pairs)}
}

func runCase(exe, input string, expected int) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(strings.TrimSpace(out.String())), &got); err != nil {
		return fmt.Errorf("cannot parse output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	cases := make([]CaseE, 0, 100)
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCaseE(rng))
	}
	
	for i, tc := range cases {
		if err := runCase(exe, tc.input, tc.expected); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

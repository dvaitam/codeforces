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

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func solveF(n int, a []int64, edges [][2]int) string {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	down := make([]int64, n+1)
	diam := make([]int64, n+1)
	var finalAns int64

	// DFS 1 to compute subtree properties (rooted at 1)
	var dfs1 func(u, p int)
	dfs1 = func(u, p int) {
		var d1, d2 int64
		var maxSubDiam int64
		for _, v := range adj[u] {
			if v == p {
				continue
			}
			dfs1(v, u)
			if down[v] > d1 {
				d2 = d1
				d1 = down[v]
			} else if down[v] > d2 {
				d2 = down[v]
			}
			maxSubDiam = max(maxSubDiam, diam[v])
		}
		down[u] = a[u] + d1
		diam[u] = max(maxSubDiam, a[u]+d1+d2)
	}
	dfs1(1, -1)

	// In the real problem, we need to find TWO disjoint paths that maximize total weight.
	// The problem asks for "Alice and Bob... collect as many chocolates... total number...".
	// This is equivalent to finding two vertex-disjoint paths in the tree with maximum total weight sum.
	// The provided "correct solution" implements an O(N) re-rooting DP to essentially iterate over all edges,
	// removing the edge to split the tree into two components, and taking max(diam(T1) + diam(T2)).
	// This covers the case where the two paths are in separate subtrees created by cutting an edge.
	// HOWEVER, there is another case: One path is in a subtree, and the other path is disjoint but in the same component?
	// Actually, any two disjoint paths can be separated by removing some edge.
	// So iterating over all edges and taking sum of diameters of the two resulting trees is sufficient and correct.
	
	// Re-implementing the O(N) DP logic correctly here as Oracle.
	
	// We use the same re-rooting technique but clean.
	
	type State struct {
		d  int64 // max path starting at root going down
		dm int64 // diameter in subtree
	}
	
	// Collect child states for re-rooting
	// We need to pass "upper" tree info to children.
	// Upper tree info consists of:
	// 1. Longest path starting at parent going upwards (upPath)
	// 2. Diameter of the upper tree (upDiam)
	
	var dfs2 func(u, p int, upPath, upDiam int64)
	dfs2 = func(u, p int, upPath, upDiam int64) {
		// Update answer considering cut (u, p)
		// Tree 1: subtree u (diam[u])
		// Tree 2: rest of tree (upDiam)
		// Note: For root, p=-1, upDiam=0 (invalid cut but handled by max).
		if p != -1 {
			finalAns = max(finalAns, diam[u]+upDiam)
		}

		// Prepare for children
		children := []int{}
		for _, v := range adj[u] {
			if v != p {
				children = append(children, v)
			}
		}
		k := len(children)
		
		prefD := make([]int64, k+1)
		suffD := make([]int64, k+1)
		prefDm := make([]int64, k+1)
		suffDm := make([]int64, k+1)
		
		// Fill down paths
		for i := 0; i < k; i++ {
			v := children[i]
			prefD[i+1] = max(prefD[i], down[v])
			prefDm[i+1] = max(prefDm[i], diam[v])
		}
		for i := k - 1; i >= 0; i-- {
			v := children[i]
			suffD[i] = max(suffD[i+1], down[v])
			suffDm[i] = max(suffDm[i+1], diam[v])
		}

		// Also we need top 2 down paths to form path through u
		// We can get top 2 from prefixes/suffixes or just finding top 2.
		// Finding top 2 is cleaner for pathThroughU calculation.
		
		dList := make([]int64, k)
		for i, v := range children {
			dList[i] = down[v]
		}
		
		// Function to get max path through u using children excluding child i
		// This is max(pathAbove + a[u] + maxSibDown, a[u] + sumTop2SibDown)
		
		// To optimize, we can rely on prefD/suffD for maxSibDown.
		// For sumTop2SibDown, it's slightly more complex with prefixes.
		// Let's compute top 3 down paths globally for u's children to handle exclusion.
		
		type ValIdx struct { val int64; idx int }
		top3 := []ValIdx{{-1, -1}, {-1, -1}, {-1, -1}}
		
		for i, v := range children {
			d := down[v]
			if d > top3[0].val {
				top3[2], top3[1], top3[0] = top3[1], top3[0], ValIdx{d, i}
			} else if d > top3[1].val {
				top3[2], top3[1] = top3[1], ValIdx{d, i}
			} else if d > top3[2].val {
				top3[2] = ValIdx{d, i}
			}
		}

		for i, v := range children {
			// Calculate newUpPath for child v
			// upPath is max(pathAbove, maxSibDown) + a[u]
			maxSibD := max(prefD[i], suffD[i+1])
			nextUpPath := max(upPath, maxSibD) + a[u]
			
			// Calculate newUpDiam for child v
			// Candidates for diameter in the tree excluding v:
			// 1. upDiam (diameter fully in upper part)
			// 2. maxSibDiam (diameter fully in one of the siblings)
			// 3. Path starting in upper part, going through u, ending in a sibling
			//    = upPath + a[u] + maxSibDown
			// 4. Path starting in one sibling, going through u, ending in another sibling
			//    = maxPairSibDown + a[u]
			
			maxSibDm := max(prefDm[i], suffDm[i+1])
			
			// Top 2 siblings excluding current i
			sum2 := int64(0)
			cnt := 0
			for _, item := range top3 {
				if item.idx != i && item.idx != -1 {
					sum2 += item.val
					cnt++
					if cnt == 2 { break }
				}
			}
			pathThroughU_SibSib := int64(0)
			if cnt == 2 {
				pathThroughU_SibSib = sum2 + a[u]
			} else if cnt == 1 {
				// Only 1 sibling exists (excluding v), so path is just that sibling + u
				pathThroughU_SibSib = sum2 + a[u] 
			} else {
				// No siblings
				pathThroughU_SibSib = a[u]
			}
			// Wait, path connecting two siblings requires 2 siblings.
			// If only 1 sibling, it's just a path ending at u?
			// The definition of diameter is max distance between any two nodes.
			// If only 1 sibling, diameter could be that path.
			// However, cand3 covers "up to sibling".
			// If no siblings, cand3 is upPath + a[u].
			
			pathThroughU_UpSib := upPath + a[u] + maxSibD
			
			nextUpDiam := max(upDiam, max(maxSibDm, max(pathThroughU_UpSib, pathThroughU_SibSib)))
			
			dfs2(v, u, nextUpPath, nextUpDiam)
		}
	}
	
	dfs2(1, -1, 0, 0)
	
	return fmt.Sprintf("%d\n", finalAns)
}

func generateCaseF(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 2 // Reasonable size for fuzzing
	weights := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		weights[i] = int64(rng.Intn(20) + 1)
	}
	edges := make([][2]int, 0, n-1)
	// Random tree generation
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(weights[i]))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	
	// Use the Oracle (re-implemented solveF) to get expected output
	expect := solveF(n, weights, edges)
	return sb.String(), expect
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

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCaseF(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
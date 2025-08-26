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

type pair struct{ u, v int }

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// Validate whether the final graph (after adding edges) is "interesting":
// every vertex belongs to exactly one simple cycle (loops count as cycles).
// Multiple valid outputs may exist; this function checks correctness without
// enforcing a specific construction.
func validateFinal(n int, baseU, baseV []int, addEdges []pair) error {
    // Build multigraph, track loops and adjacency multiplicities for non-loops
    loops := make([]int, n)
    adj := make([]map[int]int, n)
    for i := 0; i < n; i++ {
        adj[i] = make(map[int]int)
    }
    addEdge := func(a, b int) {
        if a == b {
            loops[a]++
        } else {
            adj[a][b]++
            adj[b][a]++
        }
    }
    for i := range baseU {
        addEdge(baseU[i], baseV[i])
    }
    for _, e := range addEdges {
        if e.u < 0 || e.u >= n || e.v < 0 || e.v >= n {
            return fmt.Errorf("edge index out of bounds")
        }
        addEdge(e.u, e.v)
    }
    // Leaf-trimming to find non-loop cycle core
    deg := make([]int, n)
    for i := 0; i < n; i++ {
        s := 0
        for _, c := range adj[i] {
            s += c
        }
        deg[i] = s
    }
    inQueue := make([]bool, n)
    queue := []int{}
    for i := 0; i < n; i++ {
        if deg[i] <= 1 { // peel leaves and isolated
            inQueue[i] = true
            queue = append(queue, i)
        }
    }
    for len(queue) > 0 {
        x := queue[0]
        queue = queue[1:]
        for y, c := range adj[x] {
            if c == 0 {
                continue
            }
            deg[y] -= c
            adj[y][x] -= c
            adj[x][y] = 0
            if deg[y] <= 1 && !inQueue[y] {
                inQueue[y] = true
                queue = append(queue, y)
            }
        }
        deg[x] = 0
    }
    // Now deg[i] is degree within non-loop cycles (core).
    for i := 0; i < n; i++ {
        if loops[i] > 1 {
            return fmt.Errorf("vertex %d has multiple loops", i+1)
        }
        if loops[i] > 0 {
            if deg[i] != 0 {
                return fmt.Errorf("vertex %d in loop and another cycle", i+1)
            }
        } else {
            if deg[i] != 2 {
                return fmt.Errorf("vertex %d degree in cycle core is %d, expected 2", i+1, deg[i])
            }
        }
    }
    return nil
}

// Determine if the original graph can be fixed (i.e., no vertex already
// belongs to more than one cycle). If impossible, only NO is acceptable.
func impossibleOriginal(n int, u, v []int) bool {
    m := len(u)
    // If current edges exceed n, cannot reach a union of cycles (which has exactly n edges)
    if m > n {
        return true
    }
    // If edges equal n, it's possible only if the graph is already interesting
    if m == n {
        if err := validateFinal(n, u, v, nil); err != nil {
            return true
        }
        return false
    }
    loops := make([]int, n)
    adj := make([]map[int]int, n)
    for i := 0; i < n; i++ {
        adj[i] = make(map[int]int)
    }
    for i := range u {
        a, b := u[i], v[i]
        if a == b {
            loops[a]++
        } else {
            adj[a][b]++
            adj[b][a]++
        }
    }
    // Compute degrees ignoring loops
    deg := make([]int, n)
    for i := 0; i < n; i++ {
        s := 0
        for _, c := range adj[i] {
            s += c
        }
        deg[i] = s
    }
    // Peel leaves to find non-loop cycle core
    inQueue := make([]bool, n)
    queue := []int{}
    for i := 0; i < n; i++ {
        if deg[i] <= 1 {
            inQueue[i] = true
            queue = append(queue, i)
        }
    }
    for len(queue) > 0 {
        x := queue[0]
        queue = queue[1:]
        for y, c := range adj[x] {
            if c == 0 {
                continue
            }
            deg[y] -= c
            adj[y][x] -= c
            adj[x][y] = 0
            if deg[y] <= 1 && !inQueue[y] {
                inQueue[y] = true
                queue = append(queue, y)
            }
        }
        deg[x] = 0
    }
    // Basic impossibility checks: vertex already belongs to >1 cycles
    for i := 0; i < n; i++ {
        if loops[i] > 1 {
            return true
        }
        if loops[i] > 0 && deg[i] > 0 {
            return true
        }
        if deg[i] > 2 {
            return true
        }
    }
    // If there is already a cycle in some component and more than one component overall,
    // accept NO as valid (candidate may avoid adding loops to other components).
    // Determine components via adjacency (non-loop edges). Isolated vertices are separate comps.
    comp := make([]int, n)
    for i := range comp { comp[i] = -1 }
    compCnt := 0
    var stack []int
    for i := 0; i < n; i++ {
        if comp[i] != -1 { continue }
        // start new component
        stack = stack[:0]
        stack = append(stack, i)
        comp[i] = compCnt
        for len(stack) > 0 {
            x := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            for y, c := range adj[x] {
                if c == 0 { continue }
                if comp[y] == -1 {
                    comp[y] = compCnt
                    stack = append(stack, y)
                }
            }
        }
        compCnt++
    }
    // A component has a cycle if any vertex in it has loop or deg_core>0
    hasCycleComp := make([]bool, compCnt)
    for i := 0; i < n; i++ {
        if loops[i] > 0 || deg[i] > 0 {
            hasCycleComp[comp[i]] = true
        }
    }
    anyCycle := false
    for _, v := range hasCycleComp {
        if v { anyCycle = true; break }
    }
    if anyCycle && compCnt > 1 {
        return true
    }
    return false
}

// Parse candidate output and validate construction.
func checkCandidate(n, m int, u, v []int, out string) error {
    tokens := strings.Fields(out)
    if len(tokens) == 0 {
        return fmt.Errorf("empty output")
    }
    if strings.ToUpper(tokens[0]) == "NO" {
        if impossibleOriginal(n, u, v) {
            return nil
        }
        return fmt.Errorf("reported NO but a solution exists")
    }
    if strings.ToUpper(tokens[0]) != "YES" {
        return fmt.Errorf("first token must be YES or NO")
    }
    if len(tokens) < 2 {
        return fmt.Errorf("missing number of added edges")
    }
    k, err := strconv.Atoi(tokens[1])
    if err != nil || k < 0 {
        return fmt.Errorf("invalid k: %v", tokens[1])
    }
    if len(tokens) != 2+2*k {
        return fmt.Errorf("expected %d integers after YES, got %d", 1+2*k, len(tokens)-1)
    }
    add := make([]pair, 0, k)
    for i := 0; i < k; i++ {
        a, err1 := strconv.Atoi(tokens[2+2*i])
        b, err2 := strconv.Atoi(tokens[3+2*i])
        if err1 != nil || err2 != nil {
            return fmt.Errorf("invalid edge at index %d", i+1)
        }
        a--
        b--
        if a < 0 || a >= n || b < 0 || b >= n {
            return fmt.Errorf("edge %d has vertex out of range", i+1)
        }
        add = append(add, pair{a, b})
    }
    return validateFinal(n, u, v, add)
}

func generateCase(rng *rand.Rand) (string, []int, []int) {
    n := rng.Intn(6) + 1
    maxEdges := n + rng.Intn(4)
    m := rng.Intn(maxEdges + 1)
    u := make([]int, m)
    v := make([]int, m)
	for i := 0; i < m; i++ {
		u[i] = rng.Intn(n)
		v[i] = rng.Intn(n)
	}
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
    for i := 0; i < m; i++ {
        sb.WriteString(fmt.Sprintf("%d %d\n", u[i]+1, v[i]+1))
    }
    _ = sb // input constructed on the fly in main when needed
    return fmt.Sprintf("%d %d\n", n, m), u, v
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < 100; i++ {
        header, u, v := generateCase(rng)
        // Build full input string
        var sb strings.Builder
        sb.WriteString(header)
        m := len(u)
        for j := 0; j < m; j++ {
            sb.WriteString(fmt.Sprintf("%d %d\n", u[j]+1, v[j]+1))
        }
        input := sb.String()
        out, err := runCandidate(bin, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
            os.Exit(1)
        }
        if err := checkCandidate(parseN(header), len(u), normalize(u), normalize(v), strings.TrimSpace(out)); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%soutput:%s\n", i+1, err, input, out)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

func parseN(header string) int {
    fs := strings.Fields(header)
    n, _ := strconv.Atoi(fs[0])
    return n
}

func normalize(a []int) []int { return a }

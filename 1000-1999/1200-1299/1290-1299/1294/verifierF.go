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

type Edge struct{ u, v int }

type Test struct{
    n int
    edges []Edge
}

func (t Test) Input() string {
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d\n", t.n))
    for _, e := range t.edges {
        sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
    }
    return sb.String()
}

func runProg(bin, input string) (string, error) {
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

func genTest(rng *rand.Rand) Test {
    n := rng.Intn(8)+3
    edges := make([]Edge, n-1)
    for i:=2;i<=n;i++ {
        p := rng.Intn(i-1)+1
        edges[i-2] = Edge{p,i}
    }
    return Test{n, edges}
}

func buildAdj(n int, edges []Edge) [][]int {
    adj := make([][]int, n+1)
    for _, e := range edges {
        adj[e.u] = append(adj[e.u], e.v)
        adj[e.v] = append(adj[e.v], e.u)
    }
    return adj
}

func pathEdges(adj [][]int, s, t int) map[[2]int]struct{} {
    n := len(adj)-1
    parent := make([]int, n+1)
    for i := range parent { parent[i] = -1 }
    queue := []int{s}
    parent[s] = 0
    for q := 0; q < len(queue); q++ {
        u := queue[q]
        if u == t { break }
        for _, v := range adj[u] {
            if parent[v] != -1 { continue }
            parent[v] = u
            queue = append(queue, v)
        }
    }
    edges := make(map[[2]int]struct{})
    for v := t; v != s; v = parent[v] {
        u := parent[v]
        if u == -1 { break }
        if u < v {
            edges[[2]int{u,v}] = struct{}{}
        } else {
            edges[[2]int{v,u}] = struct{}{}
        }
    }
    return edges
}

func unionLen(adj [][]int, a,b,c int) int {
    m := make(map[[2]int]struct{})
    for k := range pathEdges(adj,a,b) { m[k]=struct{}{} }
    for k := range pathEdges(adj,b,c) { m[k]=struct{}{} }
    for k := range pathEdges(adj,a,c) { m[k]=struct{}{} }
    return len(m)
}

func bestLen(adj [][]int) int {
    n := len(adj)-1
    best := 0
    for a:=1;a<=n;a++ {
        for b:=a+1;b<=n;b++ {
            for c:=b+1;c<=n;c++ {
                l := unionLen(adj,a,b,c)
                if l > best { best = l }
            }
        }
    }
    return best
}

func runCase(bin string, tc Test) error {
    adj := buildAdj(tc.n, tc.edges)
    expect := bestLen(adj)
    out, err := runProg(bin, tc.Input())
    if err != nil { return err }
    tokens := strings.Fields(out)
    if len(tokens) != 4 {
        return fmt.Errorf("expected 4 numbers in output, got %d", len(tokens))
    }
    res, err := strconv.Atoi(tokens[0])
    if err != nil { return fmt.Errorf("invalid number: %v", err) }
    a, _ := strconv.Atoi(tokens[1])
    b, _ := strconv.Atoi(tokens[2])
    c, _ := strconv.Atoi(tokens[3])
    if a<1||a>tc.n||b<1||b>tc.n||c<1||c>tc.n { return fmt.Errorf("nodes out of range") }
    if a==b || a==c || b==c { return fmt.Errorf("nodes not distinct") }
    gotLen := unionLen(adj,a,b,c)
    if gotLen != res {
        return fmt.Errorf("reported length %d but actual %d", res, gotLen)
    }
    if res != expect {
        return fmt.Errorf("answer not optimal: expected %d got %d", expect, res)
    }
    return nil
}

func main(){
    if len(os.Args)!=2 {
        fmt.Println("Usage: go run verifierF.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    const cases = 100
    for i:=0;i<cases;i++ {
        tc := genTest(rng)
        if err := runCase(bin, tc); err != nil {
            fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, tc.Input())
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", cases)
}


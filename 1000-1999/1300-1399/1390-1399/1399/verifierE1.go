package main

import (
    "bufio"
    "container/heap"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

type edgeE1 struct {
    to, id int
    w int64
}

type testCaseE1 struct {
    n   int
    S   int64
    edges [][3]int
}

func parseTestcasesE1(path string) ([]testCaseE1, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    scanner := bufio.NewScanner(f)
    var cases []testCaseE1
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }
        parts := strings.Fields(line)
        if len(parts) != 2 {
            return nil, fmt.Errorf("bad header: %s", line)
        }
        n, _ := strconv.Atoi(parts[0])
        S, _ := strconv.ParseInt(parts[1], 10, 64)
        edges := make([][3]int, 0, n-1)
        for i := 0; i < n-1; i++ {
            if !scanner.Scan() {
                return nil, fmt.Errorf("unexpected EOF")
            }
            eparts := strings.Fields(scanner.Text())
            if len(eparts) != 3 {
                return nil, fmt.Errorf("bad edge line")
            }
            u, _ := strconv.Atoi(eparts[0])
            v, _ := strconv.Atoi(eparts[1])
            w, _ := strconv.Atoi(eparts[2])
            edges = append(edges, [3]int{u, v, w})
        }
        cases = append(cases, testCaseE1{n, S, edges})
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    return cases, nil
}

// solver from 1399E1.go
func solveE1(n int, S int64, edges [][3]int) int64 {
    adj := make([][]edgeE1, n+1)
    m := n - 1
    weights := make([]int64, m)
    for i := 0; i < m; i++ {
        u := edges[i][0]
        v := edges[i][1]
        w := int64(edges[i][2])
        adj[u] = append(adj[u], edgeE1{to: v, id: i, w: w})
        adj[v] = append(adj[v], edgeE1{to: u, id: i, w: w})
        weights[i] = w
    }
    leafCount := make([]int64, m)
    type Frame struct{ v, parent, eid, idx int; leaves int64 }
    stack := make([]Frame, 0, n)
    stack = append(stack, Frame{v: 1, parent: 0, eid: -1, idx: 0, leaves: 0})
    for len(stack) > 0 {
        f := &stack[len(stack)-1]
        if f.idx < len(adj[f.v]) {
            e := adj[f.v][f.idx]
            f.idx++
            if e.to == f.parent {
                continue
            }
            stack = append(stack, Frame{v: e.to, parent: f.v, eid: e.id, idx: 0, leaves: 0})
        } else {
            if f.leaves == 0 {
                f.leaves = 1
            }
            if f.eid >= 0 {
                leafCount[f.eid] = f.leaves
            }
            count := f.leaves
            stack = stack[:len(stack)-1]
            if len(stack) > 0 {
                stack[len(stack)-1].leaves += count
            }
        }
    }
    var total int64
    pq := &MaxHeap{}
    heap.Init(pq)
    for i := 0; i < m; i++ {
        total += weights[i] * leafCount[i]
        saving := (weights[i] - weights[i]/2) * leafCount[i]
        heap.Push(pq, &Item{saving: saving, id: i})
    }
    var moves int64
    for total > S {
        it := heap.Pop(pq).(*Item)
        if it.saving <= 0 {
            break
        }
        total -= it.saving
        moves++
        i := it.id
        weights[i] /= 2
        newSave := (weights[i] - weights[i]/2) * leafCount[i]
        heap.Push(pq, &Item{saving: newSave, id: i})
    }
    return moves
}

// Max-heap utilities

type Item struct {
    saving int64
    id     int
}

type MaxHeap []*Item

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i].saving > h[j].saving }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(*Item)) }
func (h *MaxHeap) Pop() interface{} {
    old := *h
    n := len(old)
    it := old[n-1]
    *h = old[:n-1]
    return it
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
    var errBuf bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errBuf
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierE1.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    cases, err := parseTestcasesE1("testcasesE1.txt")
    if err != nil {
        fmt.Println("failed to parse testcases:", err)
        os.Exit(1)
    }
    for idx, tc := range cases {
        var sb strings.Builder
        sb.WriteString("1\n")
        sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.S))
        for _, e := range tc.edges {
            sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
        }
        expected := strconv.FormatInt(solveE1(tc.n, tc.S, tc.edges), 10)
        got, err := run(bin, sb.String())
        if err != nil {
            fmt.Printf("case %d failed: %v\n", idx+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != expected {
            fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expected, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(cases))
}


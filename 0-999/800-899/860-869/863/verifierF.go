package main

import (
    "bytes"
    "container/heap"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
)

type Edge struct {
    to, rev, cap int
    cost       int
}

type Graph [][]*Edge

func (g Graph) AddEdge(u, v, cap, cost int) {
    g[u] = append(g[u], &Edge{to: v, rev: len(g[v]), cap: cap, cost: cost})
    g[v] = append(g[v], &Edge{to: u, rev: len(g[u]) - 1, cap: 0, cost: -cost})
}

type Item struct {
    v    int
    dist int
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
    old := *pq
    n := len(old)
    it := old[n-1]
    *pq = old[:n-1]
    return it
}

func minCostFlow(g Graph, s, t, maxf int) (int, int) {
    n := len(g)
    h := make([]int, n)
    prevv := make([]int, n)
    preve := make([]int, n)
    flow, cost := 0, 0
    const INF = int(1 << 60)
    for flow < maxf {
        dist := make([]int, n)
        for i := range dist { dist[i] = INF }
        dist[s] = 0
        pq := &PriorityQueue{}
        heap.Init(pq)
        heap.Push(pq, Item{v: s, dist: 0})
        for pq.Len() > 0 {
            it := heap.Pop(pq).(Item)
            v := it.v
            if dist[v] < it.dist { continue }
            for i,e := range g[v] {
                if e.cap > 0 && dist[e.to] > dist[v]+e.cost+h[v]-h[e.to] {
                    dist[e.to] = dist[v] + e.cost + h[v] - h[e.to]
                    prevv[e.to] = v
                    preve[e.to] = i
                    heap.Push(pq, Item{v:e.to, dist:dist[e.to]})
                }
            }
        }
        if dist[t] == INF { break }
        for v:=0; v<n; v++ { if dist[v] < INF { h[v] += dist[v] } }
        d := maxf - flow
        for v:=t; v!=s; v = prevv[v] {
            if d > g[prevv[v]][preve[v]].cap {
                d = g[prevv[v]][preve[v]].cap
            }
        }
        flow += d
        cost += d * h[t]
        for v:=t; v!=s; v = prevv[v] {
            e := g[prevv[v]][preve[v]]
            e.cap -= d
            g[v][e.rev].cap += d
        }
    }
    return flow, cost
}

func expectedF(n,q int, queries [][4]int) string {
    lb := make([]int,n)
    ub := make([]int,n)
    for i:=0;i<n;i++ { lb[i]=1; ub[i]=n }
    for _,qr := range queries {
        t,l,r,v := qr[0],qr[1]-1,qr[2]-1,qr[3]
        if t==1 {
            for j:=l;j<=r;j++ { if lb[j] < v { lb[j]=v } }
        } else {
            for j:=l;j<=r;j++ { if ub[j] > v { ub[j]=v } }
        }
    }
    for i:=0;i<n;i++ { if lb[i]>ub[i] { return "-1" } }
    V := 2*n+2
    s := 0
    posBase := 1
    valBase := posBase+n
    tIdx := valBase+n
    g := make(Graph,V)
    for i:=0;i<n;i++ {
        g.AddEdge(s,posBase+i,1,0)
        for v:=lb[i]; v<=ub[i]; v++ { g.AddEdge(posBase+i,valBase+(v-1),1,0) }
    }
    for v:=0; v<n; v++ {
        for k:=1; k<=n; k++ {
            cost := 2*k - 1
            g.AddEdge(valBase+v,tIdx,1,cost)
        }
    }
    flow,cost := minCostFlow(g,s,tIdx,n)
    if flow < n { return "-1" }
    return fmt.Sprintf("%d", cost)
}

func genTestsF() []string {
    rand.Seed(6)
    tests := make([]string,0,100)
    for len(tests)<100 {
        n := rand.Intn(3)+1
        q := rand.Intn(3)
        var queries [][4]int
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d %d\n", n,q))
        for i:=0;i<q;i++ {
            t := rand.Intn(2)+1
            l := rand.Intn(n)+1
            r := rand.Intn(n-l+1)+l
            v := rand.Intn(n)+1
            sb.WriteString(fmt.Sprintf("%d %d %d %d\n", t,l,r,v))
            queries = append(queries,[4]int{t,l,r,v})
        }
        tests = append(tests, sb.String())
    }
    return tests
}

func runBinary(bin string, input string) (string,error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func parseQueries(lines []string) [][4]int {
    res := make([][4]int,len(lines))
    for i,l := range lines {
        fmt.Sscanf(l, "%d %d %d %d", &res[i][0],&res[i][1],&res[i][2],&res[i][3])
    }
    return res
}

func main(){
    if len(os.Args)!=2 { fmt.Fprintf(os.Stderr,"Usage: go run verifierF.go <binary>\n"); os.Exit(1) }
    bin := os.Args[1]
    tests := genTestsF()
    for idx,t := range tests {
        lines := strings.Split(strings.TrimSpace(t),"\n")
        var n,q int
        fmt.Sscanf(lines[0],"%d %d", &n,&q)
        queries := parseQueries(lines[1:1+q])
        want := expectedF(n,q,queries)
        got, err := runBinary(bin,t)
        if err != nil {
            fmt.Printf("Test %d: runtime error: %v\n", idx+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got)!=want {
            fmt.Printf("Test %d failed.\nInput:\n%s\nExpected: %s\nGot: %s\n", idx+1, t, want, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed.\n", len(tests))
}


package main

import (
   "bufio"
   "fmt"
   "os"
)

const eps int64 = 1e8

// Edge holds input edge parameters
type Edge struct {
   x, y    int
   a, b, c, d int64
}

var n, m int
var edges []Edge

// Dinic implements max flow
type Dinic struct {
   n     int
   graph [][]edgeD
   level []int
   iter  []int
}

// edgeD is flow edge
type edgeD struct {
   to, rev int
   cap      int64
}

// NewDinic creates a Dinic object with n nodes
func NewDinic(n int) *Dinic {
   g := make([][]edgeD, n)
   level := make([]int, n)
   iter := make([]int, n)
   return &Dinic{n: n, graph: g, level: level, iter: iter}
}

// AddEdge adds directed edge from->to with capacity cap
func (d *Dinic) AddEdge(from, to int, cap int64) {
   d.graph[from] = append(d.graph[from], edgeD{to, len(d.graph[to]), cap})
   d.graph[to] = append(d.graph[to], edgeD{from, len(d.graph[from]) - 1, 0})
}

// bfs constructs level graph
func (d *Dinic) bfs(s, t int) {
   for i := range d.level {
       d.level[i] = -1
   }
   queue := make([]int, 0, d.n)
   d.level[s] = 0
   queue = append(queue, s)
   for qi := 0; qi < len(queue); qi++ {
       v := queue[qi]
       for _, e := range d.graph[v] {
           if e.cap > 0 && d.level[e.to] < 0 {
               d.level[e.to] = d.level[v] + 1
               queue = append(queue, e.to)
           }
       }
   }
}

// dfs finds augmenting path
func (d *Dinic) dfs(v, t int, f int64) int64 {
   if v == t {
       return f
   }
   for i := d.iter[v]; i < len(d.graph[v]); i++ {
       e := &d.graph[v][i]
       if e.cap > 0 && d.level[v] < d.level[e.to] {
           w := d.dfs(e.to, t, min(f, e.cap))
           if w > 0 {
               e.cap -= w
               d.graph[e.to][e.rev].cap += w
               return w
           }
       }
       d.iter[v]++
   }
   return 0
}

// MaxFlow computes max flow from s to t
func (d *Dinic) MaxFlow(s, t int) int64 {
   flow := int64(0)
   inf := int64(1<<62 - 1)
   for {
       d.bfs(s, t)
       if d.level[t] < 0 {
           break
       }
       for i := range d.iter {
           d.iter[i] = 0
       }
       for {
           f := d.dfs(s, t, inf)
           if f == 0 {
               break
           }
           flow += f
       }
   }
   return flow
}

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

// calc builds flow graph for time t and returns sum-flow
func calc(t int64) int64 {
   deg := make([]int64, n+1)
   d := NewDinic(n + 2)
   S, T := 0, n+1
   for _, e := range edges {
       left := e.a*t + e.b
       right := e.c*t + e.d
       d.AddEdge(e.x, e.y, right-left)
       deg[e.x] -= left
       deg[e.y] += left
   }
   sum := int64(0)
   for i := 1; i <= n; i++ {
       if deg[i] >= 0 {
           d.AddEdge(S, i, deg[i])
           sum += deg[i]
       } else {
           d.AddEdge(i, T, -deg[i])
       }
   }
   flow := d.MaxFlow(S, T)
   return sum - flow
}

// process expands interval around St and prints result
func process(St int64) {
   L, R := St, St
   for i := 30; i >= 0; i-- {
       step := int64(1) << uint(i)
       if L >= step && calc(L-step) == 0 {
           L -= step
       }
   }
   for i := 30; i >= 0; i-- {
       step := int64(1) << uint(i)
       if R+step <= eps && calc(R+step) == 0 {
           R += step
       }
   }
   res := float64(R-L) / float64(eps)
   fmt.Printf("%.6f\n", res)
   os.Exit(0)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   fmt.Fscan(reader, &n, &m)
   edges = make([]Edge, m)
   for i := 0; i < m; i++ {
       var x, y int
       var a, b, c, d_ int64
       fmt.Fscan(reader, &x, &y, &a, &b, &c, &d_)
       b *= eps
       d_ *= eps
       edges[i] = Edge{x, y, a, b, c, d_}
   }
   Left, Right := int64(0), eps
   for Left <= Right {
       lenThird := (Right - Left) / 3
       mid1 := Left + lenThird
       mid2 := Right - lenThird
       val1 := calc(mid1)
       val2 := calc(mid2)
       if val1 == 0 {
           process(mid1)
       }
       if val2 == 0 {
           process(mid2)
       }
       if val1 <= val2 {
           Right = mid2 - 1
       } else {
           Left = mid1 + 1
       }
   }
   fmt.Printf("0.000000\n")
}

package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

// Dinic max flow
type Edge struct {
   to, rev int
   cap     int64
}

type Dinic struct {
   n     int
   graph [][]Edge
   level []int
   iter  []int
}

func NewDinic(n int) *Dinic {
   d := &Dinic{
       n:     n,
       graph: make([][]Edge, n),
       level: make([]int, n),
       iter:  make([]int, n),
   }
   return d
}

func (d *Dinic) AddEdge(from, to int, cap int64) {
   d.graph[from] = append(d.graph[from], Edge{to, len(d.graph[to]), cap})
   d.graph[to] = append(d.graph[to], Edge{from, len(d.graph[from]) - 1, 0})
}

func (d *Dinic) bfs(s int) {
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

func (d *Dinic) dfs(v, t int, f int64) int64 {
   if v == t {
       return f
   }
   for i := d.iter[v]; i < len(d.graph[v]); i++ {
       e := &d.graph[v][i]
       if e.cap > 0 && d.level[v] < d.level[e.to] {
           ret := d.dfs(e.to, t, minInt64(f, e.cap))
           if ret > 0 {
               e.cap -= ret
               d.graph[e.to][e.rev].cap += ret
               return ret
           }
       }
       d.iter[v]++
   }
   return 0
}

func (d *Dinic) MaxFlow(s, t int) int64 {
   flow := int64(0)
   const inf = math.MaxInt64
   for {
       d.bfs(s)
       if d.level[t] < 0 {
           return flow
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
}

func minInt64(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   var g int64
   fmt.Fscan(reader, &n, &m, &g)
   s0 := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &s0[i])
   }
   v := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &v[i])
   }
   type Bet struct {
       t    int
       w    int64
       k    int
       idx  []int
       fr   int
   }
   bets := make([]Bet, m)
   var totalF int64 = 0
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &bets[i].t, &bets[i].w, &bets[i].k)
       bets[i].idx = make([]int, bets[i].k)
       for j := 0; j < bets[i].k; j++ {
           fmt.Fscan(reader, &bets[i].idx[j])
           bets[i].idx[j]--
       }
       fmt.Fscan(reader, &bets[i].fr)
       if bets[i].fr == 1 {
           totalF += g
       }
   }
   N := m + n + 2
   src := m + n
   sink := m + n + 1
   din := NewDinic(N)
   const INF = int64(1e18)
   var sumP int64 = 0
   for i := 0; i < m; i++ {
       Pi := bets[i].w
       if bets[i].fr == 1 {
           Pi += g
       }
       sumP += Pi
       din.AddEdge(src, i, Pi)
       // constraints
       for _, j := range bets[i].idx {
           if s0[j] == bets[i].t {
               // require x_j = 0: edge j-node to bet
               din.AddEdge(m+j, i, INF)
           } else {
               // require x_j = 1: edge bet to j-node
               din.AddEdge(i, m+j, INF)
           }
       }
   }
   for j := 0; j < n; j++ {
       if v[j] > 0 {
           din.AddEdge(m+j, sink, v[j])
       }
   }
   flow := din.MaxFlow(src, sink)
   res := sumP - flow - totalF
   fmt.Fprintln(writer, res)
}

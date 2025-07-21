package main

import (
   "bufio"
   "container/list"
   "fmt"
   "os"
   "sort"
)

// Dinic max flow implementation
type edge struct{ to, rev int; cap int }
type Dinic struct{ N int; G [][]edge; level, ptr []int }
func NewDinic(n int) *Dinic {
   g := make([][]edge, n)
   return &Dinic{N: n, G: g, level: make([]int, n), ptr: make([]int, n)}
}
func (d *Dinic) AddEdge(u, v, c int) {
   d.G[u] = append(d.G[u], edge{v, len(d.G[v]), c})
   d.G[v] = append(d.G[v], edge{u, len(d.G[u]) - 1, 0})
}
func (d *Dinic) bfs(s, t int) bool {
   for i := range d.level { d.level[i] = -1 }
   q := list.New()
   d.level[s] = 0
   q.PushBack(s)
   for q.Len() > 0 {
       u := q.Remove(q.Front()).(int)
       for _, e := range d.G[u] {
           if d.level[e.to] < 0 && e.cap > 0 {
               d.level[e.to] = d.level[u] + 1
               q.PushBack(e.to)
           }
       }
   }
   return d.level[t] >= 0
}
func (d *Dinic) dfs(u, t, f int) int {
   if u == t || f == 0 { return f }
   for &d.ptr[u]; d.ptr[u] < len(d.G[u]); d.ptr[u]++ {
       e := &d.G[u][d.ptr[u]]
       if d.level[e.to] == d.level[u]+1 && e.cap > 0 {
           pushed := d.dfs(e.to, t, min(f, e.cap))
           if pushed > 0 {
               e.cap -= pushed
               d.G[e.to][e.rev].cap += pushed
               return pushed
           }
       }
   }
   return 0
}
func (d *Dinic) MaxFlow(s, t int) int {
   flow := 0
   for d.bfs(s, t) {
       for i := range d.ptr { d.ptr[i] = 0 }
       for {
           pushed := d.dfs(s, t, 1e9)
           if pushed == 0 { break }
           flow += pushed
       }
   }
   return flow
}
func min(a, b int) int { if a < b { return a }; return b }

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, males, females int
   if _, err := fmt.Fscan(in, &n, &m, &males, &females); err != nil {
       return
   }
   grid := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &grid[i])
   }
   total := 1 + males + females
   type Sc struct{ r, c, t int }
   sc := make([]Sc, total)
   // id 0: boss
   for i := 0; i < total; i++ {
       fmt.Fscan(in, &sc[i].r, &sc[i].c, &sc[i].t)
       sc[i].r--; sc[i].c--
   }
   // map free cells
   cellIdx := make([][]int, n)
   cells := make([][2]int, 0)
   for i := 0; i < n; i++ {
       cellIdx[i] = make([]int, m)
       for j := 0; j < m; j++ {
           if grid[i][j] == '.' {
               cellIdx[i][j] = len(cells)
               cells = append(cells, [2]int{i, j})
           } else {
               cellIdx[i][j] = -1
           }
       }
   }
   nc := len(cells)
   // BFS distSteps per scayger
   INF := n*m + 5
   distSteps := make([][]int, total)
   for i := 0; i < total; i++ {
       dist := make([]int, nc)
       for j := range dist { dist[j] = INF }
       // BFS
       q := list.New()
       sr, sc0 := sc[i].r, sc[i].c
       start := cellIdx[sr][sc0]
       dist[start] = 0
       q.PushBack(start)
       for q.Len() > 0 {
           u := q.Remove(q.Front()).(int)
           ur, uc := cells[u][0], cells[u][1]
           for _, dxy := range [][2]int{{-1,0},{1,0},{0,-1},{0,1}} {
               vr, vc := ur+dxy[0], uc+dxy[1]
               if vr >= 0 && vr < n && vc >= 0 && vc < m && cellIdx[vr][vc] >= 0 {
                   v := cellIdx[vr][vc]
                   if dist[v] > dist[u]+1 {
                       dist[v] = dist[u] + 1
                       q.PushBack(v)
                   }
               }
           }
       }
       distSteps[i] = dist
   }
   // collect all possible times
   times := make([]int64, 0, total*nc)
   for i := 0; i < total; i++ {
       ti := int64(sc[i].t)
       for j := 0; j < nc; j++ {
           if distSteps[i][j] < INF {
               times = append(times, ti*int64(distSteps[i][j]))
           }
       }
   }
   if len(times) == 0 {
       fmt.Println(-1)
       return
   }
   sort.Slice(times, func(i, j int) bool { return times[i] < times[j] })
   uniq := times[:1]
   for i := 1; i < len(times); i++ {
       if times[i] != uniq[len(uniq)-1] {
           uniq = append(uniq, times[i])
       }
   }
   // side assignment based on majority
   var bossOnMale bool
   if males == females+1 {
       bossOnMale = true
   } else if females == males+1 {
       bossOnMale = false
   } else {
       fmt.Println(-1)
       return
   }
   // binary search on uniq
   lo, hi := 0, len(uniq)-1
   ans := int64(-1)
   for lo <= hi {
       mid := (lo + hi) / 2
       T := uniq[mid]
       // build flow
       var sideA, sideB []int
       // male side
       for i := 1; i <= males; i++ { sideA = append(sideA, i) }
       if bossOnMale { sideA = append(sideA, 0) }
       // female side
       for i := males + 1; i < total; i++ { sideB = append(sideB, i) }
       if !bossOnMale { sideB = append(sideB, 0) }
       szA, szB := len(sideA), len(sideB)
       if szA != szB {
           lo = mid + 1; continue
       }
       // build graph
       nc0 := nc
       S := 0
       baseA := 1
       baseCin := baseA + szA
       baseCout := baseCin + nc0
       baseB := baseCout + nc0
       Tnode := baseB + szB
       din := NewDinic(Tnode + 1)
       // edges S->A
       for i := 0; i < szA; i++ {
           din.AddEdge(S, baseA+i, 1)
       }
       // cell split cin->cout
       for c := 0; c < nc0; c++ {
           din.AddEdge(baseCin+c, baseCout+c, 1)
       }
       // edges A->cin
       for i, sid := range sideA {
           ti := int64(sc[sid].t)
           for c := 0; c < nc0; c++ {
               if int64(distSteps[sid][c])*ti <= T {
                   din.AddEdge(baseA+i, baseCin+c, 1)
               }
           }
       }
       // edges cout->B
       for j, sid := range sideB {
           ti := int64(sc[sid].t)
           for c := 0; c < nc0; c++ {
               if int64(distSteps[sid][c])*ti <= T {
                   din.AddEdge(baseCout+c, baseB+j, 1)
               }
           }
       }
       // edges B->T
       for j := 0; j < szB; j++ {
           din.AddEdge(baseB+j, Tnode, 1)
       }
       flow := din.MaxFlow(S, Tnode)
       if flow == szA {
           ans = T; hi = mid - 1
       } else {
           lo = mid + 1
       }
   }
   if ans < 0 {
       fmt.Println(-1)
   } else {
       fmt.Println(ans)
   }
}

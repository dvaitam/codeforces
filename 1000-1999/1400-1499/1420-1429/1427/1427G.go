package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Dinic max flow
type edge struct{ to, rev int; cap int }
type Dinic struct{
   n int
   g [][]edge
   level []int
   it []int
}
func NewDinic(n int) *Dinic {
   d := &Dinic{n:n, g:make([][]edge,n), level:make([]int,n), it:make([]int,n)}
   return d
}
func (d *Dinic) AddEdge(s, t, cap int) {
   d.g[s] = append(d.g[s], edge{t, len(d.g[t]), cap})
   d.g[t] = append(d.g[t], edge{s, len(d.g[s])-1, 0})
}
func (d *Dinic) bfs(s, t int) bool {
   for i := range d.level { d.level[i] = -1 }
   q := make([]int,0, d.n)
   d.level[s] = 0; q = append(q, s)
   for i := 0; i < len(q); i++ {
       u := q[i]
       for _, e := range d.g[u] {
           if e.cap>0 && d.level[e.to]<0 {
               d.level[e.to] = d.level[u]+1
               q = append(q, e.to)
           }
       }
   }
   return d.level[t]>=0
}
func (d *Dinic) dfs(u, t, f int) int {
   if u==t { return f }
   for i:=d.it[u]; i<len(d.g[u]); i++ {
       e := &d.g[u][i]
       if e.cap>0 && d.level[e.to]==d.level[u]+1 {
           ret := d.dfs(e.to, t, min(f, e.cap))
           if ret>0 {
               e.cap -= ret
               d.g[e.to][e.rev].cap += ret
               d.it[u] = i
               return ret
           }
       }
   }
   d.it[u] = len(d.g[u])
   return 0
}
func (d *Dinic) MaxFlow(s, t int) int {
   flow := 0
   for d.bfs(s,t) {
       for i:=range d.it { d.it[i]=0 }
       for {
           f := d.dfs(s,t,1<<60)
           if f==0 { break }
           flow += f
       }
   }
   return flow
}
func min(a, b int) int { if a<b { return a }; return b }

var (
   n int
   a [][]int
   idGrid [][]int
   posX, posY []int
   fixed0, fixed1 []bool
   vals []int
   delta []int
   ans int64
)

// solve thresholds in [L..R] on vars slice of global ids
func solve(L, R int, vars []int) {
   if L>R || len(vars)==0 { return }
   mid := (L+R)>>1
   // build mapping
   // map global id to local index in vars
   localID := make([]int, len(posX))
   for i := range localID { localID[i] = -1 }
   for i, id := range vars {
       localID[id] = i
   }
   V := len(vars) + 2
   s, t := 0, 1
   d := NewDinic(V)
   T := vals[mid] // threshold value
   // build graph
   for _, id := range vars {
       u := 2 + localID[id]
       // get original coords
       x, y := posX[id], posY[id]
       srcCap, sinkCap := 0, 0
       // four dirs
       for _, dxy := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
           nx, ny := x+dxy[0], y+dxy[1]
           if nx<0 || nx>=n || ny<0 || ny>=n { continue }
           v := a[nx][ny]
           if v == -1 {
               continue
           }
           // fixed painted (boundary or already painted interior)
           vid0 := idGrid[nx][ny]
           if v > 0 && vid0 < 0 {
               if v <= T {
                   srcCap++
               } else {
                   sinkCap++
               }
           } else if vid0 >= 0 {
               // interior neighbor
               if fixed0[vid0] {
                   srcCap++
               } else if fixed1[vid0] {
                   sinkCap++
               } else {
                   // both variable, add adjacency once
                   if id < vid0 {
                       vLocal := localID[vid0]
                       if vLocal >= 0 {
                           d.AddEdge(u, 2+vLocal, 1)
                           d.AddEdge(2+vLocal, u, 1)
                       }
                   }
               }
           }
       }
       if srcCap>0 { d.AddEdge(s, 2+localID[id], srcCap) }
       if sinkCap>0 { d.AddEdge(2+localID[id], t, sinkCap) }
   }
   flow := d.MaxFlow(s, t)
   ans += int64(flow) * int64(delta[mid])
   // find reachable
   vis := make([]bool, V)
   // BFS in residual
   q := []int{s}; vis[s]=true
   for len(q)>0 {
       u := q[0]; q = q[1:]
       for _, e := range d.g[u] {
           if e.cap>0 && !vis[e.to] {
               vis[e.to] = true
               q = append(q, e.to)
           }
       }
   }
   var leftVars, rightVars []int
   for _, id := range vars {
       if vis[2+localID[id]] {
           // reachable => in S1 => label1
           rightVars = append(rightVars, id)
       } else {
           leftVars = append(leftVars, id)
       }
   }
   // recurse left [L..mid-1], S0 are leftVars but fixed1 for rightVars
   for _, id := range rightVars { fixed1[id] = true }
   solve(L, mid-1, leftVars)
   for _, id := range rightVars { fixed1[id] = false }
   // recurse right [mid+1..R]
   for _, id := range leftVars { fixed0[id] = true }
   solve(mid+1, R, rightVars)
   for _, id := range leftVars { fixed0[id] = false }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   fmt.Fscan(in, &n)
   a = make([][]int, n)
   for i := 0; i < n; i++ {
       row := make([]int, n)
       for j := 0; j < n; j++ {
           fmt.Fscan(in, &row[j])
       }
       a[i] = row
   }
   // collect boundary values
   mset := make(map[int]struct{})
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if i==0||i==n-1||j==0||j==n-1 {
               if a[i][j] > 0 {
                   mset[a[i][j]] = struct{}{}
               }
           }
       }
   }
   for v := range mset { vals = append(vals, v) }
   sort.Ints(vals)
   if len(vals) < 2 {
       fmt.Println(0)
       return
   }
   m := len(vals)
   delta = make([]int, m)
   for i := 1; i < m; i++ { delta[i] = vals[i] - vals[i-1] }
   // assign interior ids
   idGrid = make([][]int, n)
   for i := 0; i < n; i++ {
       idGrid[i] = make([]int, n)
       for j := 0; j < n; j++ {
           idGrid[i][j] = -1
       }
   }
   nextID := 0
   for i := 1; i < n-1; i++ {
       for j := 1; j < n-1; j++ {
           if a[i][j] != -1 {
               idGrid[i][j] = nextID
               posX = append(posX, i)
               posY = append(posY, j)
               nextID++
           }
       }
   }
   fixed0 = make([]bool, nextID)
   fixed1 = make([]bool, nextID)
   // initial vars all interior ids
   vars := make([]int, nextID)
   for i := 0; i < nextID; i++ { vars[i] = i }
   solve(1, m-1, vars)
   // add contrast for initially painted adjacent non-broken tiles (fixed-fixed)
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if a[i][j] <= 0 { continue }
           for _, dxy := range [][2]int{{1, 0}, {0, 1}} {
               ni, nj := i+dxy[0], j+dxy[1]
               if ni<0 || ni>=n || nj<0 || nj>=n { continue }
               if a[ni][nj] > 0 {
                   diff := a[i][j] - a[ni][nj]
                   if diff < 0 { diff = -diff }
                   ans += int64(diff)
               }
           }
       }
   }
   fmt.Println(ans)
}

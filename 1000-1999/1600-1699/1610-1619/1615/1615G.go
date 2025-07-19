package main

import (
   "bufio"
   "fmt"
   "os"
)

// Maximum matching (Edmonds' blossom) for general graph
type MaximumMatching struct {
   n             int
   mate, parity  []int
   blossom, parent, last []int
   adj           [][]int
   bfsQueue      []int
   size, timer   int
}

func NewMaximumMatching(n int) *MaximumMatching {
   mm := &MaximumMatching{n: n}
   mm.mate = make([]int, n)
   mm.parity = make([]int, n)
   mm.blossom = make([]int, n)
   mm.parent = make([]int, n)
   mm.last = make([]int, n)
   mm.adj = make([][]int, n)
   for i := 0; i < n; i++ {
       mm.mate[i] = -1
   }
   return mm
}

func (mm *MaximumMatching) AddEdge(u, v int) {
   mm.adj[u] = append(mm.adj[u], v)
   mm.adj[v] = append(mm.adj[v], u)
}

func (mm *MaximumMatching) augment(u int) {
   for u != -1 {
       v := mm.parent[u]
       w := mm.mate[v]
       mm.mate[u] = v
       mm.mate[v] = u
       u = w
   }
}

func (mm *MaximumMatching) lca(u, v int) int {
   mm.timer++
   for {
       if u != -1 {
           if mm.last[u] == mm.timer {
               return u
           }
           mm.last[u] = mm.timer
           if mm.mate[u] == -1 {
               u = -1
           } else {
               u = mm.blossom[mm.parent[mm.mate[u]]]
           }
       }
       u, v = v, u
   }
}

func (mm *MaximumMatching) merge(u, v, p int) {
   for mm.blossom[u] != p {
       mm.parent[u] = v
       v = mm.mate[u]
       if mm.parity[v] == 1 {
           mm.parity[v] = 0
           mm.bfsQueue = append(mm.bfsQueue, v)
       }
       mm.blossom[u] = p
       mm.blossom[v] = p
       u = mm.parent[v]
   }
}

func (mm *MaximumMatching) bfs(root int) bool {
   // init
   for i := 0; i < mm.n; i++ {
       mm.parity[i] = -1
       mm.parent[i] = -1
       mm.blossom[i] = i
   }
   mm.bfsQueue = mm.bfsQueue[:0]
   mm.bfsQueue = append(mm.bfsQueue, root)
   mm.parity[root] = 0
   // BFS
   for qi := 0; qi < len(mm.bfsQueue); qi++ {
       u := mm.bfsQueue[qi]
       for _, v := range mm.adj[u] {
           if mm.parity[v] == -1 {
               mm.parity[v] = 1
               mm.parent[v] = u
               if mm.mate[v] == -1 {
                   mm.augment(v)
                   return true
               }
               mm.parity[mm.mate[v]] = 0
               mm.bfsQueue = append(mm.bfsQueue, mm.mate[v])
           } else if mm.parity[v] == 0 && mm.blossom[u] != mm.blossom[v] {
               p := mm.lca(mm.blossom[u], mm.blossom[v])
               mm.merge(u, v, p)
               mm.merge(v, u, p)
           }
       }
   }
   return false
}

func (mm *MaximumMatching) Solve() int {
   for i := 0; i < mm.n; i++ {
       if mm.mate[i] == -1 {
           if mm.bfs(i) {
               mm.size++
           }
       }
   }
   return mm.size
}

// OrEdge represents an OR constraint edge
type OrEdge struct{ v, id, l, r int }
type PairLR struct{ l, r int }

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   a := make([]int, n)
   k := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
       if a[i] > k {
           k = a[i]
       }
   }
   // prepare structures
   orEdges := make([][]OrEdge, k+1)
   andEdges := make([][]PairLR, k+1)
   for i := range andEdges {
       andEdges[i] = make([]PairLR, k+1)
       for j := range andEdges[i] {
           andEdges[i][j] = PairLR{-1, -1}
       }
   }
   visitedEdge := make([]bool, n)
   paired := make([]bool, k+1)
   exposed := make([]bool, k+1)
   visited := make([]bool, k+1)
   color := make([]int, k+1)

   // build edges
   m := 0
   for i := -1; i < n; {
       j := i + 1
       for j < n && a[j] == 0 {
           j++
       }
       zeroes := j - i - 1
       var endpoints [2]struct{v, pos int}
       cnt := 0
       if i >= 0 {
           endpoints[cnt] = struct{v, pos int}{a[i], i + 1}
           cnt++
       }
       if j < n {
           endpoints[cnt] = struct{v, pos int}{a[j], j - 1}
           cnt++
       }
       if cnt > 0 {
           al, l := endpoints[0].v, endpoints[0].pos
           ar, r := endpoints[cnt-1].v, endpoints[cnt-1].pos
           if zeroes%2==1 || (zeroes==0 && cnt==2 && al==ar) {
               orEdges[al] = append(orEdges[al], OrEdge{ar, m, l, r})
               orEdges[ar] = append(orEdges[ar], OrEdge{al, m, r, l})
               m++
           } else if zeroes>0 {
               andEdges[al][ar] = PairLR{l, r}
               andEdges[ar][al] = PairLR{r, l}
           }
       }
       i = j
   }

   // color components
   var dfsColor func(u, c int)
   dfsColor = func(u, c int) {
       color[u] = c
       for _, e := range orEdges[u] {
           if color[e.v] == 0 {
               dfsColor(e.v, c)
           }
       }
   }
   for i := 1; i <= k; i++ {
       if color[i] == 0 {
           dfsColor(i, i)
           vertices, edges := 0, 0
           for j := 1; j <= k; j++ {
               if color[j] == i {
                   edges += len(orEdges[j])
                   vertices++
               }
           }
           edges /= 2
           exposed[i] = (edges == vertices-1)
       }
   }

   // matching on components
   mm := NewMaximumMatching(k + 1)
   edgeSet := make(map[[2]int]bool)
   for i := 1; i <= k; i++ {
       for j := 1; j <= k; j++ {
           pr := andEdges[i][j]
           if pr.l >= 0 {
               u, v := color[i], color[j]
               if exposed[u] && exposed[v] {
                   key := [2]int{u, v}
                   if !edgeSet[key] {
                       mm.AddEdge(u, v)
                       edgeSet[key] = true
                   }
               }
           }
       }
   }
   mm.Solve()

   // dfs solve uses orEdges
   var dfsSolve func(u int) bool
   dfsSolve = func(u int) bool {
       covered := false
       visited[u] = true
       for _, e := range orEdges[u] {
           if visitedEdge[e.id] {
               continue
           }
           visitedEdge[e.id] = true
           if !visited[e.v] {
               if dfsSolve(e.v) {
                   if !covered {
                       covered = true
                       a[e.l] = u
                   }
               } else {
                   a[e.r] = e.v
               }
           } else if !covered {
               covered = true
               a[e.l] = u
           }
       }
       return covered
   }

   // assign AND matched pairs
   for i := 1; i <= k; i++ {
       for j := 1; j <= k; j++ {
           pr := andEdges[i][j]
           if pr.l >= 0 {
               u, v := color[i], color[j]
               if exposed[u] && exposed[v] && mm.mate[u] == v {
                   a[pr.l] = i
                   a[pr.r] = j
                   dfsSolve(i)
                   dfsSolve(j)
                   exposed[u] = false
                   exposed[v] = false
               }
           }
       }
   }
   // remaining exposed roots
   for i := 1; i <= k; i++ {
       if color[i]==i && mm.mate[i]==-1 {
           dfsSolve(i)
       }
   }

   // fill leftover zeros
   rem := make(map[int]struct{})
   for i := 1; i <= n; i++ {
       rem[i] = struct{}{}
   }
   for _, v := range a {
       delete(rem, v)
   }
   for i := 0; i < n; i++ {
       if a[i] == 0 {
           var v int
           for x := range rem {
               v = x
               break
           }
           delete(rem, v)
           a[i] = v
           if i+1<n && a[i+1]==0 {
               a[i+1] = v
               i++
           }
       }
   }
   // output
   for i := 0; i < n; i++ {
       fmt.Fprint(out, a[i], " ")
   }
   fmt.Fprintln(out)
}

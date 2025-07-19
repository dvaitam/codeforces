package main

import (
   "bufio"
   "fmt"
   "os"
)

// Edge represents a directed edge with capacity and reverse index.
type Edge struct {
   to, rev, cap int
}

var (
   graph [][]Edge
   level []int
   ptr   []int
   q     []int
)

// addEdge adds a directed edge u->v with capacity c.
func addEdge(u, v, c int) {
   graph[u] = append(graph[u], Edge{v, len(graph[v]), c})
   graph[v] = append(graph[v], Edge{u, len(graph[u]) - 1, 0})
}

// bfs constructs level graph and returns whether t is reachable.
func bfs(s, t int) bool {
   for i := range level {
       level[i] = -1
   }
   head, tail := 0, 0
   q[tail] = s
   tail++
   level[s] = 0
   for head < tail {
       v := q[head]
       head++
       for _, e := range graph[v] {
           if level[e.to] < 0 && e.cap > 0 {
               level[e.to] = level[v] + 1
               q[tail] = e.to
               tail++
           }
       }
   }
   return level[t] >= 0
}

// dfs finds an augmenting path and returns flow.
func dfs(v, t, pushed int) int {
   if pushed == 0 {
       return 0
   }
   if v == t {
       return pushed
   }
   for ptr[v] < len(graph[v]) {
       e := &graph[v][ptr[v]]
       if level[e.to] == level[v]+1 && e.cap > 0 {
           tr := dfs(e.to, t, min(pushed, e.cap))
           if tr > 0 {
               e.cap -= tr
               graph[e.to][e.rev].cap += tr
               return tr
           }
       }
       ptr[v]++
   }
   return 0
}

// dinic runs the Dinic algorithm and returns max flow from s to t.
func dinic(s, t int) int {
   flow := 0
   for bfs(s, t) {
       for i := range ptr {
           ptr[i] = 0
       }
       for pushed := dfs(s, t, 1<<60); pushed > 0; pushed = dfs(s, t, 1<<60) {
           flow += pushed
       }
   }
   return flow
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var m, k int
   if _, err := fmt.Fscan(in, &m, &k); err != nil {
       return
   }
   a := make([]int, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(in, &a[i])
   }
   rest := m*(m-1)
   for _, v := range a {
       rest -= v
   }
   if rest < 0 || rest > a[k-1]*(m-k) {
       fmt.Fprintln(out, "no")
       return
   }
   s := k
   t := k + 1
   fake := k + 2
   n := k + 3
   graph = make([][]Edge, n)
   level = make([]int, n)
   ptr = make([]int, n)
   q = make([]int, n)

   // Build network
   for i := 0; i < k; i++ {
       addEdge(s, i, 2*(m-1-i))
   }
   addEdge(s, fake, (m-k)*(m-k-1))
   for i := 0; i < k; i++ {
       for j := i + 1; j < k; j++ {
           addEdge(i, j, 2)
       }
   }
   for i := 0; i < k; i++ {
       addEdge(i, fake, 2*(m-k))
   }
   for i := 0; i < k; i++ {
       addEdge(i, t, a[i])
   }
   addEdge(fake, t, rest)

   flow := dinic(s, t)
   if flow != m*(m-1) {
       fmt.Fprintln(out, "no")
       return
   }
   fmt.Fprintln(out, "yes")

   // Prepare answer matrix
   ans := make([][]int, m)
   for i := 0; i < m; i++ {
       ans[i] = make([]int, m)
   }

   // Record flows between main nodes
   for i := 0; i < k; i++ {
       for _, e := range graph[i] {
           j := e.to
           if j < k && i < j {
               // initial cap = 2, residual = e.cap
               ans[i][j] = e.cap        // wins of i over j
               ans[j][i] = 2 - e.cap    // wins of j over i
           }
       }
   }

   // Record flows to fake node
   fFake := make([]int, k)
   for i := 0; i < k; i++ {
       for _, e := range graph[i] {
           if e.to == fake {
               fFake[i] = 2*(m-k) - e.cap
               break
           }
       }
   }

   // Default between main and fake group
   for i := 0; i < k; i++ {
       for j := k; j < m; j++ {
           ans[i][j] = 2
           ans[j][i] = 0
       }
   }

   // Distribute losses based on flow to fake
   tmp := k
   for i := 0; i < k; i++ {
       for cnt := 0; cnt < fFake[i]; cnt++ {
           ans[i][tmp]--
           ans[tmp][i]++
           tmp++
           if tmp == m {
               tmp = k
           }
       }
   }

   // Between fake-group nodes: draws
   for i := k; i < m; i++ {
       for j := k; j < m; j++ {
           ans[i][j] = 1
       }
   }

   // Output matrix
   for i := 0; i < m; i++ {
       for j := 0; j < m; j++ {
           if i == j {
               out.WriteByte('X')
           } else if ans[i][j] == 2 {
               out.WriteByte('W')
           } else if ans[i][j] == 1 {
               out.WriteByte('D')
           } else {
               out.WriteByte('L')
           }
       }
       out.WriteByte('\n')
   }
}

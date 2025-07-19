package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n    int
   g    [][]int
   crossEdge [][]int
   a, b []int
   done []int
   cnt  int
   id   []int
   side []bool
)

func bfs(u int, dist []int) int {
   cnt++
   dist[u] = 0
   done[u] = cnt
   q := make([]int, 0, len(g))
   q = append(q, u)
   var curr int
   for i := 0; i < len(q); i++ {
       curr = q[i]
       for _, v := range g[curr] {
           if done[v] != cnt {
               done[v] = cnt
               dist[v] = dist[curr] + 1
               q = append(q, v)
           }
       }
   }
   return curr
}

func makePermutation(l, r int) {
   if l+1 == r {
       return
   }
   self := id[l]
   otherNear := g[self][0]
   otherFar := bfs(self, a)
   bfs(otherNear, a)
   dist := a[otherFar]
   bfs(otherFar, b)
   for i := l; i <= r; i++ {
       if a[id[i]]+b[id[i]] == dist {
           side[id[i]] = true
       } else {
           side[id[i]] = false
       }
   }
   our, their := l, r
   for our < their {
       if !side[id[our]] {
           our++
       } else if side[id[their]] {
           their--
       } else {
           id[our], id[their] = id[their], id[our]
           our++
           their--
       }
   }
   our = r
   their = l
   for side[id[our]] {
       our--
   }
   for !side[id[their]] {
       their++
   }
   // now our+1 == their
   for i := l; i <= our; i++ {
       u := id[i]
       // find and remove one cross edge
       for j, v := range g[u] {
           if side[v] {
               crossEdge[u] = append(crossEdge[u], v)
               // remove v from g[u]
               last := len(g[u]) - 1
               g[u][j], g[u][last] = g[u][last], g[u][j]
               g[u] = g[u][:last]
               break
           }
       }
   }
   makePermutation(l, our)
   for i := l; i <= our; i++ {
       u := id[i]
       sz := len(crossEdge[u])
       id[i+our-l+1] = crossEdge[u][sz-1]
       crossEdge[u] = crossEdge[u][:sz-1]
   }
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
   fmt.Fscan(reader, &n)
   N := 1 << n
   // init
   g = make([][]int, N)
   crossEdge = make([][]int, N)
   a = make([]int, N)
   b = make([]int, N)
   done = make([]int, N)
   id = make([]int, N)
   side = make([]bool, N)
   // read edges
   m := n << (n - 1)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       g[u] = append(g[u], v)
       g[v] = append(g[v], u)
   }
   for i := 0; i < N; i++ {
       id[i] = i
   }
   cnt = 0
   makePermutation(0, N-1)
   // output permutation
   for i := 0; i < N; i++ {
       writer.WriteString(fmt.Sprint(id[i]))
       if i+1 < N {
           writer.WriteByte(' ')
       } else {
           writer.WriteByte('\n')
       }
   }
   // check power of two
   if n == 0 || (n&(n-1)) != 0 {
       writer.WriteString("-1\n")
       return
   }
   // compute and output colors
   color := make([]int, N)
   for i := 0; i < N; i++ {
       col := 0
       for j := 0; j < n; j++ {
           if (i>>j)&1 == 1 {
               col ^= j
           }
       }
       color[id[i]] = col
   }
   for i := 0; i < N; i++ {
       writer.WriteString(fmt.Sprint(color[i]))
       if i+1 < N {
           writer.WriteByte(' ')
       } else {
           writer.WriteByte('\n')
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   fmt.Fscan(reader, &t)
   for ; t > 0; t-- {
       solve(reader, writer)
   }
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

// Operation represents a transfer from u to v of amount w
type Operation struct {
   u, v, w int
}

var (
   n, v, m int
   a, b    []int
   parent  []int
   graph   [][]int
   ops     []Operation
   s       []int
)

func find(x int) int {
   if parent[x] != x {
       parent[x] = find(parent[x])
   }
   return parent[x]
}

func union(x, y int) {
   rx, ry := find(x), find(y)
   if rx != ry {
       parent[rx] = ry
   }
}

// pre computes s[x] = a[x] - b[x] + sum of s[children]
func pre(x, p int) {
   s[x] = a[x] - b[x]
   for _, y := range graph[x] {
       if y == p {
           continue
       }
       pre(y, x)
       s[x] += s[y]
   }
}

// dfs transfers up to z to node x from its subtree
func dfs(x, p, z int) {
   if z <= 0 {
       return
   }
   for _, y := range graph[x] {
       if y == p || z <= 0 {
           continue
       }
       if s[y] <= 0 {
           continue
       }
       // w = min(s[y], z)
       w := s[y]
       if w > z {
           w = z
       }
       z -= w
       t := a[y]
       if t >= w {
           // transfer w from y to x
           ops = append(ops, Operation{y, x, w})
           a[y] -= w
           a[x] += w
           if a[y] < b[y] {
               need := b[y] - a[y]
               if need > w {
                   need = w
               }
               dfs(y, x, need)
           }
       } else {
           // transfer all t first
           ops = append(ops, Operation{y, x, t})
           a[y] -= t
           a[x] += t
           // then dfs with adjusted amount
           extra := w
           if t >= b[y] {
               extra = w + (b[y] - t)
           }
           dfs(y, x, extra)
           // finally transfer the rest
           rest := w - t
           ops = append(ops, Operation{y, x, rest})
           a[y] -= rest
           a[x] += rest
       }
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   _, _ = fmt.Fscan(in, &n, &v, &m)
   a = make([]int, n+1)
   b = make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &b[i])
   }
   parent = make([]int, n+1)
   for i := 1; i <= n; i++ {
       parent[i] = i
   }
   graph = make([][]int, n+1)
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       if find(x) != find(y) {
           union(x, y)
           graph[x] = append(graph[x], y)
           graph[y] = append(graph[y], x)
       }
   }
   // check component sums
   compSum := make([]int, n+1)
   for i := 1; i <= n; i++ {
       compSum[find(i)] += a[i] - b[i]
   }
   for i := 1; i <= n; i++ {
       if compSum[i] != 0 {
           fmt.Println("NO")
           return
       }
   }
   // prepare for transfers
   s = make([]int, n+1)
   ops = make([]Operation, 0)
   for i := 1; i <= n; i++ {
       if a[i] < b[i] {
           pre(i, 0)
           dfs(i, 0, b[i]-a[i])
       }
   }
   // output
   fmt.Println(len(ops))
   for _, op := range ops {
       fmt.Printf("%d %d %d\n", op.u, op.v, op.w)
   }
}

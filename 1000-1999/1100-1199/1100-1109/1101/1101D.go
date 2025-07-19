package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   a := make([]int, n)
   ma := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
       if a[i] > ma {
           ma = a[i]
       }
   }
   // precompute divisors for values up to ma
   g0 := make([][]int, ma+1)
   for i := 2; i <= ma; i++ {
       for j := i; j <= ma; j += i {
           g0[j] = append(g0[j], i)
       }
   }
   // divisors per node
   divs := make([][]int, n)
   divSet := make([]map[int]bool, n)
   for i := 0; i < n; i++ {
       divs[i] = g0[a[i]]
       m := make(map[int]bool, len(divs[i]))
       for _, p := range divs[i] {
           m[p] = true
       }
       divSet[i] = m
   }
   // graph
   adj := make([][]int, n)
   for i := 0; i < n-1; i++ {
       u, v := 0, 0
       fmt.Fscan(in, &u, &v)
       u--
       v--
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // DP arrays: best and second best edges count
   f1 := make([]map[int]int, n)
   f2 := make([]map[int]int, n)
   parent := make([]int, n)
   for i := range parent {
       parent[i] = -1
   }
   // iterative DFS post-order
   type nodeState struct{ u int; visited bool }
   stack := []nodeState{{u: 0, visited: false}}
   parent[0] = -2
   ans := 0
   for len(stack) > 0 {
       cur := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       u := cur.u
       if !cur.visited {
           // pre
           stack = append(stack, nodeState{u: u, visited: true})
           for _, v := range adj[u] {
               if parent[v] == -1 {
                   parent[v] = u
                   stack = append(stack, nodeState{u: v, visited: false})
               }
           }
       } else {
           // post-order: process u
           f1[u] = make(map[int]int)
           f2[u] = make(map[int]int)
           // combine children
           for _, v := range adj[u] {
               if parent[v] != u {
                   continue
               }
               // for each prime in child v
               for _, p := range divs[v] {
                   if !divSet[u][p] {
                       continue
                   }
                   // child best edges for p
                   childBest := f1[v][p]
                   l := childBest + 1
                   // update best and second best for u
                   if f1[u][p] < l {
                       f2[u][p] = f1[u][p]
                       f1[u][p] = l
                   } else if f2[u][p] < l {
                       f2[u][p] = l
                   }
               }
           }
           // update answer, count nodes = edges1 + edges2 + 1
           for _, p := range divs[u] {
               v1 := f1[u][p]
               v2 := f2[u][p]
               if v1+v2+1 > ans {
                   ans = v1 + v2 + 1
               }
           }
       }
   }
   fmt.Println(ans)
}

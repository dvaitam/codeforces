package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   if n <= 2 {
       fmt.Fprintln(out, 0)
       return
   }
   adj := make([][]int, n)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       u--; v--
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // BFS to find farthest from 0
   far := func(src int) (int, []int, []int) {
       dist := make([]int, n)
       par := make([]int, n)
       for i := range dist {
           dist[i] = -1
       }
       q := make([]int, n)
       qi, qj := 0, 0
       q[qj] = src; qj++
       dist[src] = 0; par[src] = -1
       var last int
       for qi < qj {
           u := q[qi]; qi++
           last = u
           for _, w := range adj[u] {
               if dist[w] == -1 {
                   dist[w] = dist[u] + 1
                   par[w] = u
                   q[qj] = w; qj++
               }
           }
       }
       return last, dist, par
   }
   u, _, _ := far(0)
   v, _, par := far(u)
   // retrieve path u->v
   isOnPath := make([]bool, n)
   path := []int{}
   for x := v; x != -1; x = par[x] {
       path = append(path, x)
   }
   // path now v->u, reverse to u->v
   for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
       path[i], path[j] = path[j], path[i]
   }
   for _, x := range path {
       isOnPath[x] = true
   }
   visited := make([]bool, n)
   // mark path nodes visited to skip
   for _, x := range path {
       visited[x] = true
   }
   // compute side component sizes for each path node
   m := len(path)
   sList := make([]int64, m)
   var sumS2 int64
   stack := make([]int, 0, n)
   for i, vtx := range path {
       var compSize int64
       for _, w := range adj[vtx] {
           if !visited[w] {
               // explore this side component
               stack = stack[:0]
               visited[w] = true
               stack = append(stack, w)
               var cnt int64
               for len(stack) > 0 {
                   x := stack[len(stack)-1]
                   stack = stack[:len(stack)-1]
                   cnt++
                   for _, y := range adj[x] {
                       if !visited[y] {
                           visited[y] = true
                           stack = append(stack, y)
                       }
                   }
               }
               compSize += cnt
           }
       }
       s := compSize + 1
       sList[i] = s
       sumS2 += s * s
   }
   nn := int64(n)
   // total pairs across components
   sumPairs := (nn*nn - sumS2) / 2
   // exclude endpoint pair detour of length 1 edge
   extra := sumPairs - sList[0]*sList[m-1]
   // original simple paths of length>=2 edges: total pairs minus edges
   orig := nn*(nn-1)/2 - int64(n-1)
   ans := orig + extra
   fmt.Fprintln(out, ans)
}

package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var N, M int
       fmt.Fscan(reader, &N, &M)
       A := make([]int, M)
       B := make([]int, M)
       for i := 0; i < M; i++ {
           fmt.Fscan(reader, &A[i], &B[i])
           A[i]--
           B[i]--
       }
       edges := make([][]int, N)
       for i := 0; i < N; i++ {
           edges[i] = []int{}
       }
       for i := 0; i < M; i++ {
           u, v := A[i], B[i]
           edges[u] = append(edges[u], v)
           edges[v] = append(edges[v], u)
       }
       for i := 0; i < N; i++ {
           sort.Ints(edges[i])
       }
       // linked list of unvisited nodes
       next := make([]int, N)
       prev := make([]int, N)
       inUnvis := make([]bool, N)
       for i := 0; i < N; i++ {
           next[i] = i + 1
           prev[i] = i - 1
           inUnvis[i] = true
       }
       if N > 0 {
           next[N-1] = -1
       }
       head := 0
       var remove = func(u int) {
           inUnvis[u] = false
           if u == head {
               head = next[u]
           }
           if prev[u] != -1 {
               next[prev[u]] = next[u]
           }
           if next[u] != -1 {
               prev[next[u]] = prev[u]
           }
       }
       P := make([]int, N)
       Q := make([]int, N)
       num := 1
       var dfs func(n, u, ls int) bool
       dfs = func(n, u, ls int) bool {
           k := 0
           ok := false
           var nx []int
           var nd []int
           for x := head; x != -1; x = next[x] {
               for k < len(edges[n]) && edges[n][k] < x {
                   k++
               }
               if k < len(edges[n]) && edges[n][k] == x {
                   continue
               }
               nx = append(nx, x)
           }
           if ls == 1 {
               nd = append(nd, u)
           }
           for _, x := range nx {
               remove(x)
           }
           for i, v := range nx {
               if i == len(nx)-1 && len(nd) == 0 {
                   dfs(v, n, 1)
                   ok = true
               } else {
                   if dfs(v, n, 0) {
                       nd = append(nd, v)
                   }
               }
           }
           if len(nd) > 0 {
               P[n] = num
               for i, v := range nd {
                   P[v] = num + i + 1
               }
               Q[n] = num + len(nd)
               for i, v := range nd {
                   Q[v] = num + i
               }
               num += len(nd) + 1
               ok = true
           }
           if !ok && u == -1 {
               P[n] = num
               Q[n] = num
               num++
               ok = true
           }
           return !ok
       }
       for i := 0; i < N; i++ {
           if inUnvis[i] {
               remove(i)
               dfs(i, -1, 0)
           }
       }
       // output P
       for i := 0; i < N; i++ {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, P[i])
       }
       writer.WriteByte('\n')
       // output Q
       for i := 0; i < N; i++ {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, Q[i])
       }
       writer.WriteByte('\n')
   }
}

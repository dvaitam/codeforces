package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, m, k int
   G        [][]int
   st       []int
   dfn      []int
   col      []int
   col1     []int
   col2     []int
)

func dfs2(u, fa int) {
   st = append(st, u)
   dfn[u] = len(st)
   for _, v := range G[u] {
       if v == fa {
           continue
       }
       if dfn[v] == 0 {
           dfs2(v, u)
       } else {
           length := dfn[u] - dfn[v] + 1
           if length <= k {
               // found cycle
               fmt.Println(2)
               fmt.Println(length)
               for i := dfn[v] - 1; i < dfn[u]; i++ {
                   fmt.Printf("%d ", st[i])
               }
               fmt.Println()
               os.Exit(0)
           }
       }
   }
   st = st[:len(st)-1]
}

func dfs1(u, fa int) {
   if col[u] == 1 {
       col1 = append(col1, u)
   } else {
       col2 = append(col2, u)
   }
   need := (k + 1) / 2
   if len(col1) >= need {
       fmt.Println(1)
       for i := 0; i < need; i++ {
           fmt.Printf("%d ", col1[i])
       }
       fmt.Println()
       os.Exit(0)
   }
   if len(col2) >= need {
       fmt.Println(1)
       for i := 0; i < need; i++ {
           fmt.Printf("%d ", col2[i])
       }
       fmt.Println()
       os.Exit(0)
   }
   for _, v := range G[u] {
       if v == fa || col[v] != 0 {
           continue
       }
       col[v] = 3 - col[u]
       dfs1(v, u)
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   // fast input
   _, _ = fmt.Fscan(reader, &n, &m, &k)
   G = make([][]int, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       _, _ = fmt.Fscan(reader, &u, &v)
       G[u] = append(G[u], v)
       G[v] = append(G[v], u)
   }
   st = make([]int, 0, n)
   dfn = make([]int, n+1)
   // try find short cycle
   dfs2(1, 0)
   // no short cycle, output independent set
   fmt.Println(1)
   col = make([]int, n+1)
   col[1] = 1
   dfs1(1, 0)
}

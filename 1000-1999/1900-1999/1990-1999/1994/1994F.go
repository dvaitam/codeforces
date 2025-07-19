package main

import (
   "bufio"
   "fmt"
   "os"
)

type edge struct {
   v, id int
}

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

func nextInt() int {
   var b byte
   var err error
   // skip non-numeric
   for {
       b, err = reader.ReadByte()
       if err != nil {
           return 0
       }
       if (b >= '0' && b <= '9') || b == '-' {
           break
       }
   }
   sign := 1
   if b == '-' {
       sign = -1
       b, _ = reader.ReadByte()
   }
   x := 0
   for ; b >= '0' && b <= '9'; b, _ = reader.ReadByte() {
       x = x*10 + int(b-'0')
   }
   return x * sign
}

func main() {
   defer writer.Flush()
   t := nextInt()
   for tc := 0; tc < t; tc++ {
       n := nextInt()
       m := nextInt()
       finalG := make([][]edge, n)
       optionalG := make([][]edge, n)
       for i := 0; i < m; i++ {
           u := nextInt() - 1
           v := nextInt() - 1
           c := nextInt()
           if c == 1 {
               finalG[u] = append(finalG[u], edge{v, i})
               finalG[v] = append(finalG[v], edge{u, i})
           } else {
               optionalG[u] = append(optionalG[u], edge{v, i})
               optionalG[v] = append(optionalG[v], edge{u, i})
           }
       }
       trav := make([]bool, n)
       var dfs func(int) int
       dfs = func(u int) int {
           trav[u] = true
           odd := len(finalG[u]) & 1
           for _, e := range optionalG[u] {
               if trav[e.v] {
                   continue
               }
               if dfs(e.v) == 1 {
                   odd ^= 1
                   finalG[u] = append(finalG[u], e)
                   finalG[e.v] = append(finalG[e.v], edge{u, e.id})
               }
           }
           return odd
       }
       hasSol := true
       for i := 0; i < n; i++ {
           if !trav[i] {
               if dfs(i) == 1 {
                   hasSol = false
               }
           }
       }
       if !hasSol {
           fmt.Fprintln(writer, "NO")
           continue
       }
       fmt.Fprintln(writer, "YES")
       used := make([]bool, m)
       ptr := make([]int, n)
       var ans []int
       var solve func(int)
       solve = func(u int) {
           for ptr[u] < len(finalG[u]) {
               e := finalG[u][ptr[u]]
               ptr[u]++
               if used[e.id] {
                   continue
               }
               used[e.id] = true
               solve(e.v)
           }
           ans = append(ans, u+1)
       }
       solve(0)
       // output edges count = len(ans)-1
       fmt.Fprintln(writer, len(ans)-1)
       for i, v := range ans {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
       writer.WriteByte('\n')
   }
}

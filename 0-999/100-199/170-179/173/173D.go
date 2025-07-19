package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   rdr := bufio.NewReader(os.Stdin)
   wr := bufio.NewWriter(os.Stdout)
   defer wr.Flush()

   var n, m int
   fmt.Fscan(rdr, &n, &m)
   adj := make([][]int, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(rdr, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   col := make([]int, n+1)
   for i := 1; i <= n; i++ {
       col[i] = -1
   }
   // 2-coloring by BFS
   queue := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       if col[i] != -1 {
           continue
       }
       col[i] = 0
       queue = queue[:0]
       queue = append(queue, i)
       for qi := 0; qi < len(queue); qi++ {
           u := queue[qi]
           for _, v := range adj[u] {
               if col[v] == -1 {
                   col[v] = 1 - col[u]
                   queue = append(queue, v)
               }
           }
       }
   }
   // count colors
   cnt := [2]int{}
   for i := 1; i <= n; i++ {
       cnt[col[i]]++
   }
   // degrees
   du := make([]int, n+1)
   for i := 1; i <= n; i++ {
       du[i] = len(adj[i])
   }
   idx := make([]int, n+1)
   tot := 0

   var paint func(x int)
   paint = func(x int) {
       tot++
       idx[x] = tot
       vist := make([]bool, n+1)
       for _, v := range adj[x] {
           vist[v] = true
       }
       t := 0
       for i := 1; i <= n && t < 2; i++ {
           if idx[i] == 0 && col[i] != col[x] && !vist[i] {
               tot++
               idx[i] = tot
               t++
           }
       }
   }

   var output func()
   output = func() {
       fmt.Fprintln(wr, "YES")
       // assign remaining
       for i := 1; i <= n; i++ {
           if idx[i] == 0 && col[i] == 0 {
               tot++
               idx[i] = tot
           }
       }
       for i := 1; i <= n; i++ {
           if idx[i] == 0 && col[i] == 1 {
               tot++
               idx[i] = tot
           }
       }
       // print groups
       for i := 1; i <= n; i++ {
           grp := (idx[i]-1)/3 + 1
           fmt.Fprint(wr, grp)
           if i < n {
               fmt.Fprint(wr, " ")
           }
       }
       fmt.Fprintln(wr)
   }

   if cnt[0]%3 == 0 {
       output()
       return
   }
   if cnt[0]%3 == 2 {
       // flip colors
       for i := 1; i <= n; i++ {
           col[i] = 1 - col[i]
       }
       cnt[0], cnt[1] = cnt[1], cnt[0]
   }
   // try single paint on color 0
   for u := 1; u <= n; u++ {
       if col[u] != 0 {
           continue
       }
       if du[u]+2 <= cnt[1] {
           paint(u)
           output()
           return
       }
   }
   // try two paints on color 1
   tcnt := 0
   for u := 1; u <= n; u++ {
       if col[u] != 1 {
           continue
       }
       if du[u]+2 <= cnt[0] {
           paint(u)
           tcnt++
           if tcnt == 2 {
               output()
               return
           }
       }
   }
   fmt.Fprintln(wr, "NO")
}

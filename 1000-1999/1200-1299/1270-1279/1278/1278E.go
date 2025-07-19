package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

func readInt() int {
   sign := 1
   c, err := reader.ReadByte()
   for err == nil && (c < '0' || c > '9') {
       if c == '-' {
           sign = -1
       }
       c, err = reader.ReadByte()
   }
   x := 0
   for err == nil && c >= '0' && c <= '9' {
       x = x*10 + int(c-'0')
       c, err = reader.ReadByte()
   }
   return x * sign
}

func main() {
   defer writer.Flush()
   n := readInt()
   adj := make([][]int, n+1)
   for i := 1; i < n; i++ {
       u := readInt()
       v := readInt()
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   sz := make([]int, n+1)
   sz2 := make([]int, n+1)
   l := make([]int, n+1)
   r := make([]int, n+1)

   var dfs1 func(u, p int)
   dfs1 = func(u, p int) {
       sz[u] = 1
       for _, v := range adj[u] {
           if v == p {
               continue
           }
           dfs1(v, u)
           sz[u] += sz[v]
           sz2[u]++
       }
   }
   dfs1(1, 0)

   total := 2 * n
   r[1] = total
   l[1] = total - sz2[1] - 1

   var dfs2 func(u, p int)
   dfs2 = func(u, p int) {
       nw2 := l[u] + 1
       nw := l[u] - 1
       for _, v := range adj[u] {
           if v == p {
               continue
           }
           r[v] = nw2
           nw2++
           l[v] = nw - sz2[v]
           nw = nw - sz[v]*2 + 1
       }
       for _, v := range adj[u] {
           if v == p {
               continue
           }
           dfs2(v, u)
       }
   }
   dfs2(1, 0)

   for i := 1; i <= n; i++ {
       fmt.Fprintf(writer, "%d %d\n", l[i], r[i])
   }
}

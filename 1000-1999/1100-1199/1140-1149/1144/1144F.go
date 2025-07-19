package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   rdr := bufio.NewReader(os.Stdin)
   wrt := bufio.NewWriter(os.Stdout)
   defer wrt.Flush()

   // fast integer read
   readInt := func() int {
       var c byte
       var err error
       // skip non-numeric
       for {
           c, err = rdr.ReadByte()
           if err != nil {
               return 0
           }
           if (c >= '0' && c <= '9') || c == '-' {
               break
           }
       }
       neg := false
       if c == '-' {
           neg = true
           c, _ = rdr.ReadByte()
       }
       x := 0
       for ; c >= '0' && c <= '9'; c, _ = rdr.ReadByte() {
           x = x*10 + int(c - '0')
       }
       if neg {
           return -x
       }
       return x
   }

   n := readInt()
   m := readInt()
   adj := make([][]int, n+1)
   edges := make([][2]int, m)
   for i := 0; i < m; i++ {
       u := readInt()
       v := readInt()
       edges[i] = [2]int{u, v}
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }

   col := make([]int, n+1)
   // BFS for bipartiteness
   queue := make([]int, 0, n)
   head := 0
   col[1] = 1
   queue = append(queue, 1)
   for head < len(queue) {
       u := queue[head]
       head++
       for _, v := range adj[u] {
           if col[v] == 0 {
               col[v] = 3 - col[u]
               queue = append(queue, v)
           } else if col[v] == col[u] {
               fmt.Fprintln(wrt, "NO")
               return
           }
       }
   }

   // output
   fmt.Fprintln(wrt, "YES")
   res := make([]byte, m)
   for i, e := range edges {
       if col[e[0]] == 1 {
           res[i] = '0'
       } else {
           res[i] = '1'
       }
   }
   wrt.Write(res)
   wrt.WriteByte('\n')
}

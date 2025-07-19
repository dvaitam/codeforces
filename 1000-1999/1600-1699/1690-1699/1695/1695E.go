package main

import (
   "bufio"
   "fmt"
   "os"
)

// Edge represents a neighbor and edge index
type Edge struct {
   to, idx int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   m := n
   totalV := 2 * m

   // build graph
   graph := make([][]Edge, totalV)
   taken := make([]bool, m)
   vis := make([]bool, totalV)

   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       x--
       y--
       graph[x] = append(graph[x], Edge{to: y, idx: i})
       if x != y {
           graph[y] = append(graph[y], Edge{to: x, idx: i})
       }
   }

   // dfs to extract cycles
   var cycles [][]int
   var dfs func(c *[]int, u int)
   dfs = func(c *[]int, u int) {
       *c = append(*c, u)
       if vis[u] {
           return
       }
       vis[u] = true
       for _, e := range graph[u] {
           if !taken[e.idx] {
               taken[e.idx] = true
               dfs(c, e.to)
               *c = append(*c, u)
           }
       }
   }
   for i := 0; i < totalV; i++ {
       if !vis[i] && len(graph[i]) > 0 {
           c := make([]int, 0)
           dfs(&c, i)
           if len(c) > 0 {
               c = c[:len(c)-1]
               cycles = append(cycles, c)
           }
       }
   }

   // prepare results
   r0 := make([]int, 0)
   r1 := make([]int, 0)
   a := [2][]byte{{}, {}}
   b := [2][]byte{{}, {}}

   for _, c := range cycles {
       sz := len(c)
       if sz <= 2 {
           fmt.Fprintln(writer, -1)
           return
       }
       h := sz / 2
       for j := 0; j < h; j++ {
           r0 = append(r0, c[j])
           r1 = append(r1, c[sz-1-j])
           if j%2 == 0 {
               a[0] = append(a[0], 'R')
               a[1] = append(a[1], 'R')
               b[0] = append(b[0], 'L')
               b[1] = append(b[1], 'L')
           } else {
               a[0] = append(a[0], 'L')
               a[1] = append(a[1], 'L')
               b[0] = append(b[0], 'R')
               b[1] = append(b[1], 'R')
           }
           // first adjustment
           if j == 0 {
               a[0][len(a[0])-1] = 'U'
               a[1][len(a[1])-1] = 'D'
           }
           // last adjustment
           if j == h-1 {
               if h%2 == 0 {
                   a[0][len(a[0])-1] = 'U'
                   a[1][len(a[1])-1] = 'D'
               } else {
                   b[0][len(b[0])-1] = 'U'
                   b[1][len(b[1])-1] = 'D'
               }
           }
       }
   }

   // output
   fmt.Fprintln(writer, 2, len(r0))
   fmt.Fprintln(writer)
   // sequences (1-based)
   for i, v := range r0 {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v+1)
   }
   fmt.Fprintln(writer)
   for i, v := range r1 {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v+1)
   }
   fmt.Fprintln(writer)
   fmt.Fprintln(writer)
   // movements
   writer.Write(a[0])
   fmt.Fprintln(writer)
   writer.Write(a[1])
   fmt.Fprintln(writer)
   fmt.Fprintln(writer)
   writer.Write(b[0])
   fmt.Fprintln(writer)
   writer.Write(b[1])
   fmt.Fprintln(writer)
}

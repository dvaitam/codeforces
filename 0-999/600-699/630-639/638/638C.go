package main

import (
   "bufio"
   "fmt"
   "os"
)

type Edge struct {
   to, idx int
}

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
   n      int
   adj    [][]Edge
   ans    [][]int
   mx     int
)

// fast read integer
func readInt() int {
   x := 0; sign := 1
   b, err := reader.ReadByte()
   for err == nil && (b < '0' || b > '9') {
       if b == '-' {
           sign = -1
       }
       b, err = reader.ReadByte()
   }
   for err == nil && b >= '0' && b <= '9' {
       x = x*10 + int(b-'0')
       b, err = reader.ReadByte()
   }
   return x * sign
}

func dfs(u, last, parent int) {
   tot := 0
   for _, e := range adj[u] {
       v := e.to
       if v == parent {
           continue
       }
       tot++
       if tot == last {
           tot++
       }
       if tot > mx {
           mx = tot
       }
       ans[tot] = append(ans[tot], e.idx)
       dfs(v, tot, u)
   }
}

func main() {
   defer writer.Flush()
   n = readInt()
   adj = make([][]Edge, n+1)
   // ans size at most n
   ans = make([][]int, n+2)
   for i := 1; i < n; i++ {
       u := readInt()
       v := readInt()
       adj[u] = append(adj[u], Edge{to: v, idx: i})
       adj[v] = append(adj[v], Edge{to: u, idx: i})
   }
   dfs(1, 0, 0)
   // output
   fmt.Fprintln(writer, mx)
   for i := 1; i <= mx; i++ {
       line := ans[i]
       // print count and elements
       writer.WriteString(fmt.Sprintf("%d", len(line)))
       for _, id := range line {
           writer.WriteByte(' ')
           writer.WriteString(fmt.Sprintf("%d", id))
       }
       writer.WriteByte('\n')
   }
}

package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const mod = 998244353

var rdr = bufio.NewReader(os.Stdin)
var wr = bufio.NewWriter(os.Stdout)

// readInt reads next integer from standard input.
func readInt() int {
   var c byte
   var x int
   var neg bool
   for {
       b, err := rdr.ReadByte()
       if err != nil {
           return x
       }
       c = b
       if (c >= '0' && c <= '9') || c == '-' {
           break
       }
   }
   if c == '-' {
       neg = true
   } else {
       x = int(c - '0')
   }
   for {
       b, err := rdr.ReadByte()
       if err != nil {
           break
       }
       c = b
       if c < '0' || c > '9' {
           break
       }
       x = x*10 + int(c-'0')
   }
   if neg {
       return -x
   }
   return x
}

// Node holds value and coordinates.
type Node struct {
   v, x, y int
}

func main() {
   defer wr.Flush()
   n := readInt()
   m := readInt()
   nm := n * m
   nodes := make([]Node, nm)
   idx := 0
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           nodes[idx].v = readInt()
           nodes[idx].x = i
           nodes[idx].y = j
           idx++
       }
   }
   sort.Slice(nodes, func(i, j int) bool {
       return nodes[i].v < nodes[j].v
   })
   inv := make([]int, nm+1)
   if nm >= 0 {
       inv[0] = 1
   }
   if nm >= 1 {
       inv[1] = 1
   }
   for i := 2; i <= nm; i++ {
       inv[i] = inv[mod%i] * (mod - mod/i) % mod
   }
   ans := make([][]int, n+1)
   for i := range ans {
       ans[i] = make([]int, m+1)
   }
   var s, X, X2, Y, Y2, sum int
   for i := 0; i < nm; {
       j := i
       for j+1 < nm && nodes[j+1].v == nodes[i].v {
           j++
       }
       // compute answers for this group
       for k := i; k <= j; k++ {
           nd := nodes[k]
           cur := (s*nd.x*nd.x%mod + X2 + X*nd.x%mod + s*nd.y*nd.y%mod + Y2 + Y*nd.y%mod + sum) % mod
           ans[nd.x][nd.y] = cur * inv[s] % mod
       }
       // update accumulators
       for k := i; k <= j; k++ {
           nd := nodes[k]
           s++
           X = (X - 2*nd.x%mod + mod) % mod
           X2 = (X2 + nd.x*nd.x%mod) % mod
           Y = (Y - 2*nd.y%mod + mod) % mod
           Y2 = (Y2 + nd.y*nd.y%mod) % mod
           sum = (sum + ans[nd.x][nd.y]) % mod
       }
       i = j + 1
   }
   r := readInt()
   c := readInt()
   fmt.Fprintln(wr, ans[r][c])
}

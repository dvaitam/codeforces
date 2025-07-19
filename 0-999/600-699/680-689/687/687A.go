package main

import (
   "bufio"
   "fmt"
   "os"
)

var f []int

// find with path compression (iterative)
func find(x int) int {
   for f[x] != x {
       f[x] = f[f[x]]
       x = f[x]
   }
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   // initialize DSU for 1..2n
   f = make([]int, 2*n+1)
   for i := 1; i <= 2*n; i++ {
       f[i] = i
   }
   // read constraints
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       // x and y must be in opposite groups
       rx := find(x + n)
       ry := find(y)
       f[rx] = ry
       rx = find(y + n)
       ry = find(x)
       f[rx] = ry
   }
   // check for impossibility
   ok := true
   for i := 1; i <= n; i++ {
       if find(i) == find(i+n) {
           ok = false
           break
       }
   }
   if !ok {
       fmt.Fprintln(out, -1)
       return
   }
   // assign groups
   vis := make([]bool, 2*n+1)
   gr := make([]bool, 2*n+1)
   for i := 1; i <= n; i++ {
       t := find(i)
       rt := find(i + n)
       if !vis[t] && !vis[rt] {
           vis[t] = true
           vis[rt] = true
           gr[t] = true
       }
   }
   var v1, v2 []int
   for i := 1; i <= n; i++ {
       if gr[find(i)] {
           v1 = append(v1, i)
       } else {
           v2 = append(v2, i)
       }
   }
   // output
   fmt.Fprintln(out, len(v1))
   for i, x := range v1 {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, x)
   }
   fmt.Fprintln(out)
   fmt.Fprintln(out, len(v2))
   for i, x := range v2 {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, x)
   }
   fmt.Fprintln(out)
}

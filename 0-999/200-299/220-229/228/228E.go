package main

import (
   "bufio"
   "fmt"
   "os"
)

var n, m int
var par, dist []int
var use []bool

// fp finds the root of x with path compression,
// updating dist[x] to be the parity from x to its root.
func fp(x int) int {
   if par[x] == x {
       return x
   }
   pa := par[x]
   root := fp(pa)
   dist[x] ^= dist[pa]
   par[x] = root
   return root
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &n, &m)
   par = make([]int, n+1)
   dist = make([]int, n+1)
   use = make([]bool, n+1)
   for i := 1; i <= n; i++ {
       par[i] = i
       dist[i] = 0
   }
   for i := 0; i < m; i++ {
       var a, b, c int
       fmt.Fscan(reader, &a, &b, &c)
       // flip c for parity tracking
       c ^= 1
       pa := fp(a)
       pb := fp(b)
       if pa == pb {
           // check consistency
           if (dist[a]^dist[b]) != c {
               fmt.Fprintln(writer, "Impossible")
               return
           }
           continue
       }
       par[pa] = pb
       dist[pa] = c ^ dist[a] ^ dist[b]
   }
   cnt := 0
   for i := 1; i <= n; i++ {
       fp(i)
       if dist[i]&1 == 1 {
           use[i] = true
           cnt++
       }
   }
   fmt.Fprintln(writer, cnt)
   for i := 1; i <= n; i++ {
       if use[i] {
           cnt--
           if cnt > 0 {
               fmt.Fprintf(writer, "%d ", i)
           } else {
               fmt.Fprintf(writer, "%d\n", i)
           }
       }
   }
}

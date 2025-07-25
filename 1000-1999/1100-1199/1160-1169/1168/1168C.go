package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, q int
   fmt.Fscan(in, &n, &q)
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   const B = 19
   // f[i*B + b]: minimal index >= i reachable from i that has bit b set
   size := (n + 2) * B
   f := make([]int32, size)
   nex := make([]int32, B)
   inf := int32(n + 1)
   // initialize f and nex
   for b := 0; b < B; b++ {
       nex[b] = inf
   }
   // dp from right to left
   for i := n; i >= 1; i-- {
       off := i * B
       // init row to inf
       for b := 0; b < B; b++ {
           f[off+b] = inf
       }
       // list bits in a[i]
       ai := a[i]
       var bits []int
       for b := 0; b < B; b++ {
           if ai&(1<<b) != 0 {
               bits = append(bits, b)
               f[off+b] = int32(i)
           }
       }
       // merge from next positions
       for _, u := range bits {
           ni := nex[u]
           if ni > int32(n) {
               continue
           }
           noff := int(ni) * B
           for b := 0; b < B; b++ {
               if f[off+b] > f[noff+b] {
                   f[off+b] = f[noff+b]
               }
           }
       }
       // update nex
       for _, u := range bits {
           nex[u] = int32(i)
       }
   }
   // answer queries
   for qi := 0; qi < q; qi++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       offx := x * B
       ay := a[y]
       ok := false
       for b := 0; b < B; b++ {
           if ay&(1<<b) != 0 {
               if f[offx+b] <= int32(y) {
                   ok = true
                   break
               }
           }
       }
       if ok {
           fmt.Fprintln(out, "Shi")
       } else {
           fmt.Fprintln(out, "Fou")
       }
   }
}

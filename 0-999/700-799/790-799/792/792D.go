package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n uint64
   var q int
   if _, err := fmt.Fscan(reader, &n, &q); err != nil {
       return
   }
   // compute tree height H (leaves at depth H)
   // n+1 is power of two: let exp = log2(n+1), then H = exp - 1
   exp := bits.Len64(n+1) - 1
   H := exp - 1
   out := make([]uint64, q)
   for i := 0; i < q; i++ {
       var u uint64
       var s string
       fmt.Fscan(reader, &u)
       fmt.Fscan(reader, &s)
       // build path from root to u: 0 for L, 1 for R
       path := make([]uint8, 0, H)
       l, r := uint64(1), n
       h := H
       for {
           mid := (l + r) >> 1
           if u == mid {
               break
           } else if u < mid {
               path = append(path, 0)
               r = mid - 1
           } else {
               path = append(path, 1)
               l = mid + 1
           }
           h--
       }
       // height of subtree at current node
       curH := h
       // process commands
       for _, c := range s {
           switch c {
           case 'L':
               if curH > 0 {
                   path = append(path, 0)
                   curH--
               }
           case 'R':
               if curH > 0 {
                   path = append(path, 1)
                   curH--
               }
           case 'U':
               if len(path) > 0 {
                   // move up
                   path = path[:len(path)-1]
                   curH++
               }
           }
       }
       // compute label at final path
       l, r = 1, n
       for _, b := range path {
           mid := (l + r) >> 1
           if b == 0 {
               // left
               r = mid - 1
           } else {
               // right
               l = mid + 1
           }
       }
       // final node is root of [l,r]
       out[i] = (l + r) >> 1
   }
   // output
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for _, v := range out {
       fmt.Fprintln(w, v)
   }
}

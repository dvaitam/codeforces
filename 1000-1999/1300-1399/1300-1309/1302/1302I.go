package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   var s string
   fmt.Fscan(reader, &s)
   if k > n {
       fmt.Fprintln(writer, 0)
       return
   }

   // rolling hash parameters
   const base1 = 91138233
   const base2 = 97266353
   h1 := make([]uint64, n+1)
   h2 := make([]uint64, n+1)
   p1 := make([]uint64, n+1)
   p2 := make([]uint64, n+1)
   p1[0], p2[0] = 1, 1
   for i := 1; i <= n; i++ {
       p1[i] = p1[i-1] * base1
       p2[i] = p2[i-1] * base2
       var c uint64 = uint64(s[i-1]-'0') + 1
       h1[i] = h1[i-1]*base1 + c
       h2[i] = h2[i-1]*base2 + c
   }

   seen := make(map[[2]uint64]struct{})
   for i := k; i <= n; i++ {
       l := i - k + 1
       x1 := h1[i] - h1[l-1]*p1[k]
       x2 := h2[i] - h2[l-1]*p2[k]
       seen[[2]uint64{x1, x2}] = struct{}{}
   }

   fmt.Fprintln(writer, len(seen))
}

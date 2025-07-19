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

   var n, d int
   if _, err := fmt.Fscan(in, &n, &d); err != nil {
       return
   }
   s := make([]int, n)
   p := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &s[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &p[i])
   }

   cnt := 0
   score := s[d-1] + p[0]
   // pointers pl unused as in original logic
   pr := n - 1
   for i := 0; i < d-1; i++ {
       if s[i] + p[pr] > score {
           cnt++
       } else {
           pr--
       }
   }
   fmt.Fprintln(out, cnt+1)
}

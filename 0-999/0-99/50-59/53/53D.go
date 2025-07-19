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

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n)
   b := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &b[i])
   }
   // record swap positions
   swp := make([]int, 0, n*n)
   for i := 0; i < n; i++ {
       if a[i] != b[i] {
           j := i
           for ; j < n; j++ {
               if a[i] == b[j] {
                   break
               }
           }
           for k := j; k > i; k-- {
               swp = append(swp, k)
               b[k], b[k-1] = b[k-1], b[k]
           }
       }
   }
   // output
   fmt.Fprintln(out, len(swp))
   for _, k := range swp {
       fmt.Fprintln(out, k, k+1)
   }
}

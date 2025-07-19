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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       ops := make([]int, 0, 2*n)

       // helper to check non-decreasing
       isSorted := func() bool {
           for i := 1; i < n; i++ {
               if a[i] < a[i-1] {
                   return false
               }
           }
           return true
       }

       for !isSorted() {
           // compute mex
           present := make([]bool, n+1)
           for _, v := range a {
               if v >= 0 && v <= n {
                   present[v] = true
               }
           }
           mex := 0
           for mex <= n && present[mex] {
               mex++
           }
           if mex < n {
               // set a[mex] = mex
               a[mex] = mex
               ops = append(ops, mex+1)
           } else {
               // mex == n, find first i where a[i] != i
               idx := -1
               for i := 0; i < n; i++ {
                   if a[i] != i {
                       idx = i
                       break
                   }
               }
               if idx == -1 {
                   break
               }
               a[idx] = mex
               ops = append(ops, idx+1)
           }
       }
       // output
       fmt.Fprintln(writer, len(ops))
       for i, p := range ops {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, p)
       }
       writer.WriteByte('\n')
   }
}

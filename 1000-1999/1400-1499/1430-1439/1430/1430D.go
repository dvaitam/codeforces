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
   fmt.Fscan(reader, &t)
   for ; t > 0; t-- {
       var n int
       var s string
       fmt.Fscan(reader, &n, &s)
       // build run lengths
       runs := make([]int, 0, n)
       last := byte(0)
       for i := 0; i < n; i++ {
           if i == 0 || s[i] != last {
               runs = append(runs, 1)
               last = s[i]
           } else {
               runs[len(runs)-1]++
           }
       }
       ops := 0
       idx := 0
       m := len(runs)
       for idx < m {
           if idx == m-1 {
               // only one run left
               ops++
               break
           }
           if runs[idx] > 1 {
               // remove one char from this run, then delete rest of this run only
               ops++
               idx++
           } else {
               // runs[idx] == 1: remove this run and one char from next run
               ops++
               runs[idx+1]--
               idx++
               // if next run became empty, skip it too
               if idx < m && runs[idx] == 0 {
                   idx++
               }
           }
       }
       fmt.Fprintln(writer, ops)
   }
}

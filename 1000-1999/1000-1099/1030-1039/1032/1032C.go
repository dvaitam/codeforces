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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // dp arrays
   dp := [6]bool{}
   for f := 1; f <= 5; f++ {
       dp[f] = true
   }
   parent := make([][6]int, n+1)

   for i := 2; i <= n; i++ {
       var dp2 [6]bool
       for f := 1; f <= 5; f++ {
           if !dp[f] {
               continue
           }
           // try assign f2 at position i
           for f2 := 1; f2 <= 5; f2++ {
               ok := false
               if a[i-1] < a[i] {
                   if f < f2 {
                       ok = true
                   }
               } else if a[i-1] > a[i] {
                   if f > f2 {
                       ok = true
                   }
               } else {
                   if f != f2 {
                       ok = true
                   }
               }
               if ok && !dp2[f2] {
                   dp2[f2] = true
                   parent[i][f2] = f
               }
           }
       }
       dp = dp2
       // if no possible fingering at i
       any := false
       for f := 1; f <= 5; f++ {
           if dp[f] {
               any = true
               break
           }
       }
       if !any {
           fmt.Fprintln(writer, -1)
           return
       }
   }
   // construct answer
   b := make([]int, n+1)
   // pick any valid f at n
   for f := 1; f <= 5; f++ {
       if dp[f] {
           b[n] = f
           break
       }
   }
   // backtrack
   for i := n; i >= 2; i-- {
       b[i-1] = parent[i][b[i]]
   }
   // output
   for i := 1; i <= n; i++ {
       if i > 1 {
           writer.WriteByte(' ')
       }
       fmt.Fprintf(writer, "%d", b[i])
   }
   writer.WriteByte('\n')
}

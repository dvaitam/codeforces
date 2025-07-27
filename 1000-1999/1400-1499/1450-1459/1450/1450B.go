package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs(a int) int {
   if a < 0 {
       return -a
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n, k int
       fmt.Fscan(reader, &n, &k)
       xs := make([]int, n)
       ys := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &xs[i], &ys[i])
       }
       ans := -1
       for i := 0; i < n; i++ {
           ok := true
           xi, yi := xs[i], ys[i]
           for j := 0; j < n; j++ {
               if abs(xi-xs[j])+abs(yi-ys[j]) > k {
                   ok = false
                   break
               }
           }
           if ok {
               ans = 1
               break
           }
       }
       fmt.Fprintln(writer, ans)
   }
}

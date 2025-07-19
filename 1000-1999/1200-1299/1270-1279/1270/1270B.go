package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var T int
   fmt.Fscan(in, &T)
   for T > 0 {
       T--
       var n int
       fmt.Fscan(in, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &a[i])
       }
       found := false
       for i := 1; i < n; i++ {
           if abs(a[i]-a[i-1]) >= 2 {
               fmt.Fprintln(out, "YES")
               fmt.Fprintln(out, i, i+1)
               found = true
               break
           }
       }
       if !found {
           fmt.Fprintln(out, "NO")
       }
   }
}

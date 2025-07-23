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
   xs := make([]int64, n)
   hs := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &xs[i], &hs[i])
   }
   if n == 1 {
       fmt.Fprintln(writer, 1)
       return
   }
   ans := 1
   last := xs[0]
   for i := 1; i < n-1; i++ {
       if xs[i]-hs[i] > last {
           ans++
           last = xs[i]
       } else if xs[i]+hs[i] < xs[i+1] {
           ans++
           last = xs[i] + hs[i]
       } else {
           last = xs[i]
       }
   }
   ans++
   fmt.Fprintln(writer, ans)
}

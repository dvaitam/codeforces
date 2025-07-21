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
   var sum int64
   var maxa int64
   for i := 0; i < n; i++ {
       var a int64
       fmt.Fscan(reader, &a)
       sum += a
       if a > maxa {
           maxa = a
       }
   }
   // Each round has n-1 players, total player slots = rounds*(n-1)
   // Need rounds*(n-1) >= sum => rounds >= ceil(sum/(n-1))
   rounds := (sum + int64(n-1) - 1) / int64(n-1)
   if rounds < maxa {
       rounds = maxa
   }
   fmt.Fprintln(writer, rounds)
}

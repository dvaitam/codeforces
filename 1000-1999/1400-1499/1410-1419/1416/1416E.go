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
   var T int
   fmt.Fscan(reader, &T)
   for T > 0 {
       T--
       var n int
       fmt.Fscan(reader, &n)
       M := 0
       var r int
       for i := 0; i < n; i++ {
           var a int
           fmt.Fscan(reader, &a)
           if r > 0 && r < a {
               v := a - r
               if v == r {
                   M += 2
               } else {
                   M++
               }
               r = v
           } else {
               if a%2 == 0 {
                   M++
                   r = a / 2
               } else {
                   // split odd as (a-1,1) to favor future crosses
                   r = 1
               }
           }
       }
       // total b length = 2n - M
       ans := 2*n - M
       fmt.Fprintln(writer, ans)
   }
}

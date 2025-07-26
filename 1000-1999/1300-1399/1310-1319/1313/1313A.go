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
       var a, b, c int
       fmt.Fscan(reader, &a, &b, &c)
       // masks sorted by increasing number of dishes to maximize number of visitors
       masks := []int{1, 2, 4, 3, 5, 6, 7}
       cnt := []int{a, b, c}
       res := 0
       for _, m := range masks {
           ok := true
           // check availability
           for i := 0; i < 3; i++ {
               if (m&(1<<i)) != 0 && cnt[i] <= 0 {
                   ok = false
                   break
               }
           }
           if ok {
               // assign this mask
               res++
               for i := 0; i < 3; i++ {
                   if (m&(1<<i)) != 0 {
                       cnt[i]--
                   }
               }
           }
       }
       fmt.Fprintln(writer, res)
   }
}

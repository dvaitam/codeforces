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
   _, err := fmt.Fscan(reader, &t)
   if err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var X int64
       fmt.Fscan(reader, &X)
       // find highest bit k
       k := 0
       for i := 60; i >= 0; i-- {
           if (X>>i)&1 == 1 {
               k = i
               break
           }
       }
       // build answer
       p := 200 - k + 1
       var ans []int
       for i := 200 - k + 1; i <= 200; i++ {
           ans = append(ans, i)
           if (X>>(200-i))&1 == 1 {
               p--
               ans = append(ans, p)
           }
       }
       // output
       fmt.Fprintln(writer, len(ans))
       for i, v := range ans {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
       fmt.Fprintln(writer)
   }
}

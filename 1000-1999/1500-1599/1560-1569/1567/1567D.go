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
       var s int64
       var n int
       fmt.Fscan(reader, &s, &n)
       // find highest power of 10 <= s
       base := int64(1)
       for base <= s {
           base *= 10
       }
       base /= 10

       res := make([]int64, 0, n)
       rem := s
       curBase := base
       left := n
       for left > 0 {
           if left == 1 {
               res = append(res, rem)
               break
           }
           if rem-curBase >= int64(left-1) {
               res = append(res, curBase)
               rem -= curBase
               left--
           } else {
               curBase /= 10
           }
       }
       // output
       for i, v := range res {
           if i > 0 {
               fmt.Fprint(writer, " ")
           }
           fmt.Fprint(writer, v)
       }
       fmt.Fprint(writer, "\n")
   }
}

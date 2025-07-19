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
       var l, r int64
       fmt.Fscan(reader, &l, &r)
       bestDiff := -1
       bestNum := l
       found := false
       for i := l; i <= r; i++ {
           x := i
           mini, maxi := 10, -1
           if x == 0 {
               mini = 0
               maxi = 0
           }
           for x > 0 {
               d := int(x % 10)
               if d < mini {
                   mini = d
               }
               if d > maxi {
                   maxi = d
               }
               x /= 10
           }
           diff := maxi - mini
           if diff == 9 {
               fmt.Fprintln(writer, i)
               found = true
               break
           }
           if diff > bestDiff {
               bestDiff = diff
               bestNum = i
           }
       }
       if !found {
           fmt.Fprintln(writer, bestNum)
       }
   }
}

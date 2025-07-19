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
   fmt.Fscan(reader, &n)
   arr := make([]int, 2*n)
   cnt := make([]int, 90)
   for i := 0; i < 2*n; i++ {
       fmt.Fscan(reader, &arr[i])
       cnt[arr[i]-10]++
   }
   assignCnt := make([]int, 90)
   t, g1, g2 := 0, 0, 0
   // Handle unique elements
   for i := 0; i < 90; i++ {
       if cnt[i] == 1 {
           if t == 0 {
               assignCnt[i] = 1
               g1++
           } else {
               g2++
           }
           cnt[i] = 0
           t = (t + 1) % 2
       }
   }
   // Handle elements with multiple occurrences
   for i := 0; i < 90; i++ {
       if cnt[i] > 1 {
           assignCnt[i] = cnt[i] / 2
           g1++
           g2++
           if cnt[i]%2 == 1 {
               if t == 0 {
                   assignCnt[i]++
               }
               t = (t + 1) % 2
           }
       }
   }
   // Output maximum product of distinct counts
   fmt.Fprintln(writer, g1*g2)
   // Output assignment for each element
   for i := 0; i < 2*n; i++ {
       idx := arr[i] - 10
       if assignCnt[idx] > 0 {
           fmt.Fprint(writer, "1 ")
           assignCnt[idx]--
       } else {
           fmt.Fprint(writer, "2 ")
       }
   }
   fmt.Fprintln(writer)
}

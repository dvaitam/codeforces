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
   var M int
   fmt.Fscan(reader, &n, &M)
   times := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &times[i])
   }

   // counts of previous students' times
   var cnt [101]int
   res := make([]int, n)
   for i := 0; i < n; i++ {
       t := times[i]
       L := M - t
       sum := 0
       kept := 0
       // take smallest times first
       for tv := 1; tv <= 100; tv++ {
           c := cnt[tv]
           if c == 0 {
               continue
           }
           total := c * tv
           if sum+total <= L {
               sum += total
               kept += c
           } else {
               // take maximal from this time
               rem := (L - sum) / tv
               if rem > 0 {
                   kept += rem
               }
               break
           }
       }
       // removals = previous count - kept
       res[i] = i - kept
       // include current student for future
       cnt[t]++
   }
   // output
   for i, v := range res {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(fmt.Sprintf("%d", v))
   }
   writer.WriteByte('\n')
}

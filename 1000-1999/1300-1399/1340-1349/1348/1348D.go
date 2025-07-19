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
   for tc := 0; tc < t; tc++ {
       var n int64
       fmt.Fscan(reader, &n)
       // determine minimum days
       days := 0
       for (1<<days)-1 < int(n) {
           days++
       }
       days--
       // output days
       fmt.Fprintln(writer, days)
       // remaining total after initial day1 mass
       rem := n - 1
       last := int64(1)
       // collect splits per day
       for i := days; i >= 1; i-- {
           // binary search m between last and last*2
           low := last
           high := last * 2
           var ans int64 = high
           for low <= high {
               mid := (low + high) / 2
               // minimum and maximum total mass contribution over i days
               mn := int64(i) * mid
               mx := (int64((1<<i)-1)) * mid
               if mn > rem {
                   high = mid - 1
               } else if mx < rem {
                   low = mid + 1
               } else {
                   ans = mid
                   break
               }
           }
           // number of splits is ans - last
           fmt.Fprint(writer, ans-last)
           if i > 1 {
               fmt.Fprint(writer, " ")
           }
           rem -= ans
           last = ans
       }
       fmt.Fprintln(writer)
   }
}

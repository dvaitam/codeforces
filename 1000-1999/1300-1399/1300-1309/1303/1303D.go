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
   fmt.Fscan(reader, &t)
   for tc := 0; tc < t; tc++ {
       var n int64
       var m int
       fmt.Fscan(reader, &n, &m)
       // count boxes by power of two
       const maxb = 61
       cnt := make([]int64, maxb)
       var sum int64
       for i := 0; i < m; i++ {
           var a int64
           fmt.Fscan(reader, &a)
           // determine bit index
           b := 0
           for (1<<uint(b)) != a {
               b++
           }
           cnt[b]++
           sum += a
       }
       if sum < n {
           fmt.Fprintln(writer, -1)
           continue
       }
       var ans int64
       var have int64
       // process bits
       for i := 0; i < maxb; i++ {
           have += cnt[i]
           if (n>>uint(i))&1 == 1 {
               if have > 0 {
                   have--
               } else {
                   // find higher bit to split
                   j := i + 1
                   for j < maxb && cnt[j] == 0 {
                       j++
                   }
                   // perform splits
                   for k := j; k > i; k-- {
                       cnt[k]--
                       cnt[k-1] += 2
                       ans++
                   }
                   // take one at bit i
                   have += cnt[i]
                   have--
               }
           }
           // carry to next bit
           have /= 2
       }
       fmt.Fprintln(writer, ans)
   }
}

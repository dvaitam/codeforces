package main

import (
   "bufio"
   "fmt"
   "os"
)

func minInt64(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var targetW int64
   if _, err := fmt.Fscan(reader, &targetW); err != nil {
       return
   }
   var cnt [9]int64
   for i := 1; i <= 8; i++ {
       fmt.Fscan(reader, &cnt[i])
   }
   // total sum of all weights
   var sum int64
   for i := 1; i <= 8; i++ {
       sum += cnt[i] * int64(i)
   }
   if sum <= targetW {
       fmt.Fprintln(writer, sum)
       return
   }
   // greedily fill to as close as possible
   var now int64
   var cnt2 [9]int64
   for i := 1; i <= 8; i++ {
       w := int64(i)
       if now + w*cnt[i] <= targetW {
           now += w * cnt[i]
           cnt2[i] = cnt[i]
           cnt[i] = 0
       } else {
           v := (targetW - now) / w
           if v < 0 {
               v = 0
           }
           now += w * v
           cnt[i] -= v
           cnt2[i] += v
       }
   }
   // dp over possible adjustments
   const AA = 1000
   var dp [2001]bool
   dp[AA] = true // offset 0
   // add remaining (cnt) and subtract used (cnt2)
   for v := 1; v <= 8; v++ {
       // add items: positive shifts
       times := cnt[v]
       if times > 100 {
           times = 100
       }
       for j := int64(0); j < times; j++ {
           for k := 900; k >= -900; k-- {
               if dp[k+AA] {
                   dp[k+v+AA] = true
               }
           }
       }
       // remove items: negative shifts
       times2 := cnt2[v]
       if times2 > 100 {
           times2 = 100
       }
       for j := int64(0); j < times2; j++ {
           for k := -900; k <= 900; k++ {
               if dp[k+AA] {
                   dp[k-int(v)+AA] = true
               }
           }
       }
   }
   // find best achievable sum
   var ans int64
   maxR := targetW - now
   if maxR < 0 {
       maxR = 0
   }
   if maxR > 2000 {
       maxR = 2000
   }
   for i := int64(0); i <= maxR; i++ {
       if dp[int(i)+AA] {
           ans = now + i
       }
   }
   fmt.Fprintln(writer, ans)
}

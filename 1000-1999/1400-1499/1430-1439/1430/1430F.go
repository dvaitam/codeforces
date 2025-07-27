package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var k int64
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   waves := make([]struct{l, r, a int64}, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &waves[i].l, &waves[i].r, &waves[i].a)
   }
   var sumA int64
   for _, w := range waves {
       sumA += w.a
   }
   dp := map[int64]int64{k: 0}
   for i, w := range waves {
       T := w.r - w.l + 1
       newdp := make(map[int64]int64)
       var gap int64
       if i > 0 {
           gap = w.l - waves[i-1].r
       }
       // helper for transitions
       try := func(b0, cost0 int64) {
           // compute shooting events needed
           var events int64
           if b0 > 0 {
               rem := w.a - b0
               if rem < 0 {
                   rem = 0
               }
               ev := (rem + k - 1) / k
               events = 1 + ev
           } else {
               events = (w.a + k - 1) / k
           }
           if events > T {
               return
           }
           // compute remaining bullets after wave
           var remNew int64
           if b0 >= w.a {
               remNew = b0 - w.a
           } else {
               remm := (w.a - b0) % k
               if remm == 0 {
                   remNew = 0
               } else {
                   remNew = k - remm
               }
           }
           if prev, ok := newdp[remNew]; !ok || cost0 < prev {
               newdp[remNew] = cost0
           }
       }
       for remPrev, costPrev := range dp {
           // no reload before wave
           try(remPrev, costPrev)
           // reload in downtime
           if gap > 0 {
               try(k, costPrev + remPrev)
           }
       }
       dp = newdp
       if len(dp) == 0 {
           fmt.Println(-1)
           return
       }
   }
   var best int64 = -1
   for _, cost := range dp {
       if best < 0 || cost < best {
           best = cost
       }
   }
   fmt.Println(sumA + best)
}

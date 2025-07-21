package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var la, lb int
   if _, err := fmt.Fscan(reader, &la, &lb); err != nil {
       return
   }
   a := make([]int, la)
   for i := 0; i < la; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // positions in b
   bpos := make([]int, 1000001)
   for i := range bpos {
       bpos[i] = -1
   }
   for i := 0; i < lb; i++ {
       var x int
       fmt.Fscan(reader, &x)
       bpos[x] = i
   }
   // build P array for two copies of a
   n := la * 2
   P := make([]int, n)
   for i := 0; i < n; i++ {
       pos := bpos[a[i%la]]
       P[i] = pos
   }
   best := 0
   l := 0
   sumDelta := 0
   // sliding window over P[0..n)
   for r := 0; r < n; r++ {
       if P[r] < 0 {
           // reset window
           l = r + 1
           sumDelta = 0
           continue
       }
       if r > l {
           // compute delta from r-1 to r
           prev := P[r-1]
           if prev < 0 {
               // start of valid run
           } else {
               // mod difference
               d := P[r] - prev
               if d < 0 {
                   d += lb
               }
               sumDelta += d
           }
       }
       // shrink if sumDelta >= lb or window too big
       for l < r && (sumDelta >= lb || r-l+1 > la) {
           // remove delta between l and l+1
           if P[l] >= 0 && P[l+1] >= 0 {
               d := P[l+1] - P[l]
               if d < 0 {
                   d += lb
               }
               sumDelta -= d
           }
           l++
       }
       // update best
       if r-l+1 > best {
           best = r - l + 1
       }
   }
   // best cannot exceed la or lb
   if best > la {
       best = la
   }
   if best > lb {
       best = lb
   }
   fmt.Println(best)
}

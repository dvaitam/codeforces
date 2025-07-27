package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var tc int
   fmt.Fscan(reader, &tc)
   for tc > 0 {
       tc--
       var n int
       fmt.Fscan(reader, &n)
       t := make([]int64, n+1)
       x := make([]int64, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &t[i], &x[i])
       }
       // sentinel for next time
       t[n] = 1<<62

       var startTime int64 = 0
       var startPos int64 = 0
       var dest int64 = 0
       var finishTime int64 = 0
       var ans int

       for i := 0; i < n; i++ {
           // if robot is idle at t[i], start new command
           if finishTime <= t[i] {
               startTime = t[i]
               startPos = dest
               dest = x[i]
               finishTime = startTime + abs(dest-startPos)
           }
           // compute interval [t[i], t[i+1]]
           L := t[i]
           R := t[i+1]
           if R > finishTime {
               R = finishTime
           }
           // find position at L and R
           var posL, posR int64
           // pos at time L
           if L <= startTime {
               posL = startPos
           } else if L >= finishTime {
               posL = dest
           } else {
               if dest >= startPos {
                   posL = startPos + (L - startTime)
               } else {
                   posL = startPos - (L - startTime)
               }
           }
           // pos at time R
           if R <= startTime {
               posR = startPos
           } else if R >= finishTime {
               posR = dest
           } else {
               if dest >= startPos {
                   posR = startPos + (R - startTime)
               } else {
                   posR = startPos - (R - startTime)
               }
           }
           // check if x[i] lies between posL and posR
           lo, hi := posL, posR
           if lo > hi {
               lo, hi = hi, lo
           }
           if x[i] >= lo && x[i] <= hi {
               ans++
           }
       }
       fmt.Fprintln(writer, ans)
   }
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func abs(a int64) int64 {
   if a < 0 {
       return -a
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n int
       var m int64
       fmt.Fscan(reader, &n, &m)
       fishSum, meatSum := int64(0), int64(0)
       minFish, maxFish := int64(0), int64(0)
       minMeat, maxMeat := int64(0), int64(0)
       outX := make([]int64, n)
       outY := make([]int64, n)
       for i := 0; i < n; i++ {
           var a, b int64
           fmt.Fscan(reader, &a, &b)
           fishSum += a
           meatSum += b
           // minimal required fish and meat per dish
           reqFish := max(m-b, 0)
           reqMeat := max(m-a, 0)
           // maximal possible fish and meat per dish
           maxFish += min(m, a)
           maxMeat += min(m, b)
           minFish += reqFish
           minMeat += reqMeat
           outX[i] = reqFish
           outY[i] = reqMeat
       }
       // leftover total fish and meat after minimal assignment
       fishSum -= minFish
       meatSum -= minMeat
       both := maxFish - minFish
       diff := abs(fishSum - meatSum)
       var best, eatFish int64
       if diff > both {
           best = diff - both
           if fishSum > meatSum {
               eatFish = both
           } else {
               eatFish = 0
           }
       } else {
           best = (both - diff) % 2
           if fishSum > meatSum {
               eatFish = diff + (both-diff)/2
           } else {
               eatFish = (both-diff)/2
           }
       }
       // construct per-dish allocations
       fmt.Fprintln(writer, best)
       for i := 0; i < n; i++ {
           left := m - outX[i] - outY[i]
           useFish := min(eatFish, left)
           outX[i] += useFish
           eatFish -= useFish
           left -= useFish
           outY[i] += left
           fmt.Fprintln(writer, outX[i], outY[i])
       }
   }
}

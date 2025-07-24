package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   fmt.Fscan(reader, &n, &k)
   c := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &c[i])
   }
   caps := make([]int, k)
   isCap := make([]bool, n+1)
   for i := 0; i < k; i++ {
       fmt.Fscan(reader, &caps[i])
       isCap[caps[i]] = true
   }
   // total sum of beauties
   var totalBeauty int64
   for i := 1; i <= n; i++ {
       totalBeauty += c[i]
   }
   // sum of capitals beauties and sum of squares
   var sumCap, sumCapSq int64
   for _, x := range caps {
       sumCap += c[x]
       sumCapSq += c[x] * c[x]
   }
   // cycle edges sum
   var cycleSum int64
   for i := 1; i < n; i++ {
       cycleSum += c[i] * c[i+1]
   }
   cycleSum += c[n] * c[1]
   // T1 = totalBeauty*sumCap - sumCapSq
   T1 := totalBeauty*sumCap - sumCapSq
   // T2: subtract neighbor contributions of capitals
   var T2 int64
   for _, x := range caps {
       prev := x - 1
       if prev < 1 {
           prev = n
       }
       next := x + 1
       if next > n {
           next = 1
       }
       T2 += c[x] * (c[prev] + c[next])
   }
   // capTotal = sum over capitals x of cx*(totalBeauty - cx - c[prev] - c[next])
   capTotal := T1 - T2
   // sum of all unordered cap-cap pairs
   capPairsSum := (sumCap*sumCap - sumCapSq) / 2
   // sum of cap-cap neighbor pairs in cycle
   var capNeighborSum int64
   for i := 1; i < n; i++ {
       if isCap[i] && isCap[i+1] {
           capNeighborSum += c[i] * c[i+1]
       }
   }
   if isCap[n] && isCap[1] {
       capNeighborSum += c[n] * c[1]
   }
   // cap contribution adjusting double counts
   capContribution := capTotal - capPairsSum + capNeighborSum
   // result
   result := cycleSum + capContribution
   fmt.Println(result)
}

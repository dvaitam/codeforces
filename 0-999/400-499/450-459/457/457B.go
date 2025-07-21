package main

import (
   "bufio"
   "fmt"
   "os"
)

func minU(a, b uint64) uint64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var m, n int
   if _, err := fmt.Fscan(reader, &m, &n); err != nil {
       return
   }
   var sumA, sumB uint64
   var maxA, maxB uint64
   for i := 0; i < m; i++ {
       var x uint64
       fmt.Fscan(reader, &x)
       sumA += x
       if x > maxA {
           maxA = x
       }
   }
   for i := 0; i < n; i++ {
       var x uint64
       fmt.Fscan(reader, &x)
       sumB += x
       if x > maxB {
           maxB = x
       }
   }
   // Strategy 1: gather to best A partition
   // cost = (sumA - maxA) + sumB
   cost1 := (sumA - maxA) + sumB
   // Strategy 2: gather to best B partition
   // cost = sumA + (sumB - maxB)
   cost2 := sumA + (sumB - maxB)
   // Strategy 3: broadcast B rows to all A partitions
   cost3 := sumB * uint64(m)
   // Strategy 4: broadcast A rows to all B partitions
   cost4 := sumA * uint64(n)
   // answer is minimal
   ans := cost1
   ans = minU(ans, cost2)
   ans = minU(ans, cost3)
   ans = minU(ans, cost4)
   fmt.Fprintln(writer, ans)

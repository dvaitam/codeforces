package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int64
   var k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   positions := make([]int64, 0, k)
   var sum int64
   for i := 0; i < k; i++ {
       var a, b int64
       fmt.Fscan(reader, &a, &b)
       // b > 0 by problem
       sum += b
       positions = append(positions, a)
   }
   // total sum
   if sum < n {
       fmt.Println(1)
       return
   }
   if sum > n {
       fmt.Println(-1)
       return
   }
   // sum == n
   // if no zeros, all positions have at least one coin -> all b_i == 1, game ends
   if int64(len(positions)) == n {
       fmt.Println(1)
       return
   }
   // find max zero-gap between positions with b>0
   sort.Slice(positions, func(i, j int) bool { return positions[i] < positions[j] })
   maxGap := int64(0)
   // consecutive
   for i := 0; i+1 < len(positions); i++ {
       gap := positions[i+1] - positions[i] - 1
       if gap > maxGap {
           maxGap = gap
       }
   }
   // wrap-around gap
   // between last and first around the table
   wrap := positions[0] + n - positions[len(positions)-1] - 1
   if wrap > maxGap {
       maxGap = wrap
   }
   // if any zero segment length >= 2, game is good
   if maxGap >= 2 {
       fmt.Println(1)
   } else {
       fmt.Println(-1)
   }
}

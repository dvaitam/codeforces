package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, l int
   fmt.Fscan(reader, &n, &l)
   nums := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &nums[i])
   }
   sort.Ints(nums)
   // initial max distance: from 0 to first, and last to l
   biggest := float64(max(nums[0], l-nums[n-1]))
   for i := 0; i+1 < n; i++ {
       diff := nums[i+1] - nums[i]
       d := float64(diff) / 2.0
       if d > biggest {
           biggest = d
       }
   }
   fmt.Printf("%.12f\n", biggest)
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

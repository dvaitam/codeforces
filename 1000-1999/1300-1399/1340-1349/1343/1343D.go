package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
      return a
   }
   return b
}

func max(a, b int) int {
   if a > b {
      return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for ; t > 0; t-- {
      var n, k int
      fmt.Fscan(reader, &n, &k)
      arr := make([]int, n)
      for i := 0; i < n; i++ {
         fmt.Fscan(reader, &arr[i])
      }
      pairs := n / 2
      size := 2*k + 5
      cnt := make([]int, size)
      diff := make([]int, size+1)
      for i := 0; i < pairs; i++ {
         x := arr[i]
         y := arr[n-1-i]
         sum := x + y
         cnt[sum]++
         l := min(x, y) + 1
         r := max(x, y) + k
         diff[l]++
         diff[r+1]--
      }
      ones := make([]int, size)
      cur := 0
      for s := 2; s <= 2*k; s++ {
         cur += diff[s]
         ones[s] = cur
      }
      ans := 2 * pairs
      for s := 2; s <= 2*k; s++ {
         ops := 2*pairs - ones[s] - cnt[s]
         if ops < ans {
            ans = ops
         }
      }
      fmt.Fprintln(writer, ans)
   }
}

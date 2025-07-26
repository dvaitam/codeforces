package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
      return
   }
   arr := make([]int, n)
   maxVal := 0
   for i := 0; i < n; i++ {
      fmt.Fscan(reader, &arr[i])
      if arr[i] > maxVal {
         maxVal = arr[i]
      }
   }
   cnt := make([]int, maxVal+1)
   for _, v := range arr {
      cnt[v]++
   }
   freq := make([]int64, maxVal+1)
   for g := 1; g <= maxVal; g++ {
      for j := g; j <= maxVal; j += g {
         freq[g] += int64(cnt[j])
      }
   }
   dp := make([]int64, maxVal+1)
   for g := maxVal; g >= 1; g-- {
      dp[g] = int64(g) * freq[g]
      for j := g * 2; j <= maxVal; j += g {
         val := dp[j] + (freq[g]-freq[j])*int64(g)
         if val > dp[g] {
            dp[g] = val
         }
      }
   }
   fmt.Println(dp[1])
}


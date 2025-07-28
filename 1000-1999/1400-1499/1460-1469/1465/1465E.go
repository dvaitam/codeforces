package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   var T int64
   fmt.Fscan(reader, &n, &T)
   var S string
   fmt.Fscan(reader, &S)
   // Count values for first n-1 characters
   cnt := make([]int64, 26)
   for i := 0; i < n-1; i++ {
       pos := int(S[i] - 'a')
       cnt[pos]++
   }
   // Value of the last character
   last := int(S[n-1] - 'a')
   wN := int64(1) << last
   // Total sum of first n-1 values
   var sum int64
   for j := 0; j < 26; j++ {
       sum += cnt[j] << j
   }
   // Need sum_{i<n} s_i * w_i = X = T - wN, with s_i in {-1,+1}
   X := T - wN
   // sum s_i * w_i = X <=> let P = sum of w_i with s_i=+1,
   // P - (sum-P) = X => 2P = X + sum => P = (X+sum)/2
   if (sum+X)&1 != 0 || sum+X < 0 {
       fmt.Fprintln(writer, "No")
       return
   }
   K := (sum + X) / 2
   // Greedy subset sum with powers of two
   for j := 25; j >= 0; j-- {
       if K <= 0 {
           break
       }
       take := min64(cnt[j], K>>j)
       K -= take << j
   }
   if K == 0 {
       fmt.Fprintln(writer, "Yes")
   } else {
       fmt.Fprintln(writer, "No")
   }
}

// min64 returns the smaller of a or b.
func min64(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

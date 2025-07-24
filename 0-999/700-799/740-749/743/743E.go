package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   seq := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &seq[i])
       seq[i]-- // 0-based
   }
   // positions for each number 0..7
   pos := make([][]int, 8)
   for i, v := range seq {
       if v >= 0 && v < 8 {
           pos[v] = append(pos[v], i+1) // use 1-based for dp
       }
   }
   // answer case1: one of each
   ans1 := 0
   for i := 0; i < 8; i++ {
       if len(pos[i]) > 0 {
           ans1++
       }
   }
   ans := ans1
   // case2: all numbers at least once
   if ans1 == 8 {
       // dp[mask][d] = min end index (1-based) after selecting blocks in mask, with d doubles used
       // mask bits 0..7
       INF := n + 1
       dp := make([][9]int, 1<<8)
       for m := range dp {
           for d := range dp[m] {
               dp[m][d] = INF
           }
       }
       dp[0][0] = 0
       // iterate masks
       for mask := 0; mask < (1<<8); mask++ {
           t := bitsOnes(mask)
           for d := 0; d <= t; d++ {
               cur := dp[mask][d]
               if cur > n {
                   continue
               }
               // try next number i
               for i := 0; i < 8; i++ {
                   if mask&(1<<i) != 0 {
                       continue
                   }
                   // pick one
                   arr := pos[i]
                   // find first > cur
                   j := sort.Search(len(arr), func(k int) bool { return arr[k] > cur })
                   if j < len(arr) {
                       m2 := mask | (1 << i)
                       if arr[j] < dp[m2][d] {
                           dp[m2][d] = arr[j]
                       }
                   }
                   // pick two
                   if len(arr) >= 2 {
                       // first > cur
                       j1 := sort.Search(len(arr), func(k int) bool { return arr[k] > cur })
                       if j1 < len(arr) {
                           // second > arr[j1]
                           j2 := sort.Search(len(arr), func(k int) bool { return arr[k] > arr[j1] })
                           if j2 < len(arr) {
                               m2 := mask | (1 << i)
                               if arr[j2] < dp[m2][d+1] {
                                   dp[m2][d+1] = arr[j2]
                               }
                           }
                       }
                   }
               }
           }
       }
       full := (1<<8 - 1)
       // best with full mask
       for d := 0; d <= 8; d++ {
           if dp[full][d] <= n {
               tot := 8 + d
               if tot > ans {
                   ans = tot
               }
           }
       }
   }
   // output
   fmt.Println(ans)
}

// bitsOnes returns the popcount of x
func bitsOnes(x int) int {
   // builtin popcount workaround
   cnt := 0
   for x > 0 {
       cnt++
       x &= x - 1
   }
   return cnt
}

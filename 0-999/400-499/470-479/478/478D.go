package main

import (
   "fmt"
)

const mod = 1000000007

func main() {
   var r, g int
   if _, err := fmt.Scan(&r, &g); err != nil {
       return
   }
   total := r + g
   // find maximum height h: h*(h+1)/2 <= total
   h := int(( -1 + intSqrt(1+8*total) ) / 2)
   // compute total blocks needed for h levels
   S := h * (h + 1) / 2
   // choose smaller group to reduce dp size
   small, big := r, g
   if small > big {
       small, big = big, small
   }
   // bounds for sum of small-colored levels
   low := S - big
   if low < 0 {
       low = 0
   }
   // dp[s] = number of ways to choose subset of levels with sum s
   dp := make([]int, small+1)
   dp[0] = 1
   for i := 1; i <= h; i++ {
       // update dp in reverse to avoid overwrite
       if i > small {
           continue
       }
       for s := small; s >= i; s-- {
           dp[s] += dp[s-i]
           if dp[s] >= mod {
               dp[s] -= mod
           }
       }
   }
   // sum valid ways
   ans := 0
   for s := low; s <= small; s++ {
       ans += dp[s]
       if ans >= mod {
           ans -= mod
       }
   }
   fmt.Println(ans)
}

// intSqrt returns floor(sqrt(x)) for non-negative x
func intSqrt(x int) int {
   // binary search
   lo, hi := 0, x
   for lo <= hi {
       mid := (lo + hi) / 2
       if mid*mid <= x {
           lo = mid + 1
       } else {
           hi = mid - 1
       }
   }
   return hi
}

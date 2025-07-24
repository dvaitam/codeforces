package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a, b, k, t int
   if _, err := fmt.Fscan(reader, &a, &b, &k, &t); err != nil {
       return
   }
   mod := int64(1e9 + 7)
   d := b - a
   // DP over difference sums
   N := 4 * k * t + 1
   base := 2 * k * t
   dp := make([]int64, N)
   newdp := make([]int64, N)
   s0 := make([]int64, N)
   s1 := make([]int64, N)
   // initial: sum = 0
   dp[base] = 1
   twok1 := int64(2*k + 1)
   for step := 1; step <= t; step++ {
       // previous sum range [prev_low, prev_high]
       prev_low := base - 2*k*(step-1)
       prev_high := base + 2*k*(step-1)
       // current sum range [curr_low, curr_high]
       curr_low := base - 2*k*step
       curr_high := base + 2*k*step
       // build prefix sums of dp
       s0[0] = dp[0]
       s1[0] = 0
       for i := 1; i < N; i++ {
           s0[i] = s0[i-1] + dp[i]
           if s0[i] >= mod {
               s0[i] -= mod
           }
           s1[i] = s1[i-1] + dp[i]*int64(i)%mod
           if s1[i] >= mod {
               s1[i] -= mod
           }
       }
       // zero newdp in current range
       for i := curr_low; i <= curr_high; i++ {
           newdp[i] = 0
       }
       // compute convolution with triangular weights
       for dIdx := curr_low; dIdx <= curr_high; dIdx++ {
           // window [L,R]
           L := dIdx - 2*k
           if L < 0 {
               L = 0
           }
           R := dIdx + 2*k
           if R >= N {
               R = N - 1
           }
           // sum of dp in [L,R]
           var sumPrev int64
           if L > 0 {
               sumPrev = s0[R] - s0[L-1]
           } else {
               sumPrev = s0[R]
           }
           if sumPrev < 0 {
               sumPrev += mod
           }
           // weighted sum W1 = sum dp[i] * |dIdx - i|
           // split into left [L..dIdx] and right [dIdx+1..R]
           // left part
           Rleft := dIdx
           if Rleft > R {
               Rleft = R
           }
           Lleft := L
           var sumLeft, sumS1Left int64
           if Rleft >= Lleft {
               if Lleft > 0 {
                   sumLeft = s0[Rleft] - s0[Lleft-1]
                   sumS1Left = s1[Rleft] - s1[Lleft-1]
               } else {
                   sumLeft = s0[Rleft]
                   sumS1Left = s1[Rleft]
               }
               if sumLeft < 0 {
                   sumLeft += mod
               }
               if sumS1Left < 0 {
                   sumS1Left += mod
               }
           }
           // right part
           Lright := dIdx + 1
           if Lright > R {
               // no right part, W1 = 0
               newdp[dIdx] = twok1 * sumPrev % mod
           } else {
               sumPrevRight := s0[R] - s0[dIdx]
               sumS1Right := s1[R] - s1[dIdx]
               if sumPrevRight < 0 {
                   sumPrevRight += mod
               }
               if sumS1Right < 0 {
                   sumS1Right += mod
               }
               part1 := (sumLeft*int64(dIdx)%mod - sumS1Left) % mod
               if part1 < 0 {
                   part1 += mod
               }
               part2 := (sumS1Right - sumPrevRight*int64(dIdx)%mod) % mod
               if part2 < 0 {
                   part2 += mod
               }
               W1 := part1 + part2
               if W1 >= mod {
                   W1 -= mod
               }
               val := twok1*sumPrev%mod - W1
               if val < 0 {
                   val += mod
               }
               newdp[dIdx] = val
           }
       }
       // swap dp and newdp
       dp, newdp = newdp, dp
   }
   // build prefix sum of final dp
   s0[0] = dp[0]
   for i := 1; i < N; i++ {
       s0[i] = s0[i-1] + dp[i]
       if s0[i] >= mod {
           s0[i] -= mod
       }
   }
   // sum dp for D > d
   start := base + d + 1
   if start < 0 {
       start = 0
   }
   if start > N-1 {
       fmt.Println(0)
       return
   }
   res := s0[N-1]
   if start > 0 {
       res = (res - s0[start-1] + mod) % mod
   }
   fmt.Println(res)
}

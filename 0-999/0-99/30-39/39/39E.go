package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var a, b int
   var n int64
   fmt.Fscan(in, &a, &b, &n)
   // Draw if infinite safe moves: only when a == 1
   if a == 1 {
       fmt.Println("Missing")
       return
   }
   // compute max b such that 2^b < n
   bmax := 0
   for ; int64(1)<<uint(bmax) < n; bmax++ {
   }
   bmax--
   if bmax < 1 {
       bmax = 1
   }
   // compute Amax[b] = max a such that a^b < n
   Amax := make([]int, bmax+2)
   for bi := 2; bi <= bmax; bi++ {
       Amax[bi] = computeAmax(n, bi)
   }
   // DP for b >= 2
   dp := make([][]bool, bmax+2)
   for bi := 2; bi <= bmax; bi++ {
       sz := Amax[bi] + 2
       dp[bi] = make([]bool, sz)
       // dp[bi][a] default false
   }
   // fill DP
   for bi := bmax; bi >= 2; bi-- {
       for ai := Amax[bi]; ai >= 2; ai-- {
           win := false
           // move a+1
           if int64(ai+1) < n && powInt(int64(ai+1), bi) < n {
               if ai+1 <= Amax[bi] {
                   if !dp[bi][ai+1] {
                       win = true
                   }
               } else {
                   // ai+1 beyond safe region? actually pow check ensures safe until ai+1 <=Amax
               }
           }
           // move b+1
           if !win {
               if bi+1 <= bmax && ai <= Amax[bi+1] {
                   if !dp[bi+1][ai] {
                       win = true
                   }
               }
           }
           dp[bi][ai] = win
       }
   }
   // prepare DP for b = 1 via DP1
   // limit where a^2 < n
   A2Limit := 0
   if 2 <= bmax {
       A2Limit = Amax[2]
   } else {
       // no b>=2 region, so only parity
       A2Limit = 0
   }
   // DP1 for a in [2..A2Limit]
   DP1 := make([]bool, A2Limit+2)
   for ai := A2Limit; ai >= 2; ai-- {
       win := false
       // move a+1
       if int64(ai+1) < n {
           var childWin bool
           if ai+1 <= A2Limit {
               childWin = DP1[ai+1]
           } else {
               // beyond DP1, only a-move chain
               L := (n - 1) - int64(ai+1)
               childWin = (L%2 != 0)
           }
           if !childWin {
               win = true
           }
       }
       // move to b=2
       if !win && 2 <= bmax && ai <= Amax[2] {
           if !dp[2][ai] {
               win = true
           }
       }
       DP1[ai] = win
   }
   // determine initial win
   var startWin bool
   if b >= 2 {
       if a <= Amax[b] {
           startWin = dp[b][a]
       } else {
           // should not happen as initial a^b<n
           startWin = false
       }
   } else {
       // b == 1
       if a <= A2Limit {
           startWin = DP1[a]
       } else {
           L := (n - 1) - int64(a)
           startWin = (L%2 != 0)
       }
   }
   if startWin {
       // Stas (first) wins => print Masha
       fmt.Println("Masha")
   } else {
       // Stas loses => print Stas
       fmt.Println("Stas")
   }
}

// computeAmax finds max a >=1 such that a^b < n
func computeAmax(n int64, b int) int {
   lo, hi := 1, int(math.Pow(float64(n-1), 1/float64(b)))+1
   if hi < 1 {
       hi = 1
   }
   for lo < hi {
       mid := (lo + hi + 1) >> 1
       if powInt(int64(mid), b) < n {
           lo = mid
       } else {
           hi = mid - 1
       }
   }
   return lo
}

// powInt computes x^b
func powInt(x int64, b int) int64 {
   res := int64(1)
   for i := 0; i < b; i++ {
       res *= x
       if res < 0 {
           return res
       }
   }
   return res
}

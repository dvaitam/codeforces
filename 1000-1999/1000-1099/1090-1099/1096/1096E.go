package main

import (
   "bufio"
   "fmt"
   "os"
)

const (
   mod = 998244353
   M   = 5220
   N   = 120
)

var C [M][N]int

// modPow computes a^b mod mod
func modPow(a, b int) int {
   res := 1
   for b > 0 {
       if b&1 == 1 {
           res = int((int64(res) * int64(a)) % mod)
       }
       a = int((int64(a) * int64(a)) % mod)
       b >>= 1
   }
   return res
}

func main() {
   // precompute combinations C[n][k]
   for i := 0; i < M; i++ {
       C[i][0] = 1
       for j := 1; j < N && j <= i; j++ {
           C[i][j] = C[i-1][j] + C[i-1][j-1]
           if C[i][j] >= mod {
               C[i][j] -= mod
           }
       }
   }

   in := bufio.NewReader(os.Stdin)
   var playerCount, sum, lower int
   fmt.Fscan(in, &playerCount, &sum, &lower)
   // special case: no lower bound means uniform win probability = 1/playerCount
   if lower == 0 {
       fmt.Println(modPow(playerCount, mod-2))
       return
   }
   totalWays := C[sum-lower+playerCount-1][playerCount-1]
   invTotal := modPow(totalWays, mod-2)
   ans := 0
   // ties: number of players tied for highest score
   for ties := 1; ties <= playerCount; ties++ {
       save := modPow(ties, mod-2) // probability to win among ties
       // topScore: score of each tied top player
       for topScore := lower; ties*topScore <= sum; topScore++ {
           rem := sum - ties*topScore
           t := playerCount - ties
           // no other players
           if t == 0 {
               if rem == 0 {
                   ans += save
                   if ans >= mod {
                       ans -= mod
                   }
               }
               continue
           }
           // count ways for others to have scores < topScore summing to rem
           res := 0
           // inclusion-exclusion on number of "bad" players with score >= topScore
           for bad := 0; bad <= t && bad*topScore <= rem; bad++ {
               nsum := rem - bad*topScore
               add := int((int64(C[nsum+t-1][t-1]) * int64(C[t][bad])) % mod)
               if bad&1 == 1 {
                   res = (res - add + mod) % mod
               } else {
                   res = (res + add) % mod
               }
           }
           // choose which players are tied
           res = int((int64(res) * int64(C[playerCount-1][ties-1])) % mod)
           ans = int((int64(ans) + int64(res)*int64(save)) % mod)
       }
   }
   // normalize by total number of states
   ans = int((int64(ans) * int64(invTotal)) % mod)
   fmt.Println(ans)
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func modpow(a, e int64) int64 {
   res := int64(1)
   for e > 0 {
       if e&1 == 1 {
           res = res * a % MOD
       }
       a = a * a % MOD
       e >>= 1
   }
   return res
}

func modinv(a int64) int64 {
   return modpow(a, MOD-2)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int
   var s int64
   fmt.Fscan(in, &n, &m, &k, &s)
   ptsMap := make(map[int64]bool)
   anomalies := make([]pair, 0, k)
   startAnom := false
   endAnom := false
   for i := 0; i < k; i++ {
       var r, c int
       fmt.Fscan(in, &r, &c)
       key := (int64(r)<<32 | int64(c))
       ptsMap[key] = true
       anomalies = append(anomalies, pair{r, c, true})
       if r == 1 && c == 1 {
           startAnom = true
       }
       if r == n && c == m {
           endAnom = true
       }
   }
   // include start and end if not in anomalies
   pts := make([]pair, 0, len(anomalies)+2)
   pts = append(pts, pair{1, 1, startAnom})
   for _, p := range anomalies {
       if (p.r == 1 && p.c == 1) || (p.r == n && p.c == m) {
           continue
       }
       pts = append(pts, p)
   }
   pts = append(pts, pair{n, m, endAnom})
   // sort by r then c
   sortPairs(pts)
   N := len(pts)
   // precompute factorials
   maxNM := n + m + 5
   fact := make([]int64, maxNM)
   invf := make([]int64, maxNM)
   fact[0] = 1
   for i := 1; i < maxNM; i++ {
       fact[i] = fact[i-1] * int64(i) % MOD
   }
   invf[maxNM-1] = modinv(fact[maxNM-1])
   for i := maxNM - 2; i >= 0; i-- {
       invf[i] = invf[i+1] * int64(i+1) % MOD
   }
   // reachable pairs
   reachable := make([][]bool, N)
   for i := range reachable {
       reachable[i] = make([]bool, N)
   }
   for i := 0; i < N; i++ {
       for j := 0; j < N; j++ {
           if pts[i].r <= pts[j].r && pts[i].c <= pts[j].c {
               reachable[i][j] = true
           }
       }
   }
   // compute waysAvoid[j][i] for j<i
   waysAvoid := make([][]int64, N)
   for j := 0; j < N; j++ {
       wa := make([]int64, N)
       for i := j + 1; i < N; i++ {
           if !reachable[j][i] {
               continue
           }
           dr := pts[i].r - pts[j].r
           dc := pts[i].c - pts[j].c
           w := comb(dr+dc, dr, fact, invf)
           // subtract intermediate
           for l := j + 1; l < i; l++ {
               if waysAvoid[j][l] != 0 && reachable[l][i] {
                   dr2 := pts[i].r - pts[l].r
                   dc2 := pts[i].c - pts[l].c
                   w = (w - waysAvoid[j][l]*comb(dr2+dc2, dr2, fact, invf)) % MOD
               }
           }
           if w < 0 {
               w += MOD
           }
           wa[i] = w
       }
       waysAvoid[j] = wa
   }
   // DP counts dp[i][t]
   // compute T_max
   T_max := 0
   tmp := s
   for tmp > 0 {
       tmp >>= 1
       T_max++
   }
   // include end anomaly reduces one more
   // dp as slice of slices
   dp := make([][]int64, N)
   for i := 0; i < N; i++ {
       dp[i] = make([]int64, T_max+2)
   }
   // initial at start index 0
   initT := 0
   if pts[0].anom {
       initT = 1
   }
   dp[0][initT] = 1
   // DP
   for i := 1; i < N; i++ {
       for j := 0; j < i; j++ {
           wji := waysAvoid[j][i]
           if wji == 0 {
               continue
           }
           for t := 0; t <= T_max; t++ {
               if dp[j][t] == 0 {
                   continue
               }
               nt := t
               if pts[i].anom {
                   nt++
               }
               if nt > T_max {
                   continue
               }
               dp[i][nt] = (dp[i][nt] + dp[j][t]*wji) % MOD
           }
       }
   }
   // total paths
   total := comb(int64(n-1+ m-1), int64(n-1), fact, invf)
   // expected numerator
   endIdx := N - 1
   var num int64
   for t := 0; t <= T_max; t++ {
       cnt := dp[endIdx][t]
       if cnt == 0 {
           continue
       }
       // value = floor(s/2^t)
       v := s >> t
       num = (num + cnt%MOD* (v % MOD)) % MOD
   }
   // answer = num * inv(total)
   ans := num * modinv(total) % MOD
   fmt.Println(ans)
}

// comb computes C(n,k)
func comb(n, k int64, fact, invf []int64) int64 {
   if k < 0 || k > n {
       return 0
   }
   return fact[n] * invf[k] % MOD * invf[n-k] % MOD
}

// pair holds a point
type pair struct{ r, c int; anom bool }

// sortPairs sorts points by r then c
func sortPairs(a []pair) {
   for i := 1; i < len(a); i++ {
       j := i
       for j > 0 && (a[j].r < a[j-1].r || (a[j].r == a[j-1].r && a[j].c < a[j-1].c)) {
           a[j], a[j-1] = a[j-1], a[j]
           j--
       }
   }
}

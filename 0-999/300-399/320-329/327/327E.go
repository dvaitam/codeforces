package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

const MOD = 1000000007

func add(a, b int64) int64 {
   a += b
   if a >= MOD {
       a -= MOD
   }
   return a
}

func mul(a, b int64) int64 {
   return (a * b) % MOD
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   var k int
   fmt.Fscan(in, &k)
   bad := make([]int64, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(in, &bad[i])
   }
   // total ways
   fac := make([]int64, n+1)
   fac[0] = 1
   for i := 1; i <= n; i++ {
       fac[i] = fac[i-1] * int64(i) % MOD
   }
   total := fac[n]
   if k == 0 {
       fmt.Println(total)
       return
   }
   // split halves
   n1 := n / 2
   n2 := n - n1
   aL := a[:n1]
   aR := a[n1:]
   // precompute all subset sums and size counts for meet-in-middle
   mapL := make(map[int64]map[int]int)
   for mask := 0; mask < (1 << n1); mask++ {
       var sum int64
       sz := bits.OnesCount(uint(mask))
       for j := 0; j < n1; j++ {
           if mask&(1<<j) != 0 {
               sum += aL[j]
           }
       }
       if mapL[sum] == nil {
           mapL[sum] = make(map[int]int)
       }
       mapL[sum][sz]++
   }
   mapR := make(map[int64]map[int]int)
   for mask := 0; mask < (1 << n2); mask++ {
       var sum int64
       sz := bits.OnesCount(uint(mask))
       for j := 0; j < n2; j++ {
           if mask&(1<<j) != 0 {
               sum += aR[j]
           }
       }
       if mapR[sum] == nil {
           mapR[sum] = make(map[int]int)
       }
       mapR[sum][sz]++
   }
   // helper compute f(x): number of permutations hitting sum x
   computeF := func(x int64) int64 {
       var res int64
       for sumL, cntL := range mapL {
           sumR := x - sumL
           cntR, ok := mapR[sumR]
           if !ok {
               continue
           }
           for sL, cL := range cntL {
               for sR, cR := range cntR {
                   s := sL + sR
                   // s! * (n-s)!
                   ways := fac[s] * fac[n-s] % MOD
                   res = (res + ways*int64(cL)*int64(cR)) % MOD
               }
           }
       }
       return res
   }
   // compute f(x) for each bad
   f := make([]int64, k)
   for i := 0; i < k; i++ {
       f[i] = computeF(bad[i])
   }
   var fxy int64
   if k == 2 {
       x, y := bad[0], bad[1]
       // ensure order x < y
       if x > y {
           x, y = y, x
       }
       // DP subsets of halves for per-mask subset count for sums <= y and sums <= x
       // precompute full mask sums and popcounts
       cntL := 1 << n1
       sumFullL := make([]int64, cntL)
       pcL := make([]int, cntL)
       for mask := 1; mask < cntL; mask++ {
           lsb := mask & -mask
           j := bits.TrailingZeros(uint(lsb))
           prev := mask ^ lsb
           sumFullL[mask] = sumFullL[prev] + aL[j]
           pcL[mask] = pcL[prev] + 1
       }
       cntR := 1 << n2
       sumFullR := make([]int64, cntR)
       pcR := make([]int, cntR)
       for mask := 1; mask < cntR; mask++ {
           lsb := mask & -mask
           j := bits.TrailingZeros(uint(lsb))
           prev := mask ^ lsb
           sumFullR[mask] = sumFullR[prev] + aR[j]
           pcR[mask] = pcR[prev] + 1
       }
       // build dpLeftCount and dpRightCount: for each mask, map sum->map size->count
       dpL := make([]map[int64]map[int]int, cntL)
       dpL[0] = map[int64]map[int]int{0: {0: 1}}
       for mask := 1; mask < cntL; mask++ {
           lsb := mask & -mask
           j := bits.TrailingZeros(uint(lsb))
           prev := mask ^ lsb
           prevMap := dpL[prev]
           curMap := make(map[int64]map[int]int, len(prevMap)+1)
           // copy
           for sum, msz := range prevMap {
               nm := make(map[int]int, len(msz))
               for sz, c := range msz {
                   nm[sz] = c
               }
               curMap[sum] = nm
           }
           // extend
           ai := aL[j]
           for sum, msz := range prevMap {
               nsum := sum + ai
               inn, ok := curMap[nsum]
               if !ok {
                   inn = make(map[int]int)
                   curMap[nsum] = inn
               }
               for sz, c := range msz {
                   inn[sz+1] += c
               }
           }
           dpL[mask] = curMap
       }
       dpR := make([]map[int64]map[int]int, cntR)
       dpR[0] = map[int64]map[int]int{0: {0: 1}}
       for mask := 1; mask < cntR; mask++ {
           lsb := mask & -mask
           j := bits.TrailingZeros(uint(lsb))
           prev := mask ^ lsb
           prevMap := dpR[prev]
           curMap := make(map[int64]map[int]int, len(prevMap)+1)
           for sum, msz := range prevMap {
               nm := make(map[int]int, len(msz))
               for sz, c := range msz {
                   nm[sz] = c
               }
               curMap[sum] = nm
           }
           ai := aR[j]
           for sum, msz := range prevMap {
               nsum := sum + ai
               inn, ok := curMap[nsum]
               if !ok {
                   inn = make(map[int]int)
                   curMap[nsum] = inn
               }
               for sz, c := range msz {
                   inn[sz+1] += c
               }
           }
           dpR[mask] = curMap
       }
       // precompute masks by full sums for R
       masksR := make(map[int64][]int)
       for mask := 0; mask < cntR; mask++ {
           s := sumFullR[mask]
           if s > y {
               continue
           }
           masksR[s] = append(masksR[s], mask)
       }
       // iterate T = (maskL,maskR) with sumFullL+sumFullR == y
       for maskL, sumL := range sumFullL {
           if sumL > y {
               continue
           }
           sumRNeeded := y - sumL
           listR, ok := masksR[sumRNeeded]
           if !ok {
               continue
           }
           tL := pcL[maskL]
           dpLM := dpL[maskL]
           for _, maskR := range listR {
               t := tL + pcR[maskR]
               // for this T, compute g(T)
               dpRM := dpR[maskR]
               // sum over subsets S of T with sum x
               for sumSL, mszL := range dpLM {
                   sumSR := x - sumSL
                   mszR, ok2 := dpRM[sumSR]
                   if !ok2 {
                       continue
                   }
                   for szL, cL := range mszL {
                       for szR, cR := range mszR {
                           s := szL + szR
                           // s! * (t-s)! * (n-t)!
                           ways := fac[s] * fac[t-s] % MOD * fac[n-t] % MOD
                           fxy = (fxy + ways*int64(cL)*int64(cR)) % MOD
                       }
                   }
               }
           }
       }
   }
   // inclusion-exclusion
   ans := total
   for i := 0; i < k; i++ {
       ans = (ans - f[i] + MOD) % MOD
   }
   if k == 2 {
       ans = (ans + fxy) % MOD
   }
   fmt.Println(ans)
}

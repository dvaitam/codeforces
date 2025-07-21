package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 51123987

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var s string
   fmt.Fscan(reader, &s)
   // Manacher's algorithm
   d1 := make([]int, n)
   l, r := 0, -1
   for i := 0; i < n; i++ {
       k := 1
       if i <= r {
           k = min(d1[l+r-i], r-i+1)
       }
       for i-k >= 0 && i+k < n && s[i-k] == s[i+k] {
           k++
       }
       d1[i] = k
       if i+k-1 > r {
           l = i - k + 1
           r = i + k - 1
       }
   }
   d2 := make([]int, n)
   l, r = 0, -1
   for i := 0; i < n; i++ {
       k := 0
       if i <= r {
           k = min(d2[l+r-i+1], r-i+1)
       }
       for i-k-1 >= 0 && i+k < n && s[i-k-1] == s[i+k] {
           k++
       }
       d2[i] = k
       if i+k-1 > r {
           l = i - k
           r = i + k - 1
       }
   }
   // difference arrays for start and end counts
   ds := make([]int64, n+2)
   de := make([]int64, n+2)
   var total int64
   for i := 0; i < n; i++ {
       // odd-length palindromes
       if d1[i] > 0 {
           total += int64(d1[i])
           L := i - (d1[i] - 1)
           R := i
           ds[L]++
           ds[R+1]--
           de[i]++
           de[i+(d1[i]-1)+1]--
       }
       // even-length palindromes
       if d2[i] > 0 {
           total += int64(d2[i])
           L := i - d2[i]
           R := i - 1
           if L <= R {
               ds[L]++
               ds[R+1]--
           }
           de[i]++
           de[i+d2[i]]--
       }
   }
   // build counts
   cntStart := make([]int64, n)
   cntEnd := make([]int64, n)
   var cur int64
   for i := 0; i < n; i++ {
       cur += ds[i]
       cntStart[i] = cur
   }
   cur = 0
   for i := 0; i < n; i++ {
       cur += de[i]
       cntEnd[i] = cur
   }
   // prefix sum of ends
   prefEnd := make([]int64, n)
   var cum int64
   for i := 0; i < n; i++ {
       cum = (cum + cntEnd[i]) % MOD
       prefEnd[i] = cum
   }
   // count disjoint pairs
   var disjoint int64
   for i := 0; i < n; i++ {
       cs := cntStart[i] % MOD
       if cs != 0 && i > 0 {
           disjoint = (disjoint + cs*prefEnd[i-1]) % MOD
       }
   }
   // total pairs
   totalMod := total % MOD
   inv2 := (MOD + 1) / 2
   t := totalMod * ((totalMod - 1 + MOD) % MOD) % MOD
   totalPairs := t * int64(inv2) % MOD
   ans := (totalPairs - disjoint + MOD) % MOD
   fmt.Println(ans)
}

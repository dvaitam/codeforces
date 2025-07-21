package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   var p int64
   fmt.Fscan(in, &n, &m, &p)
   ls := make([]int, n)
   maxL := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &ls[i])
       if ls[i] > maxL {
           maxL = ls[i]
       }
   }
   // Precompute combinatorics up to maxL
   fact := make([]int64, maxL+1)
   invFact := make([]int64, maxL+1)
   inv := make([]int64, maxL+1)
   fact[0] = 1
   for i := 1; i <= maxL; i++ {
       fact[i] = fact[i-1] * int64(i) % p
   }
   invFact[maxL] = modInv(fact[maxL], p)
   for i := maxL; i > 0; i-- {
       invFact[i-1] = invFact[i] * int64(i) % p
   }
   inv[1] = 1
   for i := 2; i <= maxL; i++ {
       inv[i] = p - (p/int64(i))*inv[p%int64(i)]%p
   }
   // falling factorial m*(m-1)*...*(m-k+1)
   fall := make([]int64, maxL+1)
   fall[0] = 1
   for k := 1; k <= maxL; k++ {
       fall[k] = fall[k-1] * (int64(m) - int64(k-1) + p) % p
   }
   // invC[k] = inverse of C(m,k)
   invC := make([]int64, maxL+1)
   for k := 0; k <= maxL; k++ {
       // C(m,k) = fall[k] * invFact[k]
       c := fall[k] * invFact[k] % p
       if c == 0 {
           invC[k] = 0
       } else {
           invC[k] = modInv(c, p)
       }
   }
   // Precompute gRows[L][k] = C(m,k)*f(L,k): number of sequences length L with exactly k distinct colors
   gRows := make([][]int64, maxL+1)
   // L=1
   gRows[1] = make([]int64, 2)
   if m > 0 {
       gRows[1][1] = int64(m) % p
   }
   for L := 2; L <= maxL; L++ {
       prev := gRows[L-1]
       row := make([]int64, L+1)
       // k from 1 to L
       for k := 1; k <= L; k++ {
           var v int64
           // extend by new color
           if k-1 >= 1 {
               v = (v + prev[k-1] * (int64(m) - int64(k-1) + p) % p) % p
           }
           // extend by existing
           if k <= L-1 {
               v = (v + prev[k] * int64(k-1) % p) % p
           }
           row[k] = v
       }
       gRows[L] = row
   }
   // DP over layers
   // dpPrevValue[k] = dp_{i-1}(S) for any S of size k, i.e. f(li-1,k)
   dpPrevValue := make([]int64, maxL+1)
   // init for first layer
   L0 := ls[0]
   prevRow := gRows[L0]
   var prevTotal int64
   for k := 1; k <= L0; k++ {
       prevTotal = (prevTotal + prevRow[k]) % p
       if invC[k] != 0 {
           dpPrevValue[k] = prevRow[k] * invC[k] % p
       } else {
           dpPrevValue[k] = 0
       }
   }
   // iterate next layers
   for i := 1; i < n; i++ {
       Li := ls[i]
       curRow := gRows[Li]
       // sum of curRow
       var sumCur int64
       for k := 1; k <= Li; k++ {
           sumCur = (sumCur + curRow[k]) % p
       }
       // compute sum of overlap for total
       var overlap int64
       // k upto min(prevL, Li)
       prevL := len(prevRow) - 1
       mx := prevL
       if Li < mx {
           mx = Li
       }
       for k := 1; k <= mx; k++ {
           overlap = (overlap + curRow[k] * dpPrevValue[k]) % p
       }
       curTotal := (prevTotal*sumCur - overlap) % p
       if curTotal < 0 {
           curTotal += p
       }
       // compute new dpPrevValue for this layer
       newDP := make([]int64, Li+1)
       for k := 1; k <= Li; k++ {
           var prevVal int64
           if k <= prevL {
               prevVal = dpPrevValue[k]
           }
           x := (prevTotal - prevVal) % p
           if x < 0 {
               x += p
           }
           // f = curRow[k] / C(m,k)
           if invC[k] != 0 {
               newDP[k] = curRow[k] * invC[k] % p * x % p
           } else {
               newDP[k] = 0
           }
       }
       prevRow = curRow
       prevTotal = curTotal
       dpPrevValue = newDP
   }
   // output
   fmt.Println(prevTotal)
}

func modInv(a, mod int64) int64 {
   // Fermat's little theorem, mod assumed prime
   return modPow(a, mod-2, mod)
}

func modPow(a, e, mod int64) int64 {
   res := int64(1)
   a %= mod
   for e > 0 {
       if e&1 != 0 {
           res = res * a % mod
       }
       a = a * a % mod
       e >>= 1
   }
   return res
}

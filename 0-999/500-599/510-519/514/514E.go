package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

// multiply two matrices a and b of size m x m
func matMul(a, b [][]int64) [][]int64 {
   m := len(a)
   c := make([][]int64, m)
   for i := 0; i < m; i++ {
       c[i] = make([]int64, m)
       for k := 0; k < m; k++ {
           if a[i][k] == 0 {
               continue
           }
           aik := a[i][k]
           for j := 0; j < m; j++ {
               c[i][j] = (c[i][j] + aik*b[k][j]) % MOD
           }
       }
   }
   return c
}

// multiply matrix a (m x m) with vector v (size m)
func matVec(a [][]int64, v []int64) []int64 {
   m := len(v)
   res := make([]int64, m)
   for i := 0; i < m; i++ {
       var sum int64
       for j := 0; j < m; j++ {
           sum = (sum + a[i][j]*v[j]) % MOD
       }
       res[i] = sum
   }
   return res
}

// fast exponentiation of matrix a to power e
func matPow(a [][]int64, e int64) [][]int64 {
   m := len(a)
   // init identity
   res := make([][]int64, m)
   for i := 0; i < m; i++ {
       res[i] = make([]int64, m)
       res[i][i] = 1
   }
   for e > 0 {
       if e&1 != 0 {
           res = matMul(res, a)
       }
       a = matMul(a, a)
       e >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var x int64
   fmt.Fscan(in, &n, &x)
   pMap := make(map[int]int)
   var d int
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &d)
       pMap[d]++
   }
   // max distance step
   D := 0
   for k := range pMap {
       if k > D {
           D = k
       }
   }
   // dp and prefix sum S for s < D
   dp := make([]int64, D)
   S := make([]int64, D)
   dp[0] = 1
   S[0] = 1
   for s := 1; s < D; s++ {
       var sum int64
       for k, cnt := range pMap {
           if k <= s {
               sum = (sum + int64(cnt)*dp[s-k]) % MOD
           }
       }
       dp[s] = sum
       S[s] = (S[s-1] + dp[s]) % MOD
   }
   if x < int64(D) {
       fmt.Println(S[x])
       return
   }
   // build transition matrix size m = D+1
   m := D + 1
   M := make([][]int64, m)
   for i := 0; i < m; i++ {
       M[i] = make([]int64, m)
   }
   // first row: dp
   for k, cnt := range pMap {
       // dp[s+1] includes p[k]*dp[s+1-k], maps to dp[s-(k-1)] at idx k-1
       idx := k - 1
       if idx >= 0 && idx < D {
           M[0][idx] = (M[0][idx] + int64(cnt)) % MOD
       }
   }
   // shift rows for dp
   for i := 1; i < D; i++ {
       M[i][i-1] = 1
   }
   // last row for S: S[s+1] = S[s] + dp[s+1]
   // copy first row to last
   for j := 0; j < D; j++ {
       M[D][j] = M[0][j]
   }
   // S shift
   M[D][D] = 1

   // initial vector V at s = D-1
   V := make([]int64, m)
   // dp[D-1], dp[D-2],...,dp[0]
   for i := 0; i < D; i++ {
       V[i] = dp[D-1-i]
   }
   V[D] = S[D-1]
   // exponentiate M to power (x-(D-1))
   e := x - int64(D-1)
   Mexp := matPow(M, e)
   Vres := matVec(Mexp, V)
   // result is S[x] at index D
   fmt.Println(Vres[D])
}

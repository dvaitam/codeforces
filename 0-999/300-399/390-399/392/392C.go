package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

const mod = 1000000007

// multiply two matrices a and b (size dxd)
func matMul(a, b [][]int64, d int) [][]int64 {
   res := make([][]int64, d)
   for i := 0; i < d; i++ {
       res[i] = make([]int64, d)
       for k := 0; k < d; k++ {
           if a[i][k] == 0 {
               continue
           }
           aik := a[i][k]
           for j := 0; j < d; j++ {
               res[i][j] = (res[i][j] + aik*b[k][j]) % mod
           }
       }
   }
   return res
}

// matrix exponentiation: raise m to power p
func matPow(m [][]int64, d int, p uint64) [][]int64 {
   // initialize result as identity
   res := make([][]int64, d)
   for i := 0; i < d; i++ {
       res[i] = make([]int64, d)
       res[i][i] = 1
   }
   for p > 0 {
       if p&1 == 1 {
           res = matMul(res, m, d)
       }
       m = matMul(m, m, d)
       p >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var ns, ks string
   if _, err := fmt.Fscan(in, &ns, &ks); err != nil {
       return
   }
   n, _ := strconv.ParseUint(ns, 10, 64)
   k, _ := strconv.Atoi(ks)
   // sequence defined for n>=1
   // handle n==1: sum = F1*1^k = 1
   if n == 1 {
       fmt.Println(1)
       return
   }
   // precompute binomial coefficients up to k
   C := make([][]int64, k+1)
   for i := 0; i <= k; i++ {
       C[i] = make([]int64, i+1)
       C[i][0], C[i][i] = 1, 1
       for j := 1; j < i; j++ {
           C[i][j] = (C[i-1][j-1] + C[i-1][j]) % mod
       }
   }
   // dimensions
   kp := k + 1
   d := 3 * kp
   offsetP := 0
   offsetQ := kp
   offsetT := 2 * kp
   // build transition matrix M
   M := make([][]int64, d)
   for i := range M {
       M[i] = make([]int64, d)
   }
   // P_t' and Q_t' and T_j'
   for t := 0; t <= k; t++ {
       // row for P_t'
       rP := offsetP + t
       for s := 0; s <= t; s++ {
           coeff := C[t][s]
           M[rP][offsetP+s] = (M[rP][offsetP+s] + coeff) % mod
           M[rP][offsetQ+s] = (M[rP][offsetQ+s] + coeff) % mod
       }
       // row for Q_t'
       rQ := offsetQ + t
       for s := 0; s <= t; s++ {
           coeff := C[t][s]
           M[rQ][offsetP+s] = (M[rQ][offsetP+s] + coeff) % mod
       }
       // row for T_t'
       rT := offsetT + t
       // keep old T_t
       M[rT][rT] = 1
       // add P_t' contribution
       for s := 0; s <= t; s++ {
           coeff := C[t][s]
           M[rT][offsetP+s] = (M[rT][offsetP+s] + coeff) % mod
           M[rT][offsetQ+s] = (M[rT][offsetQ+s] + coeff) % mod
       }
   }
   // initial vector v1 for n=1
   v1 := make([]int64, d)
   // F(1)=1, F(0)=1, n=1
   for t := 0; t <= k; t++ {
       // 1^t = 1
       v1[offsetP+t] = 1 // F(1)*1^t
       v1[offsetQ+t] = 1 // F(0)*1^t
       v1[offsetT+t] = 1 // T_t(1) = F1*1^t
   }
   // compute M^(n-1)
   Mpow := matPow(M, d, n-1)
   // multiply Mpow * v1
   var resVal int64
   rk := offsetT + k
   for j := 0; j < d; j++ {
       resVal = (resVal + Mpow[rk][j]*v1[j]) % mod
   }
   fmt.Println(resVal)
}

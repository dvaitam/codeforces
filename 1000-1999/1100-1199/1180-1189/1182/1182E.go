package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007
const MOD1 = MOD - 1

// 3x3 matrix
type Mat3 [3][3]int64
// 5x5 matrix
type Mat5 [5][5]int64

// multiply two 3x3 matrices modulo MOD1
func mulMat3(a, b Mat3) Mat3 {
   var c Mat3
   for i := 0; i < 3; i++ {
       for j := 0; j < 3; j++ {
           var sum int64
           for k := 0; k < 3; k++ {
               sum += a[i][k] * b[k][j]
           }
           c[i][j] = sum % MOD1
       }
   }
   return c
}

// exponentiate 3x3 matrix to power e modulo MOD1
func powMat3(a Mat3, e int64) Mat3 {
   // initialize to identity
   var res Mat3
   for i := 0; i < 3; i++ {
       res[i][i] = 1
   }
   for e > 0 {
       if e&1 == 1 {
           res = mulMat3(res, a)
       }
       a = mulMat3(a, a)
       e >>= 1
   }
   return res
}

// multiply two 5x5 matrices modulo MOD1
func mulMat5(a, b Mat5) Mat5 {
   var c Mat5
   for i := 0; i < 5; i++ {
       for j := 0; j < 5; j++ {
           var sum int64
           for k := 0; k < 5; k++ {
               sum += a[i][k] * b[k][j]
           }
           c[i][j] = sum % MOD1
       }
   }
   return c
}

// exponentiate 5x5 matrix to power e modulo MOD1
func powMat5(a Mat5, e int64) Mat5 {
   var res Mat5
   for i := 0; i < 5; i++ {
       res[i][i] = 1
   }
   for e > 0 {
       if e&1 == 1 {
           res = mulMat5(res, a)
       }
       a = mulMat5(a, a)
       e >>= 1
   }
   return res
}

// modular exponentiation
func modPow(a, e, mod int64) int64 {
   a %= mod
   var r int64 = 1
   for e > 0 {
       if e&1 == 1 {
           r = (r * a) % mod
       }
       a = (a * a) % mod
       e >>= 1
   }
   return r
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, f1, f2, f3, c int64
   fmt.Fscan(in, &n, &f1, &f2, &f3, &c)
   if n == 1 {
       fmt.Println(f1 % MOD)
       return
   }
   if n == 2 {
       fmt.Println(f2 % MOD)
       return
   }
   if n == 3 {
       fmt.Println(f3 % MOD)
       return
   }
   // tribonacci for e1, e2, e3
   // base matrix
   var A3 Mat3
   A3[0] = [3]int64{1, 1, 1}
   A3[1] = [3]int64{1, 0, 0}
   A3[2] = [3]int64{0, 1, 0}
   // compute A3^(n-3)
   m3 := powMat3(A3, n-3)
   // e1_n = m3[0][2], e2_n = m3[0][1], e3_n = m3[0][0]
   e1 := m3[0][2]
   e2 := m3[0][1]
   e3 := m3[0][0]
   // compute p_n via 5x5 matrix
   var M5 Mat5
   // row0: [1,1,1,2,-4]
   M5[0] = [5]int64{1, 1, 1, 2, MOD1 - 4}
   // row1: [1,0,0,0,0]
   M5[1] = [5]int64{1, 0, 0, 0, 0}
   // row2: [0,1,0,0,0]
   M5[2] = [5]int64{0, 1, 0, 0, 0}
   // row3: [0,0,0,1,1]
   M5[3] = [5]int64{0, 0, 0, 1, 1}
   // row4: [0,0,0,0,1]
   M5[4] = [5]int64{0, 0, 0, 0, 1}
   m5 := powMat5(M5, n-3)
   // initial T3 = [p3,p2,p1,S3,1] = [0,0,0,3,1]
   // compute Tn = m5 * T3
   var T3 = [5]int64{0, 0, 0, 3, 1}
   var Tn [5]int64
   for i := 0; i < 5; i++ {
       var sum int64
       for j := 0; j < 5; j++ {
           sum += m5[i][j] * T3[j]
       }
       Tn[i] = sum % MOD1
   }
   p := Tn[0]
   // compute final answer
   res := modPow(f1, e1, MOD)
   res = res * modPow(f2, e2, MOD) % MOD
   res = res * modPow(f3, e3, MOD) % MOD
   res = res * modPow(c, p, MOD) % MOD
   fmt.Println(res)
}

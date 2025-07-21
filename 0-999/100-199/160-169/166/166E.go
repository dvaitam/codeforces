package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

// 2x2 matrix
type mat [2][2]int64

// multiply two 2x2 matrices modulo mod
func mul(a, b mat) mat {
   return mat{
       {(a[0][0]*b[0][0] + a[0][1]*b[1][0]) % mod, (a[0][0]*b[0][1] + a[0][1]*b[1][1]) % mod},
       {(a[1][0]*b[0][0] + a[1][1]*b[1][0]) % mod, (a[1][0]*b[0][1] + a[1][1]*b[1][1]) % mod},
   }
}

// fast exponentiation of matrix m to power e
func powMat(m mat, e int64) mat {
   // initialize result as identity
   res := mat{{1, 0}, {0, 1}}
   for e > 0 {
       if e&1 == 1 {
           res = mul(res, m)
       }
       m = mul(m, m)
       e >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int64
   fmt.Fscan(reader, &n)
   // transition matrix: [ [0,3], [1,2] ]
   m := mat{{0, 3}, {1, 2}}
   // compute m^n
   mn := powMat(m, n)
   // starting vector [1,0]^T, answer is first element: mn[0][0]
   ans := mn[0][0] % mod
   fmt.Println(ans)
}

package main

import (
   "fmt"
   "os"
)

// Matrix is a 2x2 matrix for Fibonacci computation
type Matrix [2][2]uint64

// mul returns a * b mod mod
func mul(a, b Matrix, mod uint64) Matrix {
   return Matrix{
       {(a[0][0]*b[0][0] + a[0][1]*b[1][0]) % mod, (a[0][0]*b[0][1] + a[0][1]*b[1][1]) % mod},
       {(a[1][0]*b[0][0] + a[1][1]*b[1][0]) % mod, (a[1][0]*b[0][1] + a[1][1]*b[1][1]) % mod},
   }
}

// matrixPow raises mat to the exp-th power modulo mod
func matrixPow(mat Matrix, exp, mod uint64) Matrix {
   // identity matrix
   res := Matrix{{1, 0}, {0, 1}}
   for exp > 0 {
       if exp&1 == 1 {
           res = mul(res, mat, mod)
       }
       mat = mul(mat, mat, mod)
       exp >>= 1
   }
   return res
}

// fib returns Fib(n) mod mod, where Fib(0)=0, Fib(1)=1
func fib(n, mod uint64) uint64 {
   if n == 0 {
       return 0
   }
   // matrix for Fibonacci
   M := Matrix{{1, 1}, {1, 0}}
   // M^(n-1)
   P := matrixPow(M, n-1, mod)
   return P[0][0]
}

// powMod returns base^exp mod mod
func powMod(base, exp, mod uint64) uint64 {
   result := uint64(1) % mod
   base %= mod
   for exp > 0 {
       if exp&1 == 1 {
           result = (result * base) % mod
       }
       base = (base * base) % mod
       exp >>= 1
   }
   return result
}

func main() {
   var n, k, l, m uint64
   if _, err := fmt.Fscan(os.Stdin, &n, &k, &l, &m); err != nil {
       return
   }
   // if k has bits outside [0, l), impossible
   if l < 64 && k>>l != 0 {
       fmt.Println(0)
       return
   }
   // count bits set in k
   c1 := uint64(0)
   kk := k
   for kk > 0 {
       c1 += kk & 1
       kk >>= 1
   }
   c0 := l - c1

   // number of sequences with no adjacent ones: A = Fib(n+2)
   A := fib(n+2, m)
   // total sequences of bits: T = 2^n
   T := powMod(2, n, m)
   // sequences with at least one adjacent ones: B = T - A
   var B uint64
   if T >= A {
       B = T - A
   } else {
       B = T + m - A
   }

   // result is A^c0 * B^c1 mod m
   r0 := powMod(A, c0, m)
   r1 := powMod(B, c1, m)
   res := (r0 * r1) % m
   fmt.Println(res)
}

package main

import (
   "bufio"
   "fmt"
   "os"
   "math/bits"
)

// Fast Walsh-Hadamard Transform for XOR convolution (in-place, unnormalized)
func fwht(a []int64, mod int64) {
   n := len(a)
   for lenStep := 1; lenStep < n; lenStep <<= 1 {
       for i := 0; i < n; i += lenStep << 1 {
           for j := 0; j < lenStep; j++ {
               u := a[i+j]
               v := a[i+j+lenStep]
               a[i+j] = u + v
               if a[i+j] >= mod {
                   a[i+j] -= mod
               }
               a[i+j+lenStep] = u - v
               if a[i+j+lenStep] < 0 {
                   a[i+j+lenStep] += mod
               }
           }
       }
   }
}

// modular exponentiation
func modPow(a, e, mod int64) int64 {
   res := int64(1)
   a %= mod
   for e > 0 {
       if e&1 == 1 {
           res = (res * a) % mod
       }
       a = (a * a) % mod
       e >>= 1
   }
   return res
}

// extended gcd for modular inverse
func extGCD(a, b int64) (int64, int64, int64) {
   if b == 0 {
       return a, 1, 0
   }
   g, x1, y1 := extGCD(b, a%b)
   x := y1
   y := x1 - (a/b)*y1
   return g, x, y
}

func modInv(a, mod int64) int64 {
   g, x, _ := extGCD(a, mod)
   if g != 1 {
       return -1 // inverse doesn't exist
   }
   x %= mod
   if x < 0 {
       x += mod
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var m int
   var t, p int64
   fmt.Fscan(reader, &m, &t, &p)
   n := 1 << m
   e0 := make([]int64, n)
   for i := 0; i < n; i++ {
       var x int64
       fmt.Fscan(reader, &x)
       e0[i] = x % p
   }
   b := make([]int64, m+1)
   for i := 0; i <= m; i++ {
       var x int64
       fmt.Fscan(reader, &x)
       b[i] = x % p
   }
   // build kernel array B of size n
   B := make([]int64, n)
   for i := 0; i < n; i++ {
       B[i] = b[bits.OnesCount(uint(i))]
   }
   // transforms
   fwht(e0, p)
   fwht(B, p)
   // exponentiate in transform domain
   for i := 0; i < n; i++ {
       e0[i] = modPow(B[i], t, p) * e0[i] % p
   }
   // inverse transform
   fwht(e0, p)
   // divide by n
   invN := modInv(int64(n), p)
   for i := 0; i < n; i++ {
       res := e0[i] * invN % p
       if res < 0 {
           res += p
       }
       fmt.Fprintln(writer, res)
   }
}

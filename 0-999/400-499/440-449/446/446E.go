package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

const MOD = 1051131

func modAdd(a, b int) int {
   c := a + b
   if c >= MOD {
       c -= MOD
   }
   return c
}
func modSub(a, b int) int {
   c := a - b
   if c < 0 {
       c += MOD
   }
   return c
}
func modMul(a, b int) int {
   return int((int64(a) * int64(b)) % MOD)
}
func modPow(a int, e int64) int {
   res := 1
   base := a % MOD
   for e > 0 {
       if e&1 != 0 {
           res = modMul(res, base)
       }
       base = modMul(base, base)
       e >>= 1
   }
   return res
}
// extended gcd
func egcd(a, b int) (int, int, int) {
   if b == 0 {
       return a, 1, 0
   }
   g, x1, y1 := egcd(b, a%b)
   return g, y1, x1 - (a/b)*y1
}
func modInv(a int) int {
   _, x, _ := egcd(a, MOD)
   x %= MOD
   if x < 0 {
       x += MOD
   }
   return x
}

// Fast Walsh-Hadamard Transform for XOR convolution
func fwht(a []int, invert bool) {
   n := len(a)
   for step := 1; step < n; step <<= 1 {
       for i := 0; i < n; i += step << 1 {
           for j := 0; j < step; j++ {
               u := a[i+j]
               v := a[i+j+step]
               a[i+j] = modAdd(u, v)
               a[i+j+step] = modSub(u, v)
           }
       }
   }
   if invert {
       invN := modInv(n)
       for i := 0; i < n; i++ {
           a[i] = modMul(a[i], invN)
       }
   }
}

func main() {
   // read input
   scanner := bufio.NewScanner(os.Stdin)
   scanner.Split(bufio.ScanWords)
   scanner.Scan(); m, _ := strconv.Atoi(scanner.Text())
   scanner.Scan(); t, _ := strconv.ParseInt(scanner.Text(), 10, 64)
   scanner.Scan(); s, _ := strconv.Atoi(scanner.Text())
   // initial array
   N := 1 << m
   a0 := make([]int, N)
   for i := 0; i < s; i++ {
       scanner.Scan()
       v, _ := strconv.Atoi(scanner.Text())
       a0[i] = v % MOD
   }
   // extend a
   for i := s; i < N; i++ {
       a0[i] = (int((101*int64(a0[i-s]) + 10007) % MOD))
   }
   // transform a0
   fwht(a0, false)
   // apply kernel powers
   pow2m := N % MOD
   for i := 0; i < N; i++ {
       // H[i] = (2^m +1 - 2*i) mod MOD
       hi := pow2m + 1 - ((i * 2) % MOD)
       hi %= MOD
       if hi < 0 {
           hi += MOD
       }
       a0[i] = modMul(a0[i], modPow(hi, t))
   }
   // inverse transform
   fwht(a0, true)
   // xor-sum of answers
   var res int
   for i := 0; i < N; i++ {
       res ^= a0[i]
   }
   fmt.Println(res)
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

const SZ = 128

// multiply convolves two distributions under XOR of size SZ
func multiply(a, b []float64) []float64 {
   t := make([]float64, SZ)
   for i := 0; i < SZ; i++ {
       ai := a[i]
       if ai == 0 {
           continue
       }
       for j := 0; j < SZ; j++ {
           t[i^j] += ai * b[j]
       }
   }
   return t
}

// power raises base distribution to p-th XOR-convolution power
func power(base []float64, p int) []float64 {
   // identity: delta at 0
   res := make([]float64, SZ)
   res[0] = 1.0
   b := make([]float64, SZ)
   copy(b, base)
   for p > 0 {
       if p&1 == 1 {
           res = multiply(res, b)
       }
       b = multiply(b, b)
       p >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   // read distribution z of values 0..m
   z := make([]float64, SZ)
   for i := 0; i <= m && i < SZ; i++ {
       var v float64
       fmt.Fscan(in, &v)
       z[i] = v
   }
   // compute z^n under XOR convolution
   p := power(z, n)
   // probability of non-zero XOR sum = 1 - p[0]
   res := 1.0 - p[0]
   // print with 7 decimal places
   fmt.Printf("%.7f", res)
}

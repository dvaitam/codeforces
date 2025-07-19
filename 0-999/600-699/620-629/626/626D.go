package main

import (
   "fmt"
   "math"
)

// fft performs in-place Fast Fourier Transform on a, invert indicates inverse transform
func fft(a []complex128, invert bool) {
   n := len(a)
   // bit-reversal permutation
   for i, j := 1, 0; i < n; i++ {
       bit := n >> 1
       for j&bit != 0 {
           j ^= bit
           bit >>= 1
       }
       j |= bit
       if i < j {
           a[i], a[j] = a[j], a[i]
       }
   }
   // Cooley-Tukey
   for length := 2; length <= n; length <<= 1 {
       ang := 2 * math.Pi / float64(length)
       if !invert {
           ang = -ang
       }
       wlen := complex(math.Cos(ang), math.Sin(ang))
       for i := 0; i < n; i += length {
           var w complex128 = 1
           half := length >> 1
           for j := 0; j < half; j++ {
               u := a[i+j]
               v := a[i+j+half] * w
               a[i+j] = u + v
               a[i+j+half] = u - v
               w *= wlen
           }
       }
   }
   if invert {
       invN := complex(1/float64(n), 0)
       for i := range a {
           a[i] *= invN
       }
   }
}

// gcd returns the greatest common divisor of x and y
func gcd(x, y int64) int64 {
   for y != 0 {
       x, y = y, x%y
   }
   return x
}

func main() {
   var n int
   fmt.Scan(&n)
   in := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&in[i])
   }
   // count pairwise absolute differences
   const maxD = 10000
   tmp := make([]int64, maxD+2)
   for i := 0; i < n; i++ {
       for j := i + 1; j < n; j++ {
           d := in[i] - in[j]
           if d < 0 {
               d = -d
           }
           tmp[d]++
       }
   }
   // prepare FFT of size power of two > 2*maxD
   fftN := 1
   for fftN <= 2*maxD {
       fftN <<= 1
   }
   a := make([]complex128, fftN)
   for i := 0; i <= maxD; i++ {
       a[i] = complex(float64(tmp[i]), 0)
   }
   // convolution via FFT
   fft(a, false)
   for i := 0; i < fftN; i++ {
       a[i] *= a[i]
   }
   fft(a, true)
   // extract convolution results
   ans := make([]int64, maxD+1)
   for i := 0; i <= maxD; i++ {
       val := real(a[i])
       ans[i] = int64(val + 0.5)
   }
   // suffix sums of tmp
   for i := maxD; i >= 1; i-- {
       tmp[i] += tmp[i+1]
   }
   // count triples
   var x int64
   for i := 1; i < maxD; i++ {
       x += ans[i] * tmp[i+1]
   }
   // total number of pairs
   y0 := int64(n) * int64(n-1) / 2
   y := y0 * y0 * y0
   g := gcd(x, y)
   x /= g
   y /= g
   // print ratio with 15 decimal places
   ratio := float64(x) / float64(y)
   fmt.Printf("%.15f\n", ratio)
}

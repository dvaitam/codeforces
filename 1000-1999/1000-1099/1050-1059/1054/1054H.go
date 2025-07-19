package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

const (
   mod = 490019
   Bc  = 32768
   hf  = 0.5
)

// qpow computes x^y mod mod
func qpow(x int, y int64) int {
   res := int64(1)
   base := int64(x)
   for y > 0 {
       if y&1 == 1 {
           res = res * base % mod
       }
       base = base * base % mod
       y >>= 1
   }
   return int(res)
}

// getLim returns smallest power of two >= x
func getLim(x int) int {
   lim := 1
   for lim < x {
       lim <<= 1
   }
   return lim
}

// getRev fills bit-reversal permutation for length lim
func getRev(rev []int, lim int) {
   for i := 0; i < lim; i++ {
       rev[i] = rev[i>>1] >> 1
       if i&1 != 0 {
           rev[i] |= lim >> 1
       }
   }
}

// DFT performs FFT on arr with direction typ (1 forward, -1 inverse)
func DFT(arr []complex128, typ int, lim int, rev []int, w []complex128) {
   // bit reversal
   for i := 0; i < lim; i++ {
       j := rev[i]
       if i < j {
           arr[i], arr[j] = arr[j], arr[i]
       }
   }
   // Cooley-Tukey
   for length := 1; length < lim; length <<= 1 {
       ang := math.Pi / float64(length)
       t := complex(math.Cos(ang), math.Sin(ang)*float64(typ))
       // precompute twiddles
       for i := length - 2; i >= 0; i -= 2 {
           w[i] = w[i>>1] * t
           w[i+1] = w[i] * t
       }
       // butterfly
       step := length << 1
       for k := 0; k < lim; k += step {
           for p := 0; p < length; p++ {
               u := arr[k+p]
               v := w[p] * arr[k+length+p]
               arr[k+p] = u + v
               arr[k+length+p] = u - v
           }
       }
   }
}

// tFFT performs convolution on A and B of original length origLen
func tFFT(A []int, B []int, origLen, lim int, rev []int) {
   fa0 := make([]complex128, lim)
   fa1 := make([]complex128, lim)
   fb := make([]complex128, lim)
   w := make([]complex128, lim)
   // pack integer arrays into complex for FFT
   for i := 0; i < lim; i++ {
       var x, y, z, wv int
       if i < origLen {
           x = A[i] >> 15
           y = A[i] & (Bc - 1)
           z = B[i] >> 15
           wv = B[i] & (Bc - 1)
       }
       fa0[i] = complex(float64(x), float64(y))
       fa1[i] = complex(float64(x), -float64(y))
       fb[i] = complex(float64(z), float64(wv))
   }
   DFT(fa0, 1, lim, rev, w)
   DFT(fa1, 1, lim, rev, w)
   DFT(fb, 1, lim, rev, w)
   // multiply and scale
   for i := 0; i < lim; i++ {
       fb[i] /= complex(float64(lim), 0)
       fa0[i] *= fb[i]
       fa1[i] *= fb[i]
   }
   // inverse FFT
   DFT(fa0, -1, lim, rev, w)
   DFT(fa1, -1, lim, rev, w)
   // unpack results
   for i := 0; i < origLen; i++ {
       xr := math.Floor((real(fa0[i])+real(fa1[i]))*hf + hf)
       yr := math.Floor((imag(fa0[i])+imag(fa1[i]))*hf + hf)
       zr := math.Floor(imag(fa0[i]) + hf)
       wr := math.Floor(real(fa1[i]) + hf)
       x := int64(int(xr)) % mod
       y := int64(int(yr)) % mod
       z := (int64(int(zr)) - y) % mod
       wv2 := (int64(int(wr)) - x) % mod
       A[i] = int(((x*Bc%mod*Bc%mod + ((y+z)%mod)*Bc%mod + wv2) % mod + mod) % mod)
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m, c int
   fmt.Fscan(in, &n, &m, &c)
   origLen := mod << 1
   lim := getLim(origLen)
   rev := make([]int, lim)
   getRev(rev, lim)
   A := make([]int, lim)
   B := make([]int, lim)
   // read and encode A
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(in, &x)
       p := int((int64(i) * int64(i)) % int64(mod-1))
       exp := (int64(p)*(int64(p)-1)>>1) % int64(mod-1)
       t := qpow(c, exp)
       invt := qpow(t, int64(mod-2))
       A[p] = (A[p] + int(int64(x)*int64(invt)%mod)) % mod
   }
   // read and encode B
   for i := 0; i < m; i++ {
       var x int
       fmt.Fscan(in, &x)
       p := int((int64(i) * int64(i) * int64(i)) % int64(mod-1))
       exp := (int64(p)*(int64(p)-1)>>1) % int64(mod-1)
       t := qpow(c, exp)
       invt := qpow(t, int64(mod-2))
       B[p] = (B[p] + int(int64(x)*int64(invt)%mod)) % mod
   }
   // convolution
   tFFT(A, B, origLen, lim, rev)
   // compute answer
   var ans int64
   for i := 0; i < origLen; i++ {
       exp := (int64(i)*(int64(i)-1) >> 1) % int64(mod-1)
       pw := qpow(c, exp)
       ans = (ans + int64(A[i])*int64(pw) % mod) % mod
   }
   fmt.Fprint(out, ans)
}

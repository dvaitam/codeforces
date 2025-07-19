package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func powmod(a, x, p int) int {
   res := 1
   a %= p
   for x > 0 {
       if x&1 != 0 {
           res = int(int64(res) * int64(a) % int64(p))
       }
       a = int(int64(a) * int64(a) % int64(p))
       x >>= 1
   }
   return res
}

func isPrime(x int) bool {
   if x < 2 {
       return false
   }
   for i := 2; i*i <= x; i++ {
       if x%i == 0 {
           return false
       }
   }
   return true
}

func factor(n int) []int {
   pr := []int{}
   x := n
   for i := 2; i*i <= x; i++ {
       if x%i == 0 {
           pr = append(pr, i)
           for x%i == 0 {
               x /= i
           }
       }
   }
   if x > 1 {
       pr = append(pr, x)
   }
   return pr
}

func check(g, p int, pr []int) bool {
   for _, t := range pr {
       if powmod(g, (p-1)/t, p) == 1 {
           return false
       }
   }
   return true
}

func choose2(x, k int) int {
   return (x*(x-1)/2) % k
}

// FFT implementation
func fft(a []complex128, invert bool) {
   n := len(a)
   j := 0
   for i := 1; i < n; i++ {
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
   for length := 2; length <= n; length <<= 1 {
       angle := 2 * math.Pi / float64(length)
       if invert {
           angle = -angle
       }
       wlen := complex(math.Cos(angle), math.Sin(angle))
       for i := 0; i < n; i += length {
           w := complex(1, 0)
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

func multiply(A, B []int, p int) []int {
   n := 1
   for n < len(A)+len(B)-1 {
       n <<= 1
   }
   fa := make([]complex128, n)
   fb := make([]complex128, n)
   for i := 0; i < len(A); i++ {
       fa[i] = complex(float64(A[i]&32767), float64(A[i]>>15))
   }
   for i := len(A); i < n; i++ {
       fa[i] = 0
   }
   for i := 0; i < len(B); i++ {
       fb[i] = complex(float64(B[i]&32767), float64(B[i]>>15))
   }
   for i := len(B); i < n; i++ {
       fb[i] = 0
   }
   fft(fa, false)
   fft(fb, false)
   dfta := make([]complex128, n)
   dftb := make([]complex128, n)
   for i := 0; i < n; i++ {
       j := (n - i) & (n - 1)
       ai := fa[i]
       aj := fa[j]
       bj := fb[j]
       conjAi := complex(real(ai), -imag(ai))
       dfta[i] = (conjAi + aj) * bj * complex(0.5, 0)
       dftb[i] = (conjAi - aj) * bj * complex(0, 0.5)
   }
   fft(dfta, true)
   fft(dftb, true)
   result := make([]int, n)
   for i := 0; i < n; i++ {
       ret0 := int(math.Round(real(dfta[i])))
       ret1 := int(math.Round(imag(dfta[i]) + real(dftb[i])))
       ret2 := int(math.Round(imag(dftb[i])))
       val := (ret0 + ((ret1%p)<<15) + ((ret2%p)<<30)) % p
       if val < 0 {
           val += p
       }
       result[i] = val
   }
   return result
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n, &m)
   var s1, s2 string
   fmt.Fscan(reader, &s1, &s2)
   n = len(s1)
   if m < n+2 {
       m = n + 2
   }
   if m < 5 {
       m = 5
   }
   for p := m; ; p++ {
       if !isPrime(p) {
           continue
       }
       k := p - 1
       pr := factor(k)
       var g int
       for g = 2; ; g++ {
           if check(g, p, pr) {
               break
           }
       }
       omega := make([]int, k+1)
       romega := make([]int, k+1)
       omega[0] = 1
       omega[1] = powmod(g, (p-1)/k, p)
       for i := 2; i <= k; i++ {
           omega[i] = omega[i-1] * omega[1] % p
       }
       for i := 0; i <= k; i++ {
           romega[i] = omega[(k-i)% (k)]
       }
       // prepare A and C
       A := make([]int, 2*k)
       for i := 0; i <= 2*k-2; i++ {
           A[i] = omega[choose2(i, k)]
       }
       C := make([]int, k+1)
       for j := 0; j <= k; j++ {
           var cj int
           if j < n {
               cj = int(s1[j] - s2[j])
           }
           C[j] = (cj%p + p) % p
       }
       for j := 0; j <= k; j++ {
           C[j] = romega[choose2(j, k)] * C[j] % p
       }
       // reverse C
       for i, j := 0, k; i < j; i, j = i+1, j-1 {
           C[i], C[j] = C[j], C[i]
       }
       // convolution
       lenConv := 1
       for lenConv < len(A)+len(C)-1 {
           lenConv <<= 1
       }
       A2 := make([]int, lenConv)
       for i := range A {
           A2[i] = A[i]
       }
       C2 := make([]int, lenConv)
       for i := range C {
           C2[i] = C[i]
       }
       res := multiply(A2, C2, p)
       for t := 0; t < k; t++ {
           ret := res[k+t] * romega[choose2(t, k)] % p
           if ret == 0 && omega[t] >= 2 && omega[t] <= p-2 {
               fmt.Println(p, omega[t])
               return
           }
       }
   }
}

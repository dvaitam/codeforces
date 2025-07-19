package main

import (
   "bufio"
   "fmt"
   "os"
)

func p1(x uint32) uint32 {
   y := x + 1
   if x%2 == 0 {
       x /= 2
   } else {
       y /= 2
   }
   return x * y
}

func p2(x uint32) uint32 {
   y := x + 1
   z := x*2 + 1
   if x%2 == 0 {
       x /= 2
   } else {
       y /= 2
   }
   if z%3 == 0 {
       z /= 3
   } else if x%3 == 0 {
       x /= 3
   } else {
       y /= 3
   }
   return x * y * z
}

func p3(x uint32) uint32 {
   y := p1(x)
   return y * y
}

// solve computes the main sum part using Meissel-Lehmer style
func solve(n, A, B, C, D uint32) uint32 {
   // compute m = floor(sqrt(n)) + 1
   var m uint32 = 1
   for m*m <= n {
       m++
   }
   // allocate f and g arrays: index 0 unused, need up to m+1
   size := int(m + 2)
   f0 := make([]uint32, size)
   f1 := make([]uint32, size)
   f2 := make([]uint32, size)
   f3 := make([]uint32, size)
   g0 := make([]uint32, size)
   g1 := make([]uint32, size)
   g2 := make([]uint32, size)
   g3 := make([]uint32, size)
   // initial values
   for j := uint32(1); j*j <= n; j++ {
       idx := int(j)
       nj := n / j
       f0[idx] = nj - 1
       f1[idx] = (p1(nj) - 1) * C
       f2[idx] = (p2(nj) - 1) * B
       f3[idx] = (p3(nj) - 1) * A
   }
   // g initialization
   for i := uint32(1); i <= m; i++ {
       idx := int(i)
       g0[idx] = i - 1
       g1[idx] = (p1(i) - 1) * C
       g2[idx] = (p2(i) - 1) * B
       g3[idx] = (p3(i) - 1) * A
   }
   // inclusion-exclusion
   for i := uint32(2); i <= m; i++ {
       if g0[int(i)] == g0[int(i-1)] {
           continue
       }
       var o0, o1, o2, o3 uint32 = 1, i, i * i, i * i * i
       lim := n / (i * i)
       if lim > m-1 {
           lim = m - 1
       }
       for j := uint32(1); j <= lim; j++ {
           ji := int(j)
           if i*j < m {
               f0[ji] -= o0 * (f0[int(i*j)] - g0[int(i-1)])
               f1[ji] -= o1 * (f1[int(i*j)] - g1[int(i-1)])
               f2[ji] -= o2 * (f2[int(i*j)] - g2[int(i-1)])
               f3[ji] -= o3 * (f3[int(i*j)] - g3[int(i-1)])
           } else {
               nij := n / i / j
               f0[ji] -= o0 * (g0[int(nij)] - g0[int(i-1)])
               f1[ji] -= o1 * (g1[int(nij)] - g1[int(i-1)])
               f2[ji] -= o2 * (g2[int(nij)] - g2[int(i-1)])
               f3[ji] -= o3 * (g3[int(nij)] - g3[int(i-1)])
           }
       }
       for j := m; j >= i*i; j-- {
           ji := int(j)
           g0[ji] -= o0 * (g0[int(j/i)] - g0[int(i-1)])
           g1[ji] -= o1 * (g1[int(j/i)] - g1[int(i-1)])
           g2[ji] -= o2 * (g2[int(j/i)] - g2[int(i-1)])
           g3[ji] -= o3 * (g3[int(j/i)] - g3[int(i-1)])
       }
   }
   // apply D coefficient to constant term
   for i := uint32(1); i <= m+1; i++ {
       idx := int(i)
       f0[idx] *= D
       g0[idx] *= D
   }
   // accumulate result
   var rt uint32
   for i := uint32(1); n/i > m; i++ {
       idx := int(i)
       mid := int(m)
       rt += f0[idx] - g0[mid]
       rt += f1[idx] - g1[mid]
       rt += f2[idx] - g2[mid]
       rt += f3[idx] - g3[mid]
   }
   return rt
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, A, B, C, D uint32
   fmt.Fscan(reader, &n, &A, &B, &C, &D)
   ans := solve(n, A, B, C, D)
   // sieve primes up to m = floor(sqrt(n)) + 1 as in solve
   var m uint32 = 1
   for m*m <= n {
       m++
   }
   sizeP := int(m) + 1
   isPrime := make([]bool, sizeP)
   for i := range isPrime {
       isPrime[i] = true
   }
   if sizeP > 0 {
       isPrime[0] = false
   }
   if sizeP > 1 {
       isPrime[1] = false
   }
   for i := uint32(2); i*i <= m; i++ {
       if isPrime[int(i)] {
           for j := i * i; j <= m; j += i {
               isPrime[int(j)] = false
           }
       }
   }
   // add prime power contributions
   for i := uint32(2); i <= m; i++ {
       if !isPrime[int(i)] {
           continue
       }
       res := A*i*i*i + B*i*i + C*i + D
       tmp := n
       for tmp > 0 {
           tmp /= i
           ans += res * tmp
       }
   }
   fmt.Println(ans)
}

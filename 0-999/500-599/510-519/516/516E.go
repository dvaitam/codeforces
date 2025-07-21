package main

import (
   "bufio"
   "fmt"
   "os"
)

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

// extended gcd: returns (g, x, y) s.t. a*x + b*y = g
func extGCD(a, b int64) (int64, int64, int64) {
   if b == 0 {
       return a, 1, 0
   }
   g, x1, y1 := extGCD(b, a%b)
   return g, y1, x1 - (a/b)*y1
}

// modinv computes inverse of a mod m, assumes gcd(a,m)=1
func modinv(a, m int64) int64 {
   g, x, _ := extGCD(a, m)
   if g != 1 {
       return 0 // no inv
   }
   x %= m
   if x < 0 {
       x += m
   }
   return x
}

// compute maximum minimal meet time from sources A (positions) on cycle of length n
// to all targets on cycle m, where steps advance by n/g each time
func compute(A []int64, n, m, d int64) int64 {
   // n: primary cycle length, m: target cycle length, d: gcd(n_orig, m_orig)
   N := n / d
   M := m / d
   invN := modinv(N%M, M)
   const INF64 = 1<<63 - 1
   Cmin := INF64
   for _, xi := range A {
       // xi is position in [0..n)
       x0 := (xi) / d
       s := (x0 * invN) % M
       // Qi = xi - s * n
       Qi := xi - s*n
       if Qi < Cmin {
           Cmin = Qi
       }
   }
   // max time over targets = Cmin + (M-1)*n
   return Cmin + (M-1)*n
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int64
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   var b int
   fmt.Fscan(reader, &b)
   boys := make([]int64, b)
   for i := 0; i < b; i++ {
       fmt.Fscan(reader, &boys[i])
   }
   var g int
   fmt.Fscan(reader, &g)
   girls := make([]int64, g)
   for i := 0; i < g; i++ {
       fmt.Fscan(reader, &girls[i])
   }
   if b == 0 || g == 0 {
       fmt.Println(-1)
       return
   }
   d := gcd(n, m)
   // check connectivity: for each residue mod d, there must be at least one source
   haveB := make(map[int64]bool)
   for _, x := range boys {
       haveB[x%d] = true
   }
   haveG := make(map[int64]bool)
   for _, y := range girls {
       haveG[y%d] = true
   }
   for r := int64(0); r < d; r++ {
       if !haveB[r] && !haveG[r] {
           fmt.Println(-1)
           return
       }
   }
   // compute waves
   t1 := compute(boys, n, m, d)
   t2 := compute(girls, m, n, d)
   ans := t1
   if t2 > ans {
       ans = t2
   }
   fmt.Println(ans)
}

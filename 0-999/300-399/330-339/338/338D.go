package main

import (
   "bufio"
   "fmt"
   "os"
)

func extGCD(a, b int64) (int64, int64, int64) {
   if b == 0 {
       return a, 1, 0
   }
   g, x1, y1 := extGCD(b, a%b)
   return g, y1, x1 - (a/b)*y1
}

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int64
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   a := make([]int64, k)
   for i := int64(0); i < k; i++ {
       fmt.Fscan(in, &a[i])
   }
   // impossible if sequence longer than columns
   if m < k {
       fmt.Println("NO")
       return
   }
   // impossible if any a[i] > n
   for _, v := range a {
       if v > n {
           fmt.Println("NO")
           return
       }
   }
   var r int64 = 0
   var M int64 = 1
   for idx, ai := range a {
       // j + idx ≡ 0 mod ai  => j ≡ -idx mod ai
       b := ((-int64(idx))%ai + ai) % ai
       // merge j ≡ r mod M and j ≡ b mod ai
       g := gcd(M, ai)
       if (b-r)%g != 0 {
           fmt.Println("NO")
           return
       }
       // compute modular inverse of M/g mod ai/g
       mg := M / g
       ag := ai / g
       _, x, _ := extGCD(mg, ag)
       inv := (x%ag + ag) % ag
       // delta = (b - r) / g * inv mod (ai/g)
       delta := ((b - r) / g % ag + ag) % ag
       delta = (delta * inv) % ag
       // update r
       r = r + M*delta
       // update modulus
       // check overflow relative to n
       if mg > 0 && (n/ mg) < ai {
           fmt.Println("NO")
           return
       }
       M = M * ag
       // normalize r
       r %= M
       if r < 0 {
           r += M
       }
   }
   // minimal positive j0
   j0 := r
   if j0 == 0 {
       j0 = M
   }
   // check fits in width
   if j0+k-1 > m {
       fmt.Println("NO")
       return
   }
   // use row i = M
   L := M
   if L > n {
       fmt.Println("NO")
       return
   }
   // verify gcds
   for idx, ai := range a {
       val := gcd(L, j0+int64(idx))
       if val != ai {
           fmt.Println("NO")
           return
       }
   }
   fmt.Println("YES")
}

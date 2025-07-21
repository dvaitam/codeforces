package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var w, h int
   if _, err := fmt.Fscan(reader, &w, &h); err != nil {
       return
   }
   // total points
   N := int64(w+1) * int64(h+1) % mod
   // total ordered triples
   T := N * (N - 1) % mod * (N - 2) % mod
   // count triples with odd cross (half-integer area)
   // parity classes counts
   nx0 := w/2 + 1
   nx1 := (w + 1) - nx0
   ny0 := h/2 + 1
   ny1 := (h + 1) - ny0
   c00 := int64(nx0) * int64(ny0) % mod
   c01 := int64(nx0) * int64(ny1) % mod
   c10 := int64(nx1) * int64(ny0) % mod
   c11 := int64(nx1) * int64(ny1) % mod
   // number of ordered triples with cross mod2 == 1
   // two non-degenerate F2 triangles: (00,01,10) and (11,01,10), each 6 permutations
   odd := int64(6) * ( (c00*c01%mod)*c10%mod + (c11*c01%mod)*c10%mod ) % mod
   E := (T - odd + mod) % mod
   // subtract degenerate (collinear) triples
   var sum int64
   W1 := w + 1
   H1 := h + 1
   for dx := 0; dx <= w; dx++ {
       for dy := 0; dy <= h; dy++ {
           if dx == 0 && dy == 0 {
               continue
           }
           g := gcd(dx, dy)
           if g <= 1 {
               continue
           }
           u := int64(W1-dx) * int64(H1-dy) % mod
           sum = (sum + u*int64(g-1)) % mod
       }
   }
   // unordered triples multiply by 6 for ordered
   D := sum * 6 % mod
   ans := (E - D + mod) % mod
   fmt.Println(ans)
}

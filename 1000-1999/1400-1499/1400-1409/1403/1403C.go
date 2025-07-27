package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

var fact, invFact []int64

// modPow computes a^b mod mod
func modPow(a, b int64) int64 {
   res := int64(1)
   a %= mod
   for b > 0 {
       if b&1 == 1 {
           res = res * a % mod
       }
       a = a * a % mod
       b >>= 1
   }
   return res
}

// initComb precomputes factorials up to n
func initComb(n int) {
   fact = make([]int64, n+1)
   invFact = make([]int64, n+1)
   fact[0] = 1
   for i := 1; i <= n; i++ {
       fact[i] = fact[i-1] * int64(i) % mod
   }
   invFact[n] = modPow(fact[n], mod-2)
   for i := n; i > 0; i-- {
       invFact[i-1] = invFact[i] * int64(i) % mod
   }
}

// comb computes C(n, k) mod mod, for n>=0, k>=0, k small
func comb(n int64, k int64) int64 {
   if k < 0 || k > n {
       return 0
   }
   num := int64(1)
   for i := int64(0); i < k; i++ {
       num = num * ((n - i) % mod) % mod
   }
   return num * invFact[k] % mod
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var R int64
   var C, Q int
   fmt.Fscan(reader, &R, &C, &Q)
   initComb(C + 5)
   for qi := 0; qi < Q; qi++ {
       var t rune
       var c1, cR int64
       fmt.Fscan(reader, &t, &c1, &cR)
       var moves, ways int64
       switch t {
       case 'P':
           if c1 == cR {
               moves = R - 1
               ways = 1
           }
       case 'R':
           if c1 == cR {
               moves = 1
               ways = 1
           } else {
               moves = 2
               ways = 2
           }
       case 'B':
           // bishop moves
           if (1+c1)%2 != (R+cR)%2 {
               moves, ways = 0, 0
           } else if abs64(cR-c1) == R-1 {
               moves, ways = 1, 1
           } else {
               moves = 2
               s1 := 1 + c1
               d1 := 1 - c1
               s2 := R + cR
               d2 := R - cR
               cnt := int64(0)
               // intersection of s1 and d2
               if (s1+d2)%2 == 0 {
                   r := (s1 + d2) / 2
                   c := s1 - r
                   if r >= 1 && r <= R && c >= 1 && c <= int64(C) {
                       cnt++
                   }
               }
               // intersection of d1 and s2
               if (s2+d1)%2 == 0 {
                   r := (s2 + d1) / 2
                   c := s2 - r
                   if r >= 1 && r <= R && c >= 1 && c <= int64(C) {
                       cnt++
                   }
               }
               ways = cnt
           }
       case 'Q':
           // queen moves
           if c1 == cR || abs64(cR-c1) == R-1 {
               moves, ways = 1, 1
           } else {
               moves = 2
               // define lines for start and end
               type line struct{ tp int; v int64 }
               s1 := 1 + c1; d1 := 1 - c1
               s2 := R + cR; d2 := R - cR
               startLines := []line{{0, 1}, {1, c1}, {2, s1}, {3, d1}}
               endLines := []line{{0, R}, {1, cR}, {2, s2}, {3, d2}}
               seen := make(map[int64]bool)
               cnt := int64(0)
               for _, ls := range startLines {
                   for _, le := range endLines {
                       if ls.tp == le.tp {
                           continue
                       }
                       var r, c int64
                       ok := true
                       switch ls.tp {
                       case 0:
                           // r = ls.v
                           r = ls.v
                           switch le.tp {
                           case 1:
                               c = le.v
                           case 2:
                               c = le.v - r
                           case 3:
                               c = r - le.v
                           }
                       case 1:
                           // c = ls.v
                           c = ls.v
                           switch le.tp {
                           case 0:
                               r = le.v
                           case 2:
                               r = le.v - c
                           case 3:
                               r = le.v + c
                           }
                       case 2:
                           // r+c = ls.v
                           switch le.tp {
                           case 0:
                               r = le.v
                               c = ls.v - r
                           case 1:
                               c = le.v
                               r = ls.v - c
                           case 3:
                               if (ls.v+le.v)%2 != 0 {
                                   ok = false
                               } else {
                                   r = (ls.v + le.v) / 2
                                   c = ls.v - r
                               }
                           }
                       case 3:
                           // r-c = ls.v
                           switch le.tp {
                           case 0:
                               r = le.v
                               c = r - ls.v
                           case 1:
                               c = le.v
                               r = ls.v + c
                           case 2:
                               if (ls.v+le.v)%2 != 0 {
                                   ok = false
                               } else {
                                   r = (ls.v + le.v) / 2
                                   c = le.v - r
                               }
                           }
                       }
                       if !ok {
                           continue
                       }
                       if r < 1 || r > R || c < 1 || c > int64(C) {
                           continue
                       }
                       key := r*int64(C+1) + c
                       if !seen[key] {
                           seen[key] = true
                           cnt++
                       }
                   }
               }
               ways = cnt
           }
       case 'K':
           // king moves: Chebyshev distance always R-1
           dc := abs64(cR - c1)
           moves = R - 1
           // number of diagonal moves = dc, vertical moves = moves - dc
           ways = comb(moves, dc)
       }
       fmt.Fprintln(writer, moves, ways)
   }
}

func abs64(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}

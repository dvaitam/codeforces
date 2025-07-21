package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 51123987

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func abs(a int) int {
   if a < 0 {
       return -a
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   var s string
   fmt.Fscan(reader, &n)
   fmt.Fscan(reader, &s)
   // letters: map a->0, b->1, c->2
   // preprocess next occurrence
   nextOcc := make([][3]int, n+2)
   // initialize beyond end
   for x := 0; x < 3; x++ {
       nextOcc[n][x] = n + 1
       nextOcc[n+1][x] = n + 1
   }
   // build from n-1 down to 0; we want occ >= i+1
   for i := n - 1; i >= 0; i-- {
       for x := 0; x < 3; x++ {
           nextOcc[i][x] = nextOcc[i+1][x]
       }
       var xi int
       switch s[i] {
       case 'a': xi = 0
       case 'b': xi = 1
       default: xi = 2
       }
       nextOcc[i][xi] = i + 1
   }
   // compute balanced distributions
   type triple struct{a, b, c int}
   var dists []triple
   base := n / 3
   rem := n % 3
   if rem == 0 {
       dists = append(dists, triple{base, base, base})
   } else if rem == 1 {
       dists = append(dists, triple{base+1, base, base})
       dists = append(dists, triple{base, base+1, base})
       dists = append(dists, triple{base, base, base+1})
   } else { // rem == 2
       dists = append(dists, triple{base+1, base+1, base})
       dists = append(dists, triple{base+1, base, base+1})
       dists = append(dists, triple{base, base+1, base+1})
   }
   ans := 0
   // dp[pos][ca][cb]
   for _, d := range dists {
       va, vb, vc := d.a, d.b, d.c
       // allocate dp
       dp := make([][][]int, n+1)
       for i := 0; i <= n; i++ {
           dp[i] = make([][]int, va+1)
           for ca := 0; ca <= va; ca++ {
               dp[i][ca] = make([]int, vb+1)
           }
       }
       dp[0][0][0] = 1
       for pos := 0; pos < n; pos++ {
           for ca := 0; ca <= va; ca++ {
               for cb := 0; cb <= vb; cb++ {
                   cur := dp[pos][ca][cb]
                   if cur == 0 {
                       continue
                   }
                   ccDone := pos - ca - cb
                   // try each letter x
                   for x := 0; x < 3; x++ {
                       // remaining cap for x
                       var remCap int
                       switch x {
                       case 0:
                           remCap = va - ca
                       case 1:
                           remCap = vb - cb
                       default:
                           remCap = vc - ccDone
                       }
                       if remCap <= 0 {
                           continue
                       }
                       nxt := nextOcc[pos][x]
                       if nxt > n {
                           continue
                       }
                       // block length L from Lmin to Lmax
                       Lmin := nxt - pos
                       // Lmax limited by remCap and remainder of positions
                       maxLen := n - pos
                       if remCap < maxLen {
                           maxLen = remCap
                       }
                       for L := Lmin; L <= maxLen; L++ {
                           pos2 := pos + L
                           ca2 := ca
                           cb2 := cb
                           if x == 0 {
                               ca2 += L
                           } else if x == 1 {
                               cb2 += L
                           }
                           dp[pos2][ca2][cb2] += cur
                           if dp[pos2][ca2][cb2] >= MOD {
                               dp[pos2][ca2][cb2] -= MOD
                           }
                       }
                   }
               }
           }
       }
       ans = (ans + dp[n][va][vb]) % MOD
   }
   fmt.Fprintln(writer, ans)
}

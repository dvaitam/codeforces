package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
)

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

// floor division a/b, floor for negative a
func floorDiv(a, b int64) int64 {
   if b < 0 {
       a, b = -a, -b
   }
   if a >= 0 {
       return a / b
   }
   return -(((-a) + b - 1) / b)
}

// ceil division a/b
func ceilDiv(a, b int64) int64 {
   if b < 0 {
       a, b = -a, -b
   }
   if a >= 0 {
       return (a + b - 1) / b
   }
   return -((-a) / b)
}

func abs(a int64) int64 {
   if a < 0 {
       return -a
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   line, _ := reader.ReadString('\n')
   t, _ := strconv.Atoi(strings.TrimSpace(line))
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for ti := 0; ti < t; ti++ {
       line, _ = reader.ReadString('\n')
       parts := strings.Fields(line)
       a, _ := strconv.ParseInt(parts[0], 10, 64)
       b, _ := strconv.ParseInt(parts[1], 10, 64)
       p, _ := strconv.ParseInt(parts[2], 10, 64)
       q, _ := strconv.ParseInt(parts[3], 10, 64)
       // search in [a,b]
       L := b - a + 1
       // if small interval, brute
       const bruteforceLimit = 200000
       var bestX int64 = a
       var bestDelta int64
       if L <= bruteforceLimit {
           // brute compute best
           bestDelta = -1
           for x := a; x <= b; x++ {
               rem := (2*p*x) % (2 * q)
               if rem < 0 {
                   rem += 2 * q
               }
               delta := abs(rem - q)
               if bestDelta < 0 || delta < bestDelta || (delta == bestDelta && x < bestX) {
                   bestDelta = delta
                   bestX = x
               }
           }
           fmt.Fprintln(out, bestX)
           continue
       }
       // large interval: consider global best classes and endpoints
       g := gcd(p, q)
       A2 := p / g
       M2 := q / g
       if M2 == 1 {
           // f(x) always zero
           fmt.Fprintln(out, a)
           continue
       }
       // modular inverse of A2 mod M2
       invA2 := func() int64 {
           // extended gcd
           aa, bb := A2, M2
           var x0, x1 int64 = 1, 0
           for bb != 0 {
               qd := aa / bb
               aa, bb = bb, aa-bb*qd
               x0, x1 = x1, x0-x1*qd
           }
           // aa is gcd
           if x0 < 0 {
               x0 += M2
           }
           return x0
       }()
       bestDelta = -1
       // candidate set
       cand := make(map[int64]struct{})
       // endpoints
       cand[a] = struct{}{}
       cand[b] = struct{}{}
       // best residue classes
       var ts []int64
       if M2%2 == 0 {
           ts = []int64{M2 / 2}
       } else {
           ts = []int64{M2/2 - 0, M2/2 + 1}
       }
       // consider base residues
       for _, t0 := range ts {
           for d := int64(0); d <= 2; d++ {
               for _, t := range []int64{t0 - d, t0 + d} {
                   if t < 0 || t >= M2 {
                       continue
                   }
                   x0 := (invA2 * t) % M2
                   // smallest >= a
                   var kmin int64
                   if x0 < a {
                       kmin = (a - x0 + M2 - 1) / M2
                   }
                   x1 := x0 + kmin*M2
                   if x1 >= a && x1 <= b {
                       cand[x1] = struct{}{}
                   }
                   // largest <= b
                   var kmax int64
                   if x0 <= b {
                       kmax = (b - x0) / M2
                       x2 := x0 + kmax*M2
                       if x2 >= a && x2 <= b {
                           cand[x2] = struct{}{}
                       }
                   }
               }
           }
       }
       // evaluate candidates
       for x := range cand {
           rem := (2*p*x) % (2 * q)
           if rem < 0 {
               rem += 2 * q
           }
           delta := abs(rem - q)
           if bestDelta < 0 || delta < bestDelta || (delta == bestDelta && x < bestX) {
               bestDelta = delta
               bestX = x
           }
       }
       fmt.Fprintln(out, bestX)
   }
}

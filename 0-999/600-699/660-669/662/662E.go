package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   times := make([][3]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &times[i][0], &times[i][1], &times[i][2])
   }
   // initial solved counts and hackable lists
   k0 := [3]int{}
   H := make([][]int, 3) // store participant indices hackable
   for i := 0; i < n; i++ {
       for j := 0; j < 3; j++ {
           t := times[i][j]
           if t != 0 {
               k0[j]++
           }
           if t < 0 {
               H[j] = append(H[j], i)
           }
       }
   }
   // sort hackable by least harm (later solves): abs(t) descending
   for j := 0; j < 3; j++ {
       sort.Slice(H[j], func(a, b int) bool {
           ta := -times[H[j][a]][j]
           tb := -times[H[j][b]][j]
           return ta > tb
       })
   }
   scores := []int{500, 1000, 1500, 2000, 2500, 3000}
   bestRank := n + 1
   // iterate s combinations
   for _, s0 := range scores {
       for _, s1 := range scores {
           for _, s2 := range scores {
               s := [3]int{s0, s1, s2}
               valid := true
               minx := [3]int{}
               maxx := [3]int{}
               for j := 0; j < 3; j++ {
                   // find k interval for s[j]
                   L, R := 0, n
                   switch s[j] {
                   case 500:
                       L = n/2 + 1; R = n
                   case 1000:
                       L = n/4 + 1; R = n/2
                   case 1500:
                       L = n/8 + 1; R = n/4
                   case 2000:
                       L = n/16 + 1; R = n/8
                   case 2500:
                       L = n/32 + 1; R = n/16
                   case 3000:
                       L = 0; R = n/32
                   }
                   // require k0[j] - x in [L,R], x in [k0-R, k0-L]
                   lo := k0[j] - R
                   hi := k0[j] - L
                   if lo < 0 {
                       lo = 0
                   }
                   if hi > len(H[j]) {
                       hi = len(H[j])
                   }
                   if lo > hi {
                       valid = false
                       break
                   }
                   minx[j] = lo
                   maxx[j] = hi
               }
               if !valid {
                   continue
               }
               // try extreme hack counts: all combinations min or max
               for mask := 0; mask < 8; mask++ {
                   x := [3]int{}
                   for j := 0; j < 3; j++ {
                       if (mask>>j)&1 == 0 {
                           x[j] = minx[j]
                       } else {
                           x[j] = maxx[j]
                       }
                   }
                   // apply hacks
                   hacked := make([][3]bool, n)
                   for j := 0; j < 3; j++ {
                       for t := 0; t < x[j]; t++ {
                           i := H[j][t]
                           hacked[i][j] = true
                       }
                   }
                   // compute our score
                   Sc := 0.0
                   for j := 0; j < 3; j++ {
                       t := times[0][j]
                       if t > 0 {
                           Sc += float64(s[j]) * float64(250-t) / 250.0
                       }
                   }
                   Sc += float64(100 * (x[0] + x[1] + x[2]))
                   // count rank
                   rank := 1
                   for i := 1; i < n; i++ {
                       Pi := 0.0
                       for j := 0; j < 3; j++ {
                           t := times[i][j]
                           if t > 0 && !hacked[i][j] {
                               Pi += float64(s[j]) * float64(250-t) / 250.0
                           }
                       }
                       if Pi > Sc {
                           rank++
                           if rank >= bestRank {
                               break
                           }
                       }
                   }
                   if rank < bestRank {
                       bestRank = rank
                   }
               }
           }
       }
   }
   fmt.Println(bestRank)
}

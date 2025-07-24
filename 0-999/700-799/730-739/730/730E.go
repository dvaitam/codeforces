package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n)
   d := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i], &d[i])
   }

   // generate candidate orders using heuristics and pick max
   // cost computes total applause for given order
   cost := func(ord []int) int {
       s2 := make([]int, n)
       for i := 0; i < n; i++ {
           s2[i] = a[i]
       }
       tot := 0
       for _, j := range ord {
           x := rank(s2, j)
           s2[j] = a[j] + d[j]
           y := rank(s2, j)
           if x > y {
               tot += x - y
           } else {
               tot += y - x
           }
       }
       return tot
   }
   var cands [][]int
   // group by d sign
   drops := make([]int, 0, n)
   gains := make([]int, 0, n)
   for i := 0; i < n; i++ {
       if d[i] < 0 {
           drops = append(drops, i)
       } else {
           gains = append(gains, i)
       }
   }
   // key1: drops by final asc, gains by final desc
   d1 := make([]int, len(drops))
   copy(d1, drops)
   sort.Slice(d1, func(i, j int) bool {
       fi, fj := a[d1[i]]+d[d1[i]], a[d1[j]]+d[d1[j]]
       if fi != fj {
           return fi < fj
       }
       return d1[i] < d1[j]
   })
   g1 := make([]int, len(gains))
   copy(g1, gains)
   sort.Slice(g1, func(i, j int) bool {
       fi, fj := a[g1[i]]+d[g1[i]], a[g1[j]]+d[g1[j]]
       if fi != fj {
           return fi > fj
       }
       return g1[i] < g1[j]
   })
   cands = append(cands, append(append([]int{}, d1...), g1...))
   // key2: drops by initial desc, gains by initial asc
   d2 := make([]int, len(drops))
   copy(d2, drops)
   sort.Slice(d2, func(i, j int) bool {
       ai, aj := a[d2[i]], a[d2[j]]
       if ai != aj {
           return ai > aj
       }
       return d2[i] < d2[j]
   })
   g2 := make([]int, len(gains))
   copy(g2, gains)
   sort.Slice(g2, func(i, j int) bool {
       ai, aj := a[g2[i]], a[g2[j]]
       if ai != aj {
           return ai < aj
       }
       return g2[i] < g2[j]
   })
   cands = append(cands, append(append([]int{}, d2...), g2...))
   // other single-key sorts
   all := make([]int, n)
   for i := range all {
       all[i] = i
   }
   // d asc
   ord := make([]int, n)
   copy(ord, all)
   sort.Slice(ord, func(i, j int) bool {
       if d[ord[i]] != d[ord[j]] {
           return d[ord[i]] < d[ord[j]]
       }
       return ord[i] < ord[j]
   })
   cands = append(cands, ord)
   // d desc
   ord2 := make([]int, n)
   copy(ord2, all)
   sort.Slice(ord2, func(i, j int) bool {
       if d[ord2[i]] != d[ord2[j]] {
           return d[ord2[i]] > d[ord2[j]]
       }
       return ord2[i] < ord2[j]
   })
   cands = append(cands, ord2)
   // final asc
   ord3 := make([]int, n)
   copy(ord3, all)
   sort.Slice(ord3, func(i, j int) bool {
       fi, fj := a[ord3[i]]+d[ord3[i]], a[ord3[j]]+d[ord3[j]]
       if fi != fj {
           return fi < fj
       }
       return ord3[i] < ord3[j]
   })
   cands = append(cands, ord3)
   // final desc
   ord4 := make([]int, n)
   copy(ord4, all)
   sort.Slice(ord4, func(i, j int) bool {
       fi, fj := a[ord4[i]]+d[ord4[i]], a[ord4[j]]+d[ord4[j]]
       if fi != fj {
           return fi > fj
       }
       return ord4[i] < ord4[j]
   })
   cands = append(cands, ord4)
   // initial asc
   ord5 := make([]int, n)
   copy(ord5, all)
   sort.Slice(ord5, func(i, j int) bool {
       if a[ord5[i]] != a[ord5[j]] {
           return a[ord5[i]] < a[ord5[j]]
       }
       return ord5[i] < ord5[j]
   })
   cands = append(cands, ord5)
   // initial desc
   ord6 := make([]int, n)
   copy(ord6, all)
   sort.Slice(ord6, func(i, j int) bool {
       if a[ord6[i]] != a[ord6[j]] {
           return a[ord6[i]] > a[ord6[j]]
       }
       return ord6[i] < ord6[j]
   })
   cands = append(cands, ord6)
   // include reversals
   sz := len(cands)
   for i := 0; i < sz; i++ {
       rev := make([]int, n)
       for j := 0; j < n; j++ {
           rev[j] = cands[i][n-1-j]
       }
       cands = append(cands, rev)
   }
   // deduplicate
   seen := make(map[string]bool)
   best := 0
   for _, ord0 := range cands {
       key := fmt.Sprint(ord0)
       if seen[key] {
           continue
       }
       seen[key] = true
       c := cost(ord0)
       if c > best {
           best = c
       }
   }
   fmt.Println(best)
}

// rank computes 1-based rank of j among all teams sorted by score desc, id asc
func rank(s []int, j int) int {
   n := len(s)
   // count how many k are before j
   cnt := 0
   for k := 0; k < n; k++ {
       if k == j {
           continue
       }
       if s[k] > s[j] || (s[k] == s[j] && k < j) {
           cnt++
       }
   }
   return cnt + 1
}

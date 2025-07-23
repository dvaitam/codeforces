package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// pair represents swimming speed s and running speed r
type pair struct{
   s, r int
}

// node represents a unique competitor point for hull
type node struct{
   s, r int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // Map each (s,r) to list of indices
   idxMap := make(map[pair][]int, n)
   for i := 1; i <= n; i++ {
       var s, r int
       fmt.Fscan(in, &s, &r)
       idxMap[pair{s, r}] = append(idxMap[pair{s, r}], i)
   }
   // Unique keys
   uniq := make([]pair, 0, len(idxMap))
   for p := range idxMap {
       uniq = append(uniq, p)
   }
   // Sort by s desc, r desc
   sort.Slice(uniq, func(i, j int) bool {
       if uniq[i].s != uniq[j].s {
           return uniq[i].s > uniq[j].s
       }
       return uniq[i].r > uniq[j].r
   })
   // Pareto filter: keep only those with r > max_r so far
   pareto := make([]pair, 0, len(uniq))
   maxR := 0
   for _, p := range uniq {
       if p.r > maxR {
           pareto = append(pareto, p)
           maxR = p.r
       }
   }
   // Build convex hull on points (1/s,1/r) sorted by a asc (s desc), b asc (r desc)
   // pareto is already sorted by s desc, r desc, and r asc increases as s decreases
   hull := make([]node, 0, len(pareto))
   for _, p := range pareto {
       cur := node{p.s, p.r}
       // pop while last three make middle unnecessary
       for len(hull) >= 2 {
           j := len(hull) - 1
           i := j - 1
           if bad(hull[i], hull[j], cur) {
               // remove hull[j]
               hull = hull[:j]
               continue
           }
           break
       }
       hull = append(hull, cur)
   }
   // Collect result indices
   res := make([]int, 0, n)
   for _, nd := range hull {
       for _, id := range idxMap[pair{nd.s, nd.r}] {
           res = append(res, id)
       }
   }
   sort.Ints(res)
   // Output
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for i, v := range res {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, v)
   }
   out.WriteByte('\n')
}

// bad returns true if middle point j is unnecessary
// bad returns true if middle point pj should be removed for lower hull
func bad(pi, pj, pk node) bool {
   // compute cross product of (pi,pj,pk) in (a=1/s,b=1/r) space
   // cross ~ (Aj-Ai)*(Bk-Bj) - (Bj-Bi)*(Ak-Aj)
   // Multiply by common positive denom to avoid fractions:
   // left = (si - sj)*(rj - rk)*ri*sk
   // right= (ri - rj)*(sj - sk)*si*rk
   si, sj, sk := int64(pi.s), int64(pj.s), int64(pk.s)
   ri, rj, rk := int64(pi.r), int64(pj.r), int64(pk.r)
   left := (si - sj) * (rj - rk) * ri * sk
   right := (ri - rj) * (sj - sk) * si * rk
   // remove pj if cross < 0 => left < right
   return left < right
}

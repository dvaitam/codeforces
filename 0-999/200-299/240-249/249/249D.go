package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
   "strings"
)

// BIT for prefix maximum
type BIT struct {
   n    int
   tree []int
}

func NewBIT(n int) *BIT {
   return &BIT{n: n, tree: make([]int, n+1)}
}

// update index i with value val (keep max)
func (b *BIT) update(i, val int) {
   for ; i <= b.n; i += i & -i {
       if b.tree[i] < val {
           b.tree[i] = val
       }
   }
}

// query maximum on prefix [1..i]
func (b *BIT) query(i int) int {
   res := 0
   for ; i > 0; i -= i & -i {
       if b.tree[i] > res {
           res = b.tree[i]
       }
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var frac1, frac2 string
   fmt.Fscan(reader, &frac1, &frac2)
   parts := strings.Split(frac1, "/")
   a, _ := strconv.Atoi(parts[0])
   b, _ := strconv.Atoi(parts[1])
   parts = strings.Split(frac2, "/")
   c, _ := strconv.Atoi(parts[0])
   d, _ := strconv.Atoi(parts[1])
   pts := make([]struct {
       v1, v2 int64
       dp     int
       idx    int
       revIdx int
   }, n)
   v2vals := make([]int64, 0, n+1)
   for i := 0; i < n; i++ {
       var xi, yi int64
       fmt.Fscan(reader, &xi, &yi)
       v1 := int64(b)*yi - int64(a)*xi
       v2 := int64(d)*yi - int64(c)*xi
       pts[i].v1 = v1
       pts[i].v2 = v2
       v2vals = append(v2vals, v2)
   }
   // include origin v2 = 0 for initial reachability
   v2vals = append(v2vals, 0)
   sort.Slice(v2vals, func(i, j int) bool { return v2vals[i] < v2vals[j] })
   uniq := v2vals[:1]
   for i := 1; i < len(v2vals); i++ {
       if v2vals[i] != v2vals[i-1] {
           uniq = append(uniq, v2vals[i])
       }
   }
   m := len(uniq)
   // map v2 to index
   mp := make(map[int64]int, m)
   for i, v := range uniq {
       mp[v] = i + 1
   }
   // assign idx and revIdx
   for i := range pts {
       idx := mp[pts[i].v2]
       pts[i].idx = idx
       pts[i].revIdx = m - idx + 1
   }
   // sort by v1 ascending
   sort.Slice(pts, func(i, j int) bool {
       return pts[i].v1 < pts[j].v1
   })
   bit := NewBIT(m)
   ans := 0
   // process by groups of equal v1
   for i := 0; i < n; {
       j := i + 1
       for j < n && pts[j].v1 == pts[i].v1 {
           j++
       }
       // compute dp for group
       for k := i; k < j; k++ {
           p := &pts[k]
           dpPrev := 0
           if p.revIdx > 1 {
               dpPrev = bit.query(p.revIdx - 1)
           }
           if dpPrev == 0 {
               if p.v1 > 0 && p.v2 < 0 {
                   p.dp = 1
               } else {
                   p.dp = 0
               }
           } else {
               p.dp = dpPrev + 1
           }
           if p.dp > ans {
               ans = p.dp
           }
       }
       // update BIT
       for k := i; k < j; k++ {
           p := &pts[k]
           if p.dp > 0 {
               bit.update(p.revIdx, p.dp)
           }
       }
       i = j
   }
   // print result
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprint(writer, ans)
}

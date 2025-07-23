package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Video struct {
   l, r int
   idx  int
}
type Channel struct {
   a, b int
   c    int
   idx  int
}

// BIT for max pair (value, index)
type pair struct{ v, id int }
type BIT struct {
   n int
   t []pair
}
func NewBIT(n int) *BIT {
   return &BIT{n, make([]pair, n+1)}
}
func (b *BIT) Update(i int, p pair) {
   for x := i; x <= b.n; x += x & -x {
       if p.v > b.t[x].v {
           b.t[x] = p
       }
   }
}
// Query max in [i..n]
func (b *BIT) QuerySuffix(i int) pair {
   // convert suffix query to prefix on reversed: we scan all x>=i
   res := pair{0, 0}
   // we can traverse manually from end? Instead, build BIT that stores at reversed index.
   // But simpler: we store BIT for suffix by reversing indices when building
   // However here we assume QuerySuffix called on reversed BIT, so alias of QueryPrefix
   for x := i; x > 0; x -= x & -x {
       if b.t[x].v > res.v {
           res = b.t[x]
       }
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   vids := make([]Video, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &vids[i].l, &vids[i].r)
       vids[i].idx = i + 1
   }
   chans := make([]Channel, m)
   bs := make([]int, m)
   for j := 0; j < m; j++ {
       fmt.Fscan(in, &chans[j].a, &chans[j].b, &chans[j].c)
       chans[j].idx = j + 1
       bs[j] = chans[j].b
   }
   // Prepare for prefix (li <= a_j)
   sort.Slice(vids, func(i, j int) bool { return vids[i].l < vids[j].l })
   chByA := make([]Channel, m)
   copy(chByA, chans)
   sort.Slice(chByA, func(i, j int) bool { return chByA[i].a < chByA[j].a })
   // prefix sweep
   best := int64(0)
   bestVi, bestCj := 0, 0
   vp := 0
   maxR, maxRi := 0, 0
   for _, ch := range chByA {
       for vp < n && vids[vp].l <= ch.a {
           if vids[vp].r > maxR {
               maxR = vids[vp].r
               maxRi = vids[vp].idx
           }
           vp++
       }
       // best overlap for prefix videos
       if maxR > ch.a {
           end := maxR
           if ch.b < end {
               end = ch.b
           }
           if end > ch.a {
               length := end - ch.a
               val := int64(length) * int64(ch.c)
               if val > best {
                   best = val
                   bestVi = maxRi
                   bestCj = ch.idx
               }
           }
       }
   }
   // suffix sweep (ri >= b_j)
   sort.Slice(vids, func(i, j int) bool { return vids[i].r > vids[j].r })
   chByB := make([]Channel, m)
   copy(chByB, chans)
   sort.Slice(chByB, func(i, j int) bool { return chByB[i].b > chByB[j].b })
   vp = 0
   minL, minLi := int(1e18), 0
   for _, ch := range chByB {
       for vp < n && vids[vp].r >= ch.b {
           if vids[vp].l < minL {
               minL = vids[vp].l
               minLi = vids[vp].idx
           }
           vp++
       }
       if minL < ch.b {
           start := ch.b
           if ch.a > minL {
               start = ch.a
           }
           if ch.b > start {
               length := ch.b - start
               val := int64(length) * int64(ch.c)
               if val > best {
                   best = val
                   bestVi = minLi
                   bestCj = ch.idx
               }
           }
       }
   }
   // type B: channel covers video
   // compress b_j reversed for suffix BIT
   // We'll store channels reversed: map b->pos = K - rank(b) +1
   sort.Ints(bs)
   bs = unique(bs)
   K := len(bs)
   // mapping original b to reversed index
   // revPos(b) = K - (lower_bound(bs,b))
   bit := NewBIT(K)
   // prepare channels sorted by a ascending
   sort.Slice(chans, func(i, j int) bool { return chans[i].a < chans[j].a })
   // videos sorted by l ascending
   sort.Slice(vids, func(i, j int) bool { return vids[i].l < vids[j].l })
   cp := 0
   for _, v := range vids {
       // add channels with a_j <= l_i
       for cp < m && chans[cp].a <= v.l {
           // compute reversed position
           bi := chans[cp].b
           pos := K - sort.SearchInts(bs, bi)
           bit.Update(pos, pair{chans[cp].c, chans[cp].idx})
           cp++
       }
       // query channels with b_j >= r_i => reversed pos <= K - idx_of(r_i)
       // lower_bound on r_i
       rb := v.r
       // first index >= r_i
       lo := sort.SearchInts(bs, rb)
       if lo < K {
           pos := K - lo
           p := bit.QuerySuffix(pos)
           if p.v > 0 {
               length := v.r - v.l
               val := int64(length) * int64(p.v)
               if val > best {
                   best = val
                   bestVi = v.idx
                   bestCj = p.id
               }
           }
       }
   }
   // output
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   if best <= 0 {
       fmt.Fprintln(out, 0)
   } else {
       fmt.Fprintln(out, best)
       fmt.Fprintf(out, "%d %d\n", bestVi, bestCj)
   }
}

func unique(a []int) []int {
   j := 0
   for i := 0; i < len(a); i++ {
       if j == 0 || a[i] != a[j-1] {
           a[j] = a[i]
           j++
       }
   }
   return a[:j]
}

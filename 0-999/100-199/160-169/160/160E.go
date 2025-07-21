package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
)

type Bus struct {
   s, f int
   t    int
   id   int
}
type Query struct {
   l, r, b int
   id       int
   low, high int
}

// segment tree for range max query and point update
type SegTree struct {
   n    int
   tree []int
}
func NewSegTree(n int) *SegTree {
   size := 1
   for size < n {
       size <<= 1
   }
   return &SegTree{n: size, tree: make([]int, 2*size)}
}
func (st *SegTree) Update(pos, val int) {
   i := pos + st.n
   if st.tree[i] >= val {
       return
   }
   st.tree[i] = val
   for i >>= 1; i > 0; i >>= 1 {
       if st.tree[2*i] > st.tree[2*i+1] {
           st.tree[i] = st.tree[2*i]
       } else {
           st.tree[i] = st.tree[2*i+1]
       }
   }
}
// query max in [0..r]
func (st *SegTree) Query(r int) int {
   l := st.n
   rr := r + st.n
   maxv := 0
   // query [0, rr]
   for l <= rr {
       if l&1 == 1 {
           if st.tree[l] > maxv {
               maxv = st.tree[l]
           }
           l++
       }
       if rr&1 == 0 {
           if st.tree[rr] > maxv {
               maxv = st.tree[rr]
           }
           rr--
       }
       l >>= 1
       rr >>= 1
   }
   return maxv
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   // better scanner
   var n, m int
   fmt.Fscan(in, &n, &m)
   buses := make([]Bus, n)
   for i := 0; i < n; i++ {
       var si, fi, ti int
       fmt.Fscan(in, &si, &fi, &ti)
       buses[i] = Bus{si, fi, ti, i + 1}
   }
   queries := make([]Query, m)
   for i := 0; i < m; i++ {
       var li, ri, bi int
       fmt.Fscan(in, &li, &ri, &bi)
       queries[i] = Query{li, ri, bi, i, 0, n - 1}
   }
   // sort buses by t ascending
   sort.Slice(buses, func(i, j int) bool { return buses[i].t < buses[j].t })
   tList := make([]int, n)
   for i := range buses {
       tList[i] = buses[i].t
   }
   // initial low based on b
   for i := range queries {
       lb := sort.Search(n, func(j int) bool { return tList[j] >= queries[i].b })
       if lb >= n {
           queries[i].low = n
           queries[i].high = -1
       } else {
           queries[i].low = lb
           queries[i].high = n - 1
       }
   }
   // compress s and l coordinates
   coords := make([]int, 0, n+m)
   for _, b := range buses {
       coords = append(coords, b.s)
   }
   for _, q := range queries {
       coords = append(coords, q.l)
   }
   sort.Ints(coords)
   uniq := coords[:1]
   for i := 1; i < len(coords); i++ {
       if coords[i] != coords[i-1] {
           uniq = append(uniq, coords[i])
       }
   }
   // map values
   getPos := func(x int) int {
       return sort.SearchInts(uniq, x)
   }
   sPos := make([]int, n)
   for i := range buses {
       sPos[i] = getPos(buses[i].s)
   }
   lPos := make([]int, m)
   for i := range queries {
       lPos[i] = getPos(queries[i].l)
   }
   // parallel binary search
   type Task struct{ mid, qi int }
   for {
       // collect active queries
       tasks := make([]Task, 0, m)
       for i := range queries {
           if queries[i].low <= queries[i].high {
               mid := (queries[i].low + queries[i].high) >> 1
               tasks = append(tasks, Task{mid, i})
           }
       }
       if len(tasks) == 0 {
           break
       }
       // process in ascending mid for prefix DS
       sort.Slice(tasks, func(i, j int) bool { return tasks[i].mid < tasks[j].mid })
       st := NewSegTree(len(uniq))
       p := -1
       for _, t := range tasks {
           mid, qi := t.mid, t.qi
           // add buses[0..mid]
           for p < mid {
               p++
               st.Update(sPos[p], buses[p].f)
           }
           // query max f for s <= l
           if st.Query(lPos[qi]) >= queries[qi].r {
               // possible in this prefix, search left for earlier
               queries[qi].high = mid - 1
           } else {
               // not possible, search right
               queries[qi].low = mid + 1
           }
       }
   }
   // prepare answers
   ans := make([]int, m)
   for i := range queries {
       if queries[i].low < n {
           ans[queries[i].id] = buses[queries[i].low].id
       } else {
           ans[queries[i].id] = -1
       }
   }
   // output
   for i, v := range ans {
       if i > 0 {
           out.WriteByte(' ')
       }
       out.WriteString(strconv.Itoa(v))
   }
   out.WriteByte('\n')
}

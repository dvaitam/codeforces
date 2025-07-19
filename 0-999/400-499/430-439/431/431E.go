package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Query struct {
   typ  int
   idx  int   // for type 1: position (1-based)
   newv int   // for type 1: new value
   v    int64 // for type 2: value v
}

var (
   n, q       int
   h          []int
   queries    []Query
   disc       []int
   cntFen     []int
   sumFen     []int64
)

func lowbit(x int) int {
   return x & -x
}

// update Fenwick tree at position i (1-based) with count delta
func update(i, delta, size int) {
   d := int64(disc[i-1]) * int64(delta)
   for ; i <= size; i += lowbit(i) {
       cntFen[i] += delta
       sumFen[i] += d
   }
}

// query Fenwick tree prefix [1..i]
func queryFen(i int) (cnt int, sum int64) {
   for ; i > 0; i -= lowbit(i) {
       cnt += cntFen[i]
       sum += sumFen[i]
   }
   return
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &n, &q)
   h = make([]int, n)
   disc = make([]int, 0, n+q)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &h[i])
       disc = append(disc, h[i])
   }
   queries = make([]Query, q)
   for i := 0; i < q; i++ {
       var typ int
       fmt.Fscan(reader, &typ)
       queries[i].typ = typ
       if typ == 1 {
           var idx, newv int
           fmt.Fscan(reader, &idx, &newv)
           queries[i].idx = idx
           queries[i].newv = newv
           disc = append(disc, newv)
       } else {
           var v int64
           fmt.Fscan(reader, &v)
           queries[i].v = v
       }
   }
   // coordinate compression
   sort.Ints(disc)
   m := 0
   for i := 0; i < len(disc); i++ {
       if i == 0 || disc[i] != disc[i-1] {
           disc[m] = disc[i]
           m++
       }
   }
   disc = disc[:m]
   // init Fenwick trees
   cntFen = make([]int, m+1)
   sumFen = make([]int64, m+1)
   // build initial counts
   for i := 0; i < n; i++ {
       tid := sort.SearchInts(disc, h[i]) + 1
       update(tid, 1, m)
   }
   // process queries
   for _, qr := range queries {
       if qr.typ == 1 {
           idx := qr.idx - 1
           // remove old value
           oldv := h[idx]
           tidOld := sort.SearchInts(disc, oldv) + 1
           update(tidOld, -1, m)
           // add new value
           tidNew := sort.SearchInts(disc, qr.newv) + 1
           update(tidNew, 1, m)
           h[idx] = qr.newv
       } else {
           v := qr.v
           l, r := 1, m
           nl := 1
           for l <= r {
               mid := (l + r) >> 1
               cnt, sum := queryFen(mid)
               if v+sum > int64(cnt)*int64(disc[mid-1]) {
                   nl = mid
                   l = mid + 1
               } else {
                   r = mid - 1
               }
           }
           cnt, sum := queryFen(nl)
           avg := float64(v+sum) / float64(cnt)
           fmt.Fprintf(writer, "%f\n", avg)
       }
   }
}

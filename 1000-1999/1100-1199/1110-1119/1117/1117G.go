package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// BIT implements a Fenwick tree for int64 values.
type BIT struct {
   n    int
   tree []int64
}

// NewBIT creates a BIT that can handle indices [0..n-1].
func NewBIT(n int) *BIT {
   // allocate a bit more space for 1-based indexing
   return &BIT{n + 5, make([]int64, n+10)}
}

// Update adds x at index u (0-based).
func (b *BIT) Update(u int, x int64) {
   for idx := u + 1; idx <= b.n; idx += idx & -idx {
       b.tree[idx] += x
   }
}

// Query returns the prefix sum [0..u] (0-based).
func (b *BIT) Query(u int) int64 {
   var res int64
   for idx := u + 1; idx > 0; idx -= idx & -idx {
       res += b.tree[idx]
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, q int
   fmt.Fscan(reader, &n, &q)
   a := make([]int, n+2)
   // sentinel values
   a[0], a[n+1] = n+1, n+1
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   l := make([]int, n+2)
   r := make([]int, n+2)
   stack := make([]int, n+2)
   // compute previous greater element
   vol := 1
   stack[1] = 0
   for i := 1; i <= n; i++ {
       for vol > 0 && a[stack[vol]] <= a[i] {
           vol--
       }
       l[i] = stack[vol]
       vol++
       stack[vol] = i
   }
   // compute next greater element
   vol = 1
   stack[1] = n + 1
   for i := n; i >= 1; i-- {
       for vol > 0 && a[stack[vol]] <= a[i] {
           vol--
       }
       r[i] = stack[vol]
       vol++
       stack[vol] = i
   }

   // read queries
   type Query struct{ L, R, idx int }
   queries := make([]Query, q)
   for i := 0; i < q; i++ {
       fmt.Fscan(reader, &queries[i].L, &queries[i].R)
       queries[i].idx = i
   }
   // sort queries by R
   sort.Slice(queries, func(i, j int) bool {
       return queries[i].R < queries[j].R
   })

   // tags for right endpoints
   tags := make([][]int, n+2)
   // BIT structures
   f := NewBIT(n + 5)
   g := NewBIT(n + 5)
   h := NewBIT(n + 5)
   cnt := NewBIT(n + 5)
   sumv := NewBIT(n + 5)
   // prefix sum of l
   sm := make([]int64, n+2)
   ans := make([]int64, q)
   qi := 0

   // main processing
   for i := 1; i <= n; i++ {
       // apply deferred updates for positions ending here
       for _, u := range tags[i] {
           f.Update(u, int64(l[u]+1))
           g.Update(u, -1)
           h.Update(u, int64(i - l[u] - 1))
       }
       // new interval starting at i
       f.Update(i, int64(-l[i] - 1))
       g.Update(i, 1)
       tags[r[i]] = append(tags[r[i]], i)
       sm[i] = sm[i-1] + int64(l[i])
       cnt.Update(l[i], 1)
       sumv.Update(l[i], int64(l[i]))
       // answer queries with right endpoint i
       for qi < q && queries[qi].R == i {
           L := queries[qi].L
           R := queries[qi].R
           idx := queries[qi].idx
           var res int64
           res += h.Query(R) - h.Query(L - 1)
           res += f.Query(R) - f.Query(L - 1)
           res += (g.Query(R) - g.Query(L - 1)) * int64(R + 1)
           res += sumv.Query(L - 1) - sm[L - 1]
           res -= (cnt.Query(L - 1) - int64(L - 1)) * int64(L - 1)
           ans[idx] = res
           qi++
       }
   }
   // output answers
   for i := 0; i < q; i++ {
       fmt.Fprintf(writer, "%d ", ans[i])
   }
}

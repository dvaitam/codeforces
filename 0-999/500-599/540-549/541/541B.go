package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Fenwick tree for ints
type Fenwick struct {
   n    int
   tree []int
}

func NewFenwick(n int) *Fenwick {
   return &Fenwick{n: n, tree: make([]int, n+1)}
}

func (f *Fenwick) Add(i, v int) {
   for ; i <= f.n; i += i & -i {
       f.tree[i] += v
   }
}

func (f *Fenwick) Sum(i int) int {
   s := 0
   for ; i > 0; i -= i & -i {
       s += f.tree[i]
   }
   return s
}

func max64(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var r int64
   if _, err := fmt.Fscan(in, &n, &r); err != nil {
       return
   }
   intervals := make([]struct{l, r int64}, 0, n)
   for i := 0; i < n; i++ {
       var hi, ti int64
       fmt.Fscan(in, &hi, &ti)
       // At time t, duck covers [hi-t, ti-t], intersects 0 when hi <= t <= ti
       if ti < 0 {
           continue
       }
       l := hi
       if l < 0 {
           l = 0
       }
       intervals = append(intervals, struct{l, r int64}{l, ti})
   }
   m := len(intervals)
   if m == 0 {
       fmt.Println(0)
       return
   }
   // sort by r ascending
   sort.Slice(intervals, func(i, j int) bool {
       return intervals[i].r < intervals[j].r
   })
   // prepare arrays (1-indexed)
   rArr := make([]int64, m+1)
   // lArr not directly needed
   // prepare list sorted by l
   type pair struct{l int64; idx int}
   llist := make([]pair, m)
   for i := 1; i <= m; i++ {
       rArr[i] = intervals[i-1].r
       llist[i-1] = pair{l: intervals[i-1].l, idx: i}
   }
   sort.Slice(llist, func(i, j int) bool { return llist[i].l < llist[j].l })
   // Fenwick for active intervals by l <= current r_i
   fw := NewFenwick(m)
   count := make([]int, m+1)
   ptr := 0
   // compute count[i]
   for i := 1; i <= m; i++ {
       ri := rArr[i]
       for ptr < m && llist[ptr].l <= ri {
           fw.Add(llist[ptr].idx, 1)
           ptr++
       }
       // count intervals j>=i with l_j <= r_i: total added - sum up to i-1
       count[i] = fw.Sum(m) - fw.Sum(i-1)
   }
   // compute p[i]
   p := make([]int, m+1)
   for i := 1; i <= m; i++ {
       // find largest k with rArr[k] <= rArr[i] - r
       val := rArr[i] - r
       // search first index with rArr[k] > val
       k := sort.Search(m, func(j int) bool { return rArr[j+1] > val })
       // k is number of elements <= val
       p[i] = k
   }
   // DP
   dp := make([]int64, m+1)
   for i := 1; i <= m; i++ {
       without := dp[i-1]
       with := dp[p[i]] + int64(count[i])
       dp[i] = max64(without, with)
   }
   fmt.Println(dp[m])
}

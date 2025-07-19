package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type pair struct { first, second int }

func maxPair(a, b pair) pair {
   if a.first > b.first {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   a := make([]pair, n+5)
   for i := 1; i <= n; i++ {
       var ai int
       fmt.Fscan(reader, &ai)
       a[i] = pair{-ai, i}
   }
   // sentinel default a[n+1] = {0,0}
   // sort by first (which is -ai) increasing, so original ai descending
   arr := a[1 : n+1]
   sort.Slice(arr, func(i, j int) bool {
       return arr[i].first < arr[j].first
   })
   const K = 12
   // sparse table f[k][i]
   f := make([][]pair, K)
   for k := 0; k < K; k++ {
       f[k] = make([]pair, n+5)
   }
   // build level 0
   for i := 1; i <= n; i++ {
       f[0][i] = pair{a[i+1].first - a[i].first, i}
   }
   // build sparse table
   for k := 1; k < K; k++ {
       shift := 1 << (k - 1)
       for i := 1; i <= n; i++ {
           if i+shift <= n {
               f[k][i] = maxPair(f[k-1][i], f[k-1][i+shift])
           } else {
               f[k][i] = f[k-1][i]
           }
       }
   }
   // log table
   lg := make([]int, n+5)
   for i := 3; i <= n; i++ {
       lg[i] = lg[(i+1)>>1] + 1
   }
   // rmq
   rmq := func(l, r int) pair {
       length := r - l + 1
       k := lg[length]
       b1 := f[k][l]
       b2 := f[k][r-(1<<k)+1]
       if b1.first > b2.first {
           return b1
       }
       return b2
   }

   const INF = 1000000000
   d1, d2, d3 := -INF, -INF, -INF
   p1, p2, p3 := 0, 0, 0
   // search best split
   for i := 1; i <= n; i++ {
       for j := i + 1; j <= n; j++ {
           len1 := j - i
           if i > 2*len1 || len1 > 2*i {
               continue
           }
           // compute l and r for third segment diff
           // m = max(i, len1)
           m := i
           if len1 > m {
               m = len1
           }
           // l = ceil(m/2)
           l := (m + 1) >> 1
           // r = 2*m, capped by remaining n-j
           r := m * 2
           rem := n - j
           if r > rem {
               r = rem
           }
           if l <= r && j+l <= n {
               x := rmq(j+l, j+r)
               d1cand := a[i+1].first - a[i].first
               d2cand := a[j+1].first - a[j].first
               if d1cand > d1 || (d1cand == d1 && (d2cand > d2 || (d2cand == d2 && x.first > d3))) {
                   d1 = d1cand
                   p1 = i
                   d2 = d2cand
                   p2 = j
                   d3 = x.first
                   p3 = x.second
               }
           }
       }
   }
   // assign answers
   ans := make([]int, n+1)
   for i := 1; i <= p1; i++ {
       ans[a[i].second] = 1
   }
   for i := p1 + 1; i <= p2; i++ {
       ans[a[i].second] = 2
   }
   for i := p2 + 1; i <= p3; i++ {
       ans[a[i].second] = 3
   }
   for i := p3 + 1; i <= n; i++ {
       ans[a[i].second] = -1
   }
   // output
   for i := 1; i <= n; i++ {
       fmt.Fprint(writer, ans[i], " ")
   }
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

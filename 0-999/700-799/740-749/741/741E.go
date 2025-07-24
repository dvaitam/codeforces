package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var S, T string
   var q int
   if _, err := fmt.Fscan(reader, &S, &T, &q); err != nil {
       return
   }
   n := len(S)
   m := len(T)
   // compute Z-array for T#S to get lcp of T and S[i:]
   V := T + "#" + S
   Z := make([]int, len(V))
   l, r := 0, 0
   for i := 1; i < len(V); i++ {
       if i <= r {
           k := i - l
           if Z[k] < r-i+1 {
               Z[i] = Z[k]
           } else {
               j := r + 1
               for j < len(V) && V[j] == V[j-i] {
                   j++
               }
               Z[i] = j - i
               l, r = i, j-1
           }
       } else {
           j := 0
           for i+j < len(V) && V[j] == V[i+j] {
               j++
           }
           Z[i] = j
           if j > 0 {
               l, r = i, i+j-1
           }
       }
   }
   zTS := make([]int, n+1)
   for i := 0; i <= n; i++ {
       idx := m + 1 + i
       if idx < len(Z) {
           zTS[i] = Z[idx]
           if zTS[i] > m {
               zTS[i] = m
           }
       } else {
           zTS[i] = 0
       }
   }
   // suffix array for S
   sa, rk := buildSA(S)
   // comparator for A_i
   type idx struct{ i int }
   arr := make([]idx, n+1)
   for i := 0; i <= n; i++ {
       arr[i] = idx{i}
   }
   compare := func(a, b idx) bool {
       i, j := a.i, b.i
       if i == j {
           return false
       }
       if i < j {
           // compare A_i vs A_j
           d := j - i
           lcp := zTS[i]
           if lcp < m && lcp < d {
               return T[lcp] < S[i+lcp]
           }
           if m < d {
               // compare suffix S at i vs i+m
               return rk[i] < rk[i+m]
           }
           if m > d {
               // both in T
               return T[d] < T[0]
           }
           // m == d
           return S[i] < T[0]
       } else {
           // i > j: invert
           return !compare(b, a)
       }
   }
   sort.Slice(arr, func(a, b int) bool { return compare(arr[a], arr[b]) })
   rankA := make([]int, n+1)
   for idx, v := range arr {
       rankA[v.i] = idx
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // answer queries by brute
   for ; q > 0; q-- {
       var lq, rq, k, x, y int
       fmt.Fscan(reader, &lq, &rq, &k, &x, &y)
       bestR, bestI := n+2, -1
       for i := lq; i <= rq; i++ {
           rmd := i % k
           if rmd < x || rmd > y {
               continue
           }
           if rankA[i] < bestR || (rankA[i] == bestR && i < bestI) {
               bestR = rankA[i]
               bestI = i
           }
       }
       if bestI < 0 {
           fmt.Fprintln(writer, -1)
       } else {
           fmt.Fprintln(writer, bestI)
       }
   }
}

// buildSA builds suffix array and rank array for s
func buildSA(s string) ([]int, []int) {
   n := len(s)
   sa := make([]int, n)
   rk := make([]int, n)
   tmp := make([]int, n)
   for i := 0; i < n; i++ {
       sa[i] = i
       rk[i] = int(s[i])
   }
   for k := 1; ; k <<= 1 {
       cmp := func(i, j int) bool {
           if rk[i] != rk[j] {
               return rk[i] < rk[j]
           }
           ri := -1
           rj := -1
           if i+k < n {
               ri = rk[i+k]
           }
           if j+k < n {
               rj = rk[j+k]
           }
           return ri < rj
       }
       sort.Slice(sa, func(a, b int) bool { return cmp(sa[a], sa[b]) })
       tmp[sa[0]] = 0
       for i := 1; i < n; i++ {
           tmp[sa[i]] = tmp[sa[i-1]]
           if cmp(sa[i-1], sa[i]) {
               tmp[sa[i]]++
           }
       }
       for i := 0; i < n; i++ {
           rk[i] = tmp[i]
       }
       if rk[sa[n-1]] == n-1 {
           break
       }
   }
   return sa, rk
}

package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Rect struct {
   a, b, c, d int64
}

func solve(rs []Rect) bool {
   n := len(rs)
   if n <= 1 {
       return true
   }
   // try vertical splits
   vs := make([]Rect, n)
   copy(vs, rs)
   sort.Slice(vs, func(i, j int) bool { return vs[i].a < vs[j].a })
   prefixMax := make([]int64, n)
   prefixMax[0] = vs[0].c
   for i := 1; i < n; i++ {
       if vs[i].c > prefixMax[i-1] {
           prefixMax[i] = vs[i].c
       } else {
           prefixMax[i] = prefixMax[i-1]
       }
   }
   suffixMin := make([]int64, n)
   suffixMin[n-1] = vs[n-1].a
   for i := n - 2; i >= 0; i-- {
       if vs[i].a < suffixMin[i+1] {
           suffixMin[i] = vs[i].a
       } else {
           suffixMin[i] = suffixMin[i+1]
       }
   }
   for i := 0; i < n-1; i++ {
       if prefixMax[i] <= suffixMin[i+1] {
           if solve(vs[:i+1]) && solve(vs[i+1:]) {
               return true
           }
       }
   }
   // try horizontal splits
   hs := make([]Rect, n)
   copy(hs, rs)
   sort.Slice(hs, func(i, j int) bool { return hs[i].b < hs[j].b })
   prefixMax = make([]int64, n)
   prefixMax[0] = hs[0].d
   for i := 1; i < n; i++ {
       if hs[i].d > prefixMax[i-1] {
           prefixMax[i] = hs[i].d
       } else {
           prefixMax[i] = prefixMax[i-1]
       }
   }
   suffixMin = make([]int64, n)
   suffixMin[n-1] = hs[n-1].b
   for i := n - 2; i >= 0; i-- {
       if hs[i].b < suffixMin[i+1] {
           suffixMin[i] = hs[i].b
       } else {
           suffixMin[i] = suffixMin[i+1]
       }
   }
   for i := 0; i < n-1; i++ {
       if prefixMax[i] <= suffixMin[i+1] {
           if solve(hs[:i+1]) && solve(hs[i+1:]) {
               return true
           }
       }
   }
   return false
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   rs := make([]Rect, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &rs[i].a, &rs[i].b, &rs[i].c, &rs[i].d)
   }
   if solve(rs) {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}

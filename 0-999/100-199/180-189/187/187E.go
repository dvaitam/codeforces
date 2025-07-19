package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func solve(an []int, a []int, n, l, s int) uint64 {
   var res uint64 = uint64(a[s] + a[n] - 2*a[1])
   pos := 0
   if l < s {
       for i := 1; i <= l; i++ {
           pos++
           an[pos] = s - i
       }
       t := an[pos]
       an[pos] = 1
       for i := 2; i <= t; i++ {
           pos++
           an[pos] = i
       }
       for i := s + 1; i <= n; i++ {
           pos++
           an[pos] = i
       }
   } else {
       // prepare gaps and order
       gap := make([]int, n+2)
       m := n - s
       p := make([]int, m)
       u := make([]bool, n+2)
       // compute gaps and initial p
       ss := s + 1
       for i := ss; i <= n; i++ {
           gap[i] = a[i] - a[i-1]
           p[i-ss] = i
       }
       // sort positions by gap
       sort.Slice(p, func(i, j int) bool {
           return gap[p[i]] < gap[p[j]]
       })
       rem := l - s + 1
       pw := rem - 1
       var cur uint64
       // take smallest rem gaps twice
       for i := 0; i < rem; i++ {
           idx := p[i]
           cur += uint64(gap[idx] * 2)
           u[idx] = true
       }
       mi := cur
       mf := n + 1
       // sliding to consider taking some suffix
       for i := 0; i < rem; i++ {
           at := n - i
           if !u[at] {
               for pw >= 0 && !u[p[pw]] {
                   pw--
               }
               cur -= uint64(gap[p[pw]] * 2)
               u[p[pw]] = false
               pw--
               cur += uint64(gap[at])
           } else {
               u[at] = false
               cur -= uint64(gap[at])
           }
           if cur < mi {
               mi = cur
               mf = at
           }
       }
       res += mi
       // reconstruct selection
       for i := ss; i <= n; i++ {
           u[i] = false
       }
       edCount := n - mf + 1
       for i := mf; i <= n; i++ {
           u[i] = true
       }
       // ensure we have rem marks
       cnt := edCount
       for i := 0; cnt < rem && i < len(p); i++ {
           if !u[p[i]] {
               u[p[i]] = true
               cnt++
           }
       }
       // build order
       for i := s - 2; i >= 1; i-- {
           pos++
           an[pos] = i
       }
       i := s
       for i <= n {
           ed := i
           for ed+1 <= n && u[ed+1] {
               ed++
           }
           for j := ed; j >= i; j-- {
               pos++
               an[pos] = j
           }
           i = ed + 1
       }
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, l, s int
   fmt.Fscan(reader, &n, &l, &s)
   a := make([]int, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   r := n - l - 1
   // edge cases
   if l == 0 {
       if s == 1 {
           fmt.Fprintln(writer, a[n])
           for i := 2; i <= n; i++ {
               if i > 2 {
                   fmt.Fprint(writer, " ")
               }
               fmt.Fprint(writer, i)
           }
           fmt.Fprintln(writer)
       } else {
           fmt.Fprintln(writer, -1)
       }
       return
   }
   if r == 0 {
       if s == n {
           fmt.Fprintln(writer, a[n])
           for i := n - 1; i >= 1; i-- {
               if i < n-1 {
                   fmt.Fprint(writer, " ")
               }
               fmt.Fprint(writer, i)
           }
           fmt.Fprintln(writer)
       } else {
           fmt.Fprintln(writer, -1)
       }
       return
   }
   a1 := make([]int, n+2)
   a2 := make([]int, n+2)
   ans1 := solve(a1, a, n, l, s)
   // prepare reversed negative array
   b := make([]int, n+2)
   for i := 1; i <= n; i++ {
       b[i] = -a[n+1-i]
   }
   ans2 := solve(a2, b, n, r, n+1-s)
   // choose best
   if ans1 < ans2 {
       fmt.Fprintln(writer, ans1)
       for i := 1; i < n; i++ {
           if i > 1 {
               fmt.Fprint(writer, " ")
           }
           fmt.Fprint(writer, a1[i])
       }
       fmt.Fprintln(writer)
   } else {
       fmt.Fprintln(writer, ans2)
       for i := 1; i < n; i++ {
           if i > 1 {
               fmt.Fprint(writer, " ")
           }
           fmt.Fprint(writer, n-a2[i]+1)
       }
       fmt.Fprintln(writer)
   }
}

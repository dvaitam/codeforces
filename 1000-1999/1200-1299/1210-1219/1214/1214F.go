package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Ele struct {
   val int
   id  int
}

func lowerBound(a []Ele, x int) int {
   return sort.Search(len(a), func(i int) bool { return a[i].val >= x })
}

func upperBound(a []Ele, x int) int {
   return sort.Search(len(a), func(i int) bool { return a[i].val > x })
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var m int64
   var n int
   fmt.Fscan(in, &m, &n)
   a := make([]Ele, n)
   b := make([]Ele, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &b[i].val)
       b[i].id = i
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i].val)
       a[i].id = i
   }
   sort.Slice(a, func(i, j int) bool { return a[i].val < a[j].val })
   sort.Slice(b, func(i, j int) bool { return b[i].val < b[j].val })
   INF := int64(1<<63 - 1)
   s1 := make([]int64, n+1)
   s2 := make([]int64, n+1)
   ans1 := make([]int, n)
   ans2 := make([]int, n)
   var Ans1, Ans2 int64 = INF, INF
   var pos1, pos2 int
   // compute s1
   for i := 1; i <= n; i++ {
       v := a[i-1].val
       if v <= b[i-1].val {
           s1[0] -= int64(v)
           s1[n-i+1] += int64(v)
       } else if v > b[n-1].val {
           s1[0] += int64(v)
           s1[n-i+1] -= int64(v)
       } else {
           cnt := lowerBound(b, v)
           s1[0] += int64(v)
           idx := cnt - i + 1
           s1[idx] -= int64(v)
           s1[idx] -= int64(v)
           s1[n-i+1] += int64(v)
       }
       s1[n-i+1] -= int64(v)
       s1[n-i+1] += m
   }
   for i := 1; i <= n; i++ {
       v := b[i-1].val
       if v >= a[i-1].val {
           s1[0] += int64(v)
           s1[i] -= int64(v)
       } else if v < a[0].val {
           s1[0] -= int64(v)
           s1[i] += int64(v)
       } else {
           pos := upperBound(a, v)
           idx := i - pos + 1
           s1[0] -= int64(v)
           s1[idx] += int64(v)
           s1[idx] += int64(v)
           s1[i] -= int64(v)
       }
       s1[i] += int64(v)
   }
   for i := 0; i <= n; i++ {
       if i > 0 {
           s1[i] += s1[i-1]
       }
       if s1[i] < Ans1 {
           Ans1 = s1[i]
           pos1 = i
       }
   }
   // compute s2
   for i := 1; i <= n; i++ {
       v := a[i-1].val
       if v >= b[i-1].val {
           s2[0] += int64(v)
           s2[i] -= int64(v)
       } else if v < b[0].val {
           s2[0] -= int64(v)
           s2[i] += int64(v)
       } else {
           pos := upperBound(b, v)
           idx := i - pos + 1
           s2[0] -= int64(v)
           s2[idx] += int64(v)
           s2[idx] += int64(v)
           s2[i] -= int64(v)
       }
       s2[i] += int64(v)
       s2[i] += m
   }
   for i := 1; i <= n; i++ {
       v := b[i-1].val
       if v <= a[i-1].val {
           s2[0] -= int64(v)
           s2[n-i+1] += int64(v)
       } else if v > a[n-1].val {
           s2[0] += int64(v)
           s2[n-i+1] -= int64(v)
       } else {
           cnt := lowerBound(a, v)
           idx := cnt - i + 1
           s2[0] += int64(v)
           s2[idx] -= int64(v)
           s2[idx] -= int64(v)
           s2[n-i+1] += int64(v)
       }
       s2[n-i+1] -= int64(v)
   }
   for i := 0; i <= n; i++ {
       if i > 0 {
           s2[i] += s2[i-1]
       }
       if s2[i] < Ans2 {
           Ans2 = s2[i]
           pos2 = i
       }
   }
   // build answers
   for i := 1; i <= n; i++ {
       j := (i-1+pos1)%n
       ans1[b[j].id] = a[i-1].id + 1
   }
   for i := 1; i <= n; i++ {
       j := (i-1-pos2+n)%n
       ans2[b[j].id] = a[i-1].id + 1
   }
   // output
   if Ans1 < Ans2 {
       fmt.Fprintln(out, Ans1)
       for i, v := range ans1 {
           if i > 0 {
               out.WriteByte(' ')
           }
           fmt.Fprint(out, v)
       }
       out.WriteByte('\n')
   } else {
       fmt.Fprintln(out, Ans2)
       for i, v := range ans2 {
           if i > 0 {
               out.WriteByte(' ')
           }
           fmt.Fprint(out, v)
       }
       out.WriteByte('\n')
   }
}

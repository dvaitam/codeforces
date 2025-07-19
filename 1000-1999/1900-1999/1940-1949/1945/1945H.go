package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const inf = int(0x7fffffff)

var rdr = bufio.NewReader(os.Stdin)
var wtr = bufio.NewWriter(os.Stdout)

func main() {
   defer wtr.Flush()
   var t int
   if _, err := fmt.Fscan(rdr, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n, m int
       fmt.Fscan(rdr, &n, &m)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(rdr, &a[i])
       }
       sort.Ints(a)
       found, v1, v2 := findSplit(a, n, m)
       if !found {
           fmt.Fprintln(wtr, "NO")
       } else {
           fmt.Fprintln(wtr, "YES")
           fmt.Fprintf(wtr, "%d ", len(v1))
           for _, v := range v1 {
               fmt.Fprintf(wtr, "%d ", v)
           }
           fmt.Fprintln(wtr)
           fmt.Fprintf(wtr, "%d ", len(v2))
           for _, v := range v2 {
               fmt.Fprintf(wtr, "%d ", v)
           }
           fmt.Fprintln(wtr)
       }
   }
}

func findSplit(a []int, n, m int) (bool, []int, []int) {
   li := a[n-1]
   pos := make([][2]int, li+1)
   for i := range pos {
       pos[i][0], pos[i][1] = -1, -1
   }
   for i, v := range a {
       if pos[v][0] < 0 {
           pos[v][0] = i
       } else if pos[v][1] < 0 {
           pos[v][1] = i
       }
   }
   tot := 100
   for i := li; i >= 1; i-- {
       var q []int
       for j := i; j <= li; j += i {
           if pos[j][0] >= 0 {
               q = append(q, pos[j][0])
           }
           if pos[j][1] >= 0 {
               q = append(q, pos[j][1])
           }
       }
       L := len(q)
       for xi := 0; xi < L; xi++ {
           for yi := xi + 1; yi < L; yi++ {
               x, y := q[xi], q[yi]
               if x < y {
                   tot--
                   k := inf
                   for j := 0; j < n; j++ {
                       if j != x && j != y {
                           k &= a[j]
                       }
                   }
                   k += m
                   if i > k {
                       v1 := []int{a[x], a[y]}
                       v2 := make([]int, 0, n-2)
                       for j := 0; j < n; j++ {
                           if j != x && j != y {
                               v2 = append(v2, a[j])
                           }
                       }
                       return true, v1, v2
                   }
                   if tot <= 0 {
                       return false, nil, nil
                   }
               }
           }
       }
   }
   return false, nil, nil
}

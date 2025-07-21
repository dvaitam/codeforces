package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var l1, r1 int
   fmt.Fscan(reader, &l1, &r1)
   // collect intersections of other segments with Alexey's segment
   intervals := make([][2]int, 0, n-1)
   for i := 1; i < n; i++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       // intersection [start, end]
       if l < r1 && r > l1 {
           start := l
           if start < l1 {
               start = l1
           }
           end := r
           if end > r1 {
               end = r1
           }
           if start < end {
               intervals = append(intervals, [2]int{start, end})
           }
       }
   }
   // if no overlaps, full length is available
   if len(intervals) == 0 {
       fmt.Println(r1 - l1)
       return
   }
   // sort by start
   sort.Slice(intervals, func(i, j int) bool {
       if intervals[i][0] != intervals[j][0] {
           return intervals[i][0] < intervals[j][0]
       }
       return intervals[i][1] < intervals[j][1]
   })
   // merge and sum covered length
   covered := 0
   curS, curE := intervals[0][0], intervals[0][1]
   for _, seg := range intervals[1:] {
       s, e := seg[0], seg[1]
       if s <= curE {
           if e > curE {
               curE = e
           }
       } else {
           covered += curE - curS
           curS, curE = s, e
       }
   }
   covered += curE - curS
   // result is Alexey's length minus covered
   result := (r1 - l1) - covered
   if result < 0 {
       result = 0
   }
   fmt.Println(result)

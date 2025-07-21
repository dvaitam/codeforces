package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func LIS(a []int) int {
   d := make([]int, 0, len(a))
   for _, v := range a {
       i := sort.SearchInts(d, v)
       if i == len(d) {
           d = append(d, v)
       } else {
           d[i] = v
       }
   }
   return len(d)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var k, n, maxb int
   var t int64
   if _, err := fmt.Fscan(in, &k, &n, &maxb, &t); err != nil {
       return
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // threshold for simulation: max total elements ~2e7
   maxElems := 20000000
   for vi := 0; vi < k; vi++ {
       b := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &b[i])
       }
       // count distinct
       uniq := make(map[int]struct{}, n)
       for _, v := range b {
           uniq[v] = struct{}{}
       }
       dcount := len(uniq)
       var ans int
       if t == 1 {
           ans = LIS(b)
       } else if int64(n)*t <= int64(maxElems) {
           // simulate full
           a := make([]int, 0, int64(n)*t)
           for i := int64(0); i < t; i++ {
               a = append(a, b...)
           }
           ans = LIS(a)
       } else if t >= int64(dcount) {
           ans = dcount
       } else {
           // partial simulate up to dcount blocks
           sim := t
           if sim > int64(dcount) {
               sim = int64(dcount)
           }
           if int64(n)*sim > int64(maxElems) {
               sim = int64(maxElems) / int64(n)
               if sim < 1 {
                   sim = 1
               }
           }
           a := make([]int, 0, int64(n)*sim)
           for i := int64(0); i < sim; i++ {
               a = append(a, b...)
           }
           ans = LIS(a)
       }
       fmt.Fprintln(writer, ans)
   }
}

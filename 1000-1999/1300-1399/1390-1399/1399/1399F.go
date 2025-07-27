package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   fmt.Fscan(reader, &t)
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(reader, &n)
       segs := make([][2]int, n)
       xs := make([]int, 0, 2*n)
       for i := 0; i < n; i++ {
           l, r := 0, 0
           fmt.Fscan(reader, &l, &r)
           segs[i][0], segs[i][1] = l, r
           xs = append(xs, l, r)
       }
       sort.Ints(xs)
       m := 0
       last := -1
       for _, v := range xs {
           if v != last {
               xs[m] = v
               m++
               last = v
           }
       }
       xs = xs[:m]
       // map original to index 1..m
       idx := make(map[int]int, m)
       for i, v := range xs {
           idx[v] = i + 1
       }
       endsAt := make([][]int, m+2)
       for _, s := range segs {
           l := idx[s[0]]
           r := idx[s[1]]
           endsAt[r] = append(endsAt[r], l)
       }
       // DP f[i][j]: max arcs in [i..j]
       // use int16 to save memory
       f := make([][]int16, m+2)
       for i := 0; i <= m+1; i++ {
           f[i] = make([]int16, m+2)
       }
       for width := 0; width < m; width++ {
           for i := 1; i+width <= m; i++ {
               j := i + width
               // skip j<i
               var best int16
               if j > i {
                   best = f[i][j-1]
               }
               // consider any segment [l0,j]
               for _, l0 := range endsAt[j] {
                   if l0 < i {
                       continue
                   }
                   // f[i][l0-1] + 1 + f[l0+1][j-1]
                   cur := int16(1)
                   if l0 > i {
                       cur += f[i][l0-1]
                   }
                   if l0 < j {
                       cur += f[l0+1][j-1]
                   }
                   if cur > best {
                       best = cur
                   }
               }
               f[i][j] = best
           }
       }
       // result for full range
       if m > 0 {
           fmt.Fprintln(writer, f[1][m])
       } else {
           fmt.Fprintln(writer, 0)
       }
   }
}

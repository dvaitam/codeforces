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
   _, err := fmt.Fscan(reader, &n)
   if err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   if n == 1 {
       // Any x will form an AP of length 2
       fmt.Println(-1)
       return
   }
   sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
   if n == 2 {
       x, y := a[0], a[1]
       d := y - x
       if d == 0 {
           fmt.Println(1)
           fmt.Println(x)
       } else if d%2 == 0 {
           mid := x + d/2
           res := []int64{x - d, mid, y + d}
           sort.Slice(res, func(i, j int) bool { return res[i] < res[j] })
           fmt.Println(len(res))
           for i, v := range res {
               if i > 0 {
                   fmt.Print(" ")
               }
               fmt.Print(v)
           }
           fmt.Println()
       } else {
           res := []int64{x - d, y + d}
           fmt.Println(2)
           fmt.Println(res[0], res[1])
       }
       return
   }
   // n >= 3
   diffs := make([]int64, n-1)
   for i := 1; i < n; i++ {
       diffs[i-1] = a[i] - a[i-1]
   }
   // count distinct diffs
   // find min diff and track larger diffs
   dmin := diffs[0]
   for _, d := range diffs {
       if d < dmin {
           dmin = d
       }
   }
   cntMin, cntOther := 0, 0
   var otherD int64
   pos := -1
   for i, d := range diffs {
       if d == dmin {
           cntMin++
       } else {
           cntOther++
           otherD = d
           pos = i
       }
   }
   if cntOther == 0 {
       // already AP
       if dmin == 0 {
           fmt.Println(1)
           fmt.Println(a[0])
       } else {
           res := []int64{a[0] - dmin, a[n-1] + dmin}
           fmt.Println(2)
           fmt.Println(res[0], res[1])
       }
       return
   }
   // one irregular gap
   if cntOther == 1 && otherD == 2*dmin {
       // insert in middle
       x := a[pos] + dmin
       fmt.Println(1)
       fmt.Println(x)
       return
   }
   fmt.Println(0)
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   fmt.Fscan(in, &n, &m)
   // sum per diameter
   const maxD = 1000
   sumM := make([]int, maxD+1)
   sumC := make([]int, maxD+1)
   // color counts per diameter
   md := make([]map[int]int, maxD+1)
   cd := make([]map[int]int, maxD+1)
   for d := 0; d <= maxD; d++ {
       md[d] = make(map[int]int)
       cd[d] = make(map[int]int)
   }
   // read markers
   for i := 0; i < n; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       sumM[y]++
       md[y][x]++
   }
   // read caps
   for j := 0; j < m; j++ {
       var a, b int
       fmt.Fscan(in, &a, &b)
       sumC[b]++
       cd[b][a]++
   }
   // compute results
   var total, beautifulTotal int
   for d := 0; d <= maxD; d++ {
       if sumM[d] == 0 || sumC[d] == 0 {
           continue
       }
       // total matches for this diameter
       t := sumM[d]
       if sumC[d] < t {
           t = sumC[d]
       }
       total += t
       // beautiful matches: sum of min counts per color
       var b int
       for color, cntM := range md[d] {
           cntC := cd[d][color]
           if cntC < cntM {
               b += cntC
           } else {
               b += cntM
           }
       }
       // ensure not exceed total matches
       if b > t {
           b = t
       }
       beautifulTotal += b
   }
   fmt.Fprintln(out, total, beautifulTotal)
}

package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   fmt.Fscan(in, &t)
   for t > 0 {
       t--
       var a, b int
       fmt.Fscan(in, &a, &b)
       r := make([]int, b)
       for i := 0; i < b; i++ {
           fmt.Fscan(in, &r[i])
       }
       xa := make([]float64, a)
       ya := make([]float64, a)
       for i := 0; i < a; i++ {
           fmt.Fscan(in, &xa[i], &ya[i])
       }
       xb := make([]float64, b)
       yb := make([]float64, b)
       for i := 0; i < b; i++ {
           fmt.Fscan(in, &xb[i], &yb[i])
       }
       fmt.Fprintln(out, "YES")
       // connect each B_j to A1
       for j := 0; j < b; j++ {
           fmt.Fprintf(out, "%d %d\n", j+1, 1)
       }
       // prepare Bs repeated r_j-1 times with angle
       type pair struct{ j int; ang float64 }
       Bs := make([]pair, 0)
       for j := 0; j < b; j++ {
           for k := 0; k < r[j]-1; k++ {
               ang := math.Atan2(yb[j]-ya[0], xb[j]-xa[0])
               Bs = append(Bs, pair{j + 1, ang})
           }
       }
       // prepare As 2..a with angle
       type paar struct{ i int; ang float64 }
       As := make([]paar, 0)
       for i := 1; i < a; i++ {
           ang := math.Atan2(ya[i]-ya[0], xa[i]-xa[0])
           As = append(As, paar{i + 1, ang})
       }
       // sort by angle
       sort.Slice(Bs, func(i, j int) bool { return Bs[i].ang < Bs[j].ang })
       sort.Slice(As, func(i, j int) bool { return As[i].ang < As[j].ang })
       // pair them
       for k := 0; k < len(As) && k < len(Bs); k++ {
           fmt.Fprintf(out, "%d %d\n", Bs[k].j, As[k].i)
       }
   }
}

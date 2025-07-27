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

   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n, m int
       fmt.Fscan(in, &n, &m)
       // read matrix
       a := make([][]int, n)
       for i := 0; i < n; i++ {
           a[i] = make([]int, n)
           for j := 0; j < n; j++ {
               fmt.Fscan(in, &a[i][j])
           }
       }
       var ops string
       fmt.Fscan(in, &ops)
       // p[k]: which original dimension goes to final dim k
       p := [3]int{0, 1, 2}
       // sh[k]: shift for final dim k
       sh := [3]int{0, 0, 0}
       for _, c := range ops {
           switch c {
           case 'R':
               sh[1]++
           case 'L':
               sh[1]--
           case 'D':
               sh[0]++
           case 'U':
               sh[0]--
           case 'I':
               p[1], p[2] = p[2], p[1]
               sh[1], sh[2] = sh[2], sh[1]
           case 'C':
               p[0], p[2] = p[2], p[0]
               sh[0], sh[2] = sh[2], sh[0]
           }
       }
       // normalize shifts
       for k := 0; k < 3; k++ {
           sh[k] %= n
           if sh[k] < 0 {
               sh[k] += n
           }
       }
       // build result
       res := make([][]int, n)
       for i := range res {
           res[i] = make([]int, n)
       }
       for i := 0; i < n; i++ {
           for j := 0; j < n; j++ {
               // original dims: 0->i,1->j,2->a[i][j]-1
               v0 := i
               v1 := j
               v2 := a[i][j] - 1
               d := [3]int{v0, v1, v2}
               var f [3]int
               for k := 0; k < 3; k++ {
                   f[k] = (d[p[k]] + sh[k]) % n
               }
               res[f[0]][f[1]] = f[2] + 1
           }
       }
       // output result
       for i := 0; i < n; i++ {
           for j := 0; j < n; j++ {
               if j > 0 {
                   out.WriteByte(' ')
               }
               fmt.Fprint(out, res[i][j])
           }
           out.WriteByte('\n')
       }
   }
}

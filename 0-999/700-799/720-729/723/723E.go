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
   fmt.Fscan(in, &t)
   for ; t > 0; t-- {
       var n, m int
       fmt.Fscan(in, &n, &m)
       // adjacency matrix
       v := make([][]int, n)
       for i := range v {
           v[i] = make([]int, n)
       }
       vv := make([]int, n)
       vvv := make([]int, n)
       // read edges
       for i := 0; i < m; i++ {
           var a, b int
           fmt.Fscan(in, &a, &b)
           a--
           b--
           v[a][b] = 1
           v[b][a] = 1
       }
       // compute initial parity
       cp := 0
       for i := 0; i < n; i++ {
           sum := 0
           for j := 0; j < n; j++ {
               sum += v[i][j]
           }
           vvv[i] = sum
           vv[i] = 1 - (sum % 2)
           cp += vv[i]
       }
       // output count of vertices with even degree
       fmt.Fprintln(out, cp)
       // orient edges
       for i := 0; i < n; i++ {
           if vv[i] == 0 {
               continue
           }
           jFlag := 0
           for vvv[i] > 0 {
               a := i
               for vvv[a] > 0 && vv[a] > 0 {
                   for k := 0; k < n; k++ {
                       if v[a][k] == 1 {
                           if jFlag == 0 {
                               v[a][k] = 2
                               v[k][a] = 0
                           } else {
                               v[a][k] = 0
                               v[k][a] = 2
                           }
                           vvv[a]--
                           vvv[k]--
                           a = k
                           break
                       }
                   }
               }
               if vvv[i]%2 == 1 {
                   jFlag = 1
               } else {
                   jFlag = 0
               }
           }
       }
       // print oriented edges
       for i := 0; i < n; i++ {
           for j := 0; j < n; j++ {
               if v[i][j] == 2 {
                   fmt.Fprintln(out, i+1, j+1)
               }
               if v[i][j] == 1 {
                   fmt.Fprintln(out, i+1, j+1)
                   v[j][i] = 0
               }
           }
       }
   }
}

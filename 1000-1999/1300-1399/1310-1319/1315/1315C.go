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
       var n int
       fmt.Fscan(in, &n)
       b := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &b[i])
       }
       // mark used b's
       used := make([]bool, 2*n+1)
       for _, v := range b {
           if v >= 1 && v <= 2*n {
               used[v] = true
           }
       }
       // build remaining numbers
       rem := make([]int, 0, 2*n)
       for v := 1; v <= 2*n; v++ {
           if !used[v] {
               rem = append(rem, v)
           }
       }
       a := make([]int, 2*n)
       ok := true
       // for each pair assign b[i] and smallest rem > b[i]
       for i := 0; i < n; i++ {
           a[2*i] = b[i]
           // find smallest rem[j] > b[i]
           found := -1
           for j, rv := range rem {
               if rv > b[i] {
                   found = j
                   break
               }
           }
           if found == -1 {
               ok = false
               break
           }
           a[2*i+1] = rem[found]
           // remove rem[found]
           rem = append(rem[:found], rem[found+1:]...)
       }
       if !ok {
           fmt.Fprintln(out, -1)
       } else {
           for i, v := range a {
               if i > 0 {
                   out.WriteByte(' ')
               }
               fmt.Fprint(out, v)
           }
           out.WriteByte('\n')
       }
   }
}

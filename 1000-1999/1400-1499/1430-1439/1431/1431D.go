package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       used := make([]bool, n)
       c := 0
       res := make([]int, 0, n)
       for k := 0; k < n; k++ {
           found := -1
           for i := 0; i < n; i++ {
               if !used[i] && a[i] <= c {
                   found = i
                   break
               }
           }
           if found != -1 {
               used[found] = true
               res = append(res, found+1)
               c = 1
           } else {
               best := -1
               bi := -1
               for i := 0; i < n; i++ {
                   if !used[i] && a[i] > best {
                       best = a[i]
                       bi = i
                   }
               }
               used[bi] = true
               res = append(res, bi+1)
               c++
           }
       }
       for i, v := range res {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
       writer.WriteByte('\n')
   }
}

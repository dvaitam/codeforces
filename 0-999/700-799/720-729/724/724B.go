package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   a := make([][]int, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(in, &a[i][j])
       }
   }
   // Try all possible global column swaps (including no swap when c1 == c2)
   for c1 := 0; c1 < m; c1++ {
       for c2 := c1; c2 < m; c2++ {
           ok := true
           for i := 0; i < n && ok; i++ {
               // build row after possible global swap
               b := make([]int, m)
               for j := 0; j < m; j++ {
                   if j == c1 {
                       b[j] = a[i][c2]
                   } else if j == c2 {
                       b[j] = a[i][c1]
                   } else {
                       b[j] = a[i][j]
                   }
               }
               // collect mismatches
               var pos [2]int
               cnt := 0
               for j := 0; j < m; j++ {
                   if b[j] != j+1 {
                       if cnt < 2 {
                           pos[cnt] = j
                       }
                       cnt++
                       if cnt > 2 {
                           break
                       }
                   }
               }
               if cnt == 0 {
                   continue
               }
               if cnt == 2 {
                   j1, j2 := pos[0], pos[1]
                   // check if swapping fixes
                   if !(b[j1] == j2+1 && b[j2] == j1+1) {
                       ok = false
                   }
                   continue
               }
               // cannot fix (1 mismatch or >2 mismatches)
               ok = false
           }
           if ok {
               fmt.Println("YES")
               return
           }
       }
   }
   fmt.Println("NO")
}

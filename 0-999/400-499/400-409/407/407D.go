package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([][]int, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(reader, &a[i][j])
       }
   }
   // value range up to 160000
   const maxVal = 160000
   lastRow := make([]int, maxVal+1)
   seen := make([]int, maxVal+1)
   epoch := 0
   best := 0
   for l := 0; l < m; l++ {
       epoch++
       top := 0
       for r := l; r < m; r++ {
           for i := 0; i < n; i++ {
               v := a[i][r]
               if seen[v] != epoch {
                   seen[v] = epoch
                   lastRow[v] = i
               } else {
                   prev := lastRow[v]
                   if prev >= top {
                       top = prev + 1
                   }
                   lastRow[v] = i
               }
               h := i - top + 1
               if h > 0 {
                   w := r - l + 1
                   area := h * w
                   if area > best {
                       best = area
                   }
               }
           }
       }
   }
   fmt.Println(best)
}

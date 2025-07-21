package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // Find spacing steps where k = n/step >= 3
   steps := make([]int, 0)
   for d := 1; d*d <= n; d++ {
       if n%d == 0 {
           if d <= n/3 {
               steps = append(steps, d)
           }
           d2 := n / d
           if d2 != d && d2 <= n/3 {
               steps = append(steps, d2)
           }
       }
   }
   // Check each step for a regular polygon of good knights
   for _, step := range steps {
       good := make([]bool, step)
       for i := range good {
           good[i] = true
       }
       for idx, v := range a {
           if v == 0 {
               m := idx % step
               if good[m] {
                   good[m] = false
               }
           }
       }
       for _, ok := range good {
           if ok {
               fmt.Println("YES")
               return
           }
       }
   }
   fmt.Println("NO")
}

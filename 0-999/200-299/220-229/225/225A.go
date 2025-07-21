package main

import (
   "bufio"
   "fmt"
   "os"
)

// orientation of a dice
type ori struct {
   top, bottom, front, back, right, left int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, x int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   fmt.Fscan(in, &x)
   a := make([]int, n)
   b := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i], &b[i])
   }
   // generate all 48 orientations
   oris := make([]ori, 0, 48)
   // faces values 1..6, opposites sum to 7
   for t := 1; t <= 6; t++ {
       bt := 7 - t
       for f := 1; f <= 6; f++ {
           if f == t || f == bt {
               continue
           }
           bk := 7 - f
           // remaining for left/right
           rem := make([]int, 0, 2)
           for v := 1; v <= 6; v++ {
               if v != t && v != bt && v != f && v != bk {
                   rem = append(rem, v)
               }
           }
           // two chiral possibilities
           oris = append(oris, ori{t, bt, f, bk, rem[0], rem[1]})
           oris = append(oris, ori{t, bt, f, bk, rem[1], rem[0]})
       }
   }
   m := len(oris)
   // dp: number of ways to reach orientation j at current level (cap at 2)
   dp := make([]int, m)
   // initial for i=0
   for j, o := range oris {
       if o.top == x && o.front == a[0] && o.right == b[0] {
           dp[j] = 1
       }
   }
   // iterate layers
   for i := 1; i < n; i++ {
       next := make([]int, m)
       for j, o := range oris {
           if o.front != a[i] || o.right != b[i] {
               continue
           }
           // consider all prev
           sum := 0
           for k, p := range oris {
               if dp[k] == 0 {
                   continue
               }
               // bottom of prev p vs top of curr o
               if p.bottom != o.top {
                   sum += dp[k]
                   if sum >= 2 {
                       sum = 2
                       break
                   }
               }
           }
           next[j] = sum
       }
       dp = next
   }
   total := 0
   for _, v := range dp {
       total += v
       if total >= 2 {
           break
       }
   }
   if total == 1 {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}

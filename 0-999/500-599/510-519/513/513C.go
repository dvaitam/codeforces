package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   l := make([]int, n)
   r := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &l[i], &r[i])
   }
   var ans, dem float64
   for tap := 1; tap < (1 << n); tap++ {
       L, R := -1000000000, 1000000000
       for i := 0; i < n; i++ {
           if tap&(1<<i) != 0 {
               L = max(L, l[i])
               R = min(R, r[i])
           }
       }
       if L > R {
           continue
       }
       for b := L; b <= R; b++ {
           // Case 1: one winner outside the set
           for win := 0; win < n; win++ {
               if tap&(1<<win) != 0 {
                   continue
               }
               if r[win] <= b {
                   continue
               }
               c := float64(r[win] - max(l[win]-1, b))
               for j := 0; j < n; j++ {
                   if j == win {
                       continue
                   }
                   if tap&(1<<j) != 0 {
                       continue
                   }
                   if l[j] >= b {
                       c = 0
                       break
                   }
                   c *= float64(min(r[j]+1, b) - l[j])
               }
               ans += c * float64(b)
               dem += c
           }
           // Case 2: max from the set (only if set size > 1)
           if bits.OnesCount(uint(tap)) == 1 {
               continue
           }
           c2 := float64(1)
           for j := 0; j < n; j++ {
               if tap&(1<<j) != 0 {
                   continue
               }
               if l[j] >= b {
                   c2 = 0
                   break
               }
               c2 *= float64(min(r[j]+1, b) - l[j])
           }
           ans += c2 * float64(b)
           dem += c2
       }
   }
   if dem == 0 {
       fmt.Println("0.0000000000")
   } else {
       fmt.Printf("%.10f\n", ans/dem)
   }
}

package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// node holds a normalized pair
type node struct {
   a, b int
}

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   if a < 0 {
       return -a
   }
   return a
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n)
   b := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &b[i])
   }
   c := make([]node, n)
   ans := 0
   maxi := 0
   flag := 0
   // normalize and count zeros
   for i := 0; i < n; i++ {
       ai, bi := a[i], b[i]
       if ai != 0 && bi != 0 {
           // normalize sign
           if ai < 0 && bi < 0 {
               g := gcd(-ai, -bi)
               ai = (-ai) / g
               bi = (-bi) / g
           } else if ai < 0 {
               g := gcd(-ai, bi)
               ai = ai / g
               bi = bi / g
           } else if bi < 0 {
               g := gcd(ai, -bi)
               ai = (-ai) / g
               bi = (-bi) / g
           } else {
               g := gcd(ai, bi)
               ai = ai / g
               bi = bi / g
           }
           c[i] = node{ai, bi}
       } else {
           // one or both zero
           if ai == 0 && bi == 0 {
               ans++
               c[i] = node{0, 0}
               continue
           }
           if ai == 0 && bi != 0 {
               flag++
           }
           if bi == 0 {
               maxi++
           }
           c[i] = node{ai, bi}
       }
   }
   // sort by b, then a
   sort.Slice(c, func(i, j int) bool {
       if c[i].b == c[j].b {
           return c[i].a < c[j].a
       }
       return c[i].b < c[j].b
   })
   if flag != (n - ans) {
       if maxi < 1 {
           maxi = 1
       }
   }
   sum := 0
   // count max group of equal non-zero slopes
   for i := 1; i < n; i++ {
       if c[i].a == 0 && c[i].b == 0 {
           continue
       }
       if c[i].b == 0 {
           continue
       }
       if c[i].a != 0 && c[i].a == c[i-1].a && c[i].b == c[i-1].b {
           sum++
       } else {
           if sum > 0 {
               if maxi < sum+1 {
                   maxi = sum + 1
               }
           }
           sum = 0
       }
   }
   if sum > 0 {
       if maxi < sum+1 {
           maxi = sum + 1
       }
   }
   // result is min(maxi+ans, n)
   res := maxi + ans
   if res > n {
       res = n
   }
   fmt.Println(res)
}

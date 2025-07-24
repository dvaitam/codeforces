package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   var total int64
   for i := 0; i < n; i++ {
       var a1, b1, a2, b2 int64
       fmt.Fscan(in, &a1, &b1, &a2, &b2)
       // DP memo: mask 0..3, turn 0=Alice,1=Bonnie, prev 0=false,1=true
       const INF = int64(9e18)
       var memo [4][2][2]int64
       var seen [4][2][2]bool
       var dfs func(mask int, turn int, prev int) int64
       dfs = func(mask int, turn int, prev int) int64 {
           if mask == 0 {
               return 0
           }
           if seen[mask][turn][prev] {
               return memo[mask][turn][prev]
           }
           seen[mask][turn][prev] = true
           var res int64
           if turn == 0 {
               // Alice: maximize
               res = -INF
               // pass
               if prev == 1 {
                   if res < 0 {
                       res = 0
                   }
               } else {
                   v := dfs(mask, 1, 1)
                   if v > res {
                       res = v
                   }
               }
               // pick top if available
               if mask&1 != 0 {
                   v := a1 + dfs(mask&^1, 1, 0)
                   if v > res {
                       res = v
                   }
               } else if mask&2 != 0 {
                   // only bottom available
                   v := a2 + dfs(mask&^2, 1, 0)
                   if v > res {
                       res = v
                   }
               }
           } else {
               // Bonnie: minimize
               res = INF
               // pass
               if prev == 1 {
                   if res > 0 {
                       res = 0
                   }
               } else {
                   v := dfs(mask, 0, 1)
                   if v < res {
                       res = v
                   }
               }
               // pick top if available
               if mask&1 != 0 {
                   v := -b1 + dfs(mask&^1, 0, 0)
                   if v < res {
                       res = v
                   }
               } else if mask&2 != 0 {
                   v := -b2 + dfs(mask&^2, 0, 0)
                   if v < res {
                       res = v
                   }
               }
           }
           memo[mask][turn][prev] = res
           return res
       }
       // initial: both photos present: mask=3, Alice turn, prev=0
       total += dfs(3, 0, 0)
   }
   // output total difference
   fmt.Println(total)
}

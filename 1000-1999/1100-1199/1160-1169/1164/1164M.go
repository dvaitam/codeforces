package main

import (
   "fmt"
)

var (
   used [10]bool
   curr [7]int
   found bool
   ans int
)

func dfs(pos int) {
   if found {
       return
   }
   if pos == 7 {
       // build number
       n := 0
       sum := 0
       for i := 0; i < 7; i++ {
           d := curr[i]
           n = n*10 + d
           sum += d
       }
       // check divisibility
       for i := 0; i < 7; i++ {
           d := curr[i]
           if d == 0 || n%d != 0 {
               return
           }
       }
       ans = sum
       found = true
       return
   }
   for d := 1; d <= 9; d++ {
       if !used[d] {
           used[d] = true
           curr[pos] = d
           dfs(pos + 1)
           used[d] = false
           if found {
               return
           }
       }
   }
}

func main() {
   dfs(0)
   fmt.Println(ans)
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, a, b int
   h       []int
   shoot   []int
)

// dfs attempts to use 'balls' shots to reduce all h[1..n] below 0
func dfs(l, balls, last int) bool {
   if balls == 0 {
       for i := 1; i <= n; i++ {
           if h[i] >= 0 {
               return false
           }
       }
       return true
   }
   if l <= n && h[l] < 0 {
       return dfs(l+1, balls, last)
   }
   // determine range of possible shots
   lb := last
   if l > 2 && l > lb {
       lb = l
   }
   if lb < 2 {
       lb = 2
   }
   ub := n - 1
   if l+1 < ub {
       ub = l + 1
   }
   for i := lb; i <= ub; i++ {
       shoot[balls] = i
       // apply shot
       h[i] -= a
       h[i-1] -= b
       h[i+1] -= b
       if dfs(l, balls-1, i) {
           return true
       }
       // undo shot
       h[i] += a
       h[i-1] += b
       h[i+1] += b
   }
   return false
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   if _, err := fmt.Fscan(reader, &n, &a, &b); err != nil {
       return
   }
   h = make([]int, n+3)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &h[i])
   }
   // allocate shoot positions, size up to n*15 approx
   shoot = make([]int, 1000)
   // try increasing number of shots
   for ans := 1; ; ans++ {
       if dfs(1, ans, 2) {
           // output result
           fmt.Println(ans)
           for i := ans; i >= 1; i-- {
               fmt.Printf("%d", shoot[i])
               if i > 1 {
                   fmt.Printf(" ")
               }
           }
           fmt.Println()
           return
       }
   }
}

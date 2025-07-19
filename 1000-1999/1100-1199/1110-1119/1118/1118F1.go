package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   colors := make([]int, n)
   var totalRed, totalBlue int
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &colors[i])
       if colors[i] == 1 {
           totalRed++
       } else if colors[i] == 2 {
           totalBlue++
       }
   }
   adj := make([][]int, n)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       u--
       v--
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }

   parent := make([]int, n)
   parent[0] = -1
   stack := make([]int, 0, n)
   stack = append(stack, 0)
   order := make([]int, 0, n)
   for len(stack) > 0 {
       u := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       order = append(order, u)
       for _, v := range adj[u] {
           if v == parent[u] {
               continue
           }
           parent[v] = u
           stack = append(stack, v)
       }
   }

   redCount := make([]int, n)
   blueCount := make([]int, n)
   var ans int
   // process in post-order
   for i := len(order) - 1; i >= 0; i-- {
       u := order[i]
       // accumulate counts from children
       for _, v := range adj[u] {
           if parent[v] == u {
               redCount[u] += redCount[v]
               blueCount[u] += blueCount[v]
           }
       }
       if colors[u] == 1 {
           redCount[u]++
       } else if colors[u] == 2 {
           blueCount[u]++
       }
       if u != 0 {
           // subtree at u vs rest
           if redCount[u] == 0 && blueCount[u] == totalBlue {
               ans++
           } else if blueCount[u] == 0 && redCount[u] == totalRed {
               ans++
           }
       }
   }
   fmt.Fprint(out, ans)
}

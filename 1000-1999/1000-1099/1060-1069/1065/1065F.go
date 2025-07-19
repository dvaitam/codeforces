package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = 1000000000

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   children := make([][]int, n+1)
   for i := 2; i <= n; i++ {
       var p int
       fmt.Fscan(in, &p)
       children[p] = append(children[p], i)
   }
   // BFS order
   order := make([]int, 0, n)
   queue := make([]int, 0, n)
   queue = append(queue, 1)
   for head := 0; head < len(queue); head++ {
       u := queue[head]
       order = append(order, u)
       for _, v := range children[u] {
           queue = append(queue, v)
       }
   }
   // allocate arrays
   h := make([]int, n+1)
   a := make([]int, n+1)
   b := make([]int, n+1)
   // process in reverse BFS (bottom-up)
   for idx := len(order) - 1; idx >= 0; idx-- {
       u := order[idx]
       if len(children[u]) == 0 {
           // leaf
           h[u] = 0
           a[u] = 1
           b[u] = 1
       } else {
           // internal
           hu := INF
           au := 0
           bu := 0
           for _, v := range children[u] {
               hu = min(hu, h[v]+1)
               au += a[v]
               bu = max(bu, b[v]-a[v])
           }
           bu += au
           h[u] = hu
           a[u] = au
           b[u] = bu
           if h[u] >= k {
               a[u] = 0
           }
       }
   }
   // result at root
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, b[1])
}

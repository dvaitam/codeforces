package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewScanner(os.Stdin)
   in.Split(bufio.ScanWords)
   var readInt = func() int {
       in.Scan()
       v, _ := strconv.Atoi(in.Text())
       return v
   }
   n := readInt()
   children := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       u := readInt()
       v := readInt()
       // edge from u to v
       children[u] = append(children[u], v)
   }
   // compute depth parity by BFS
   depth := make([]int, n+1)
   queue := make([]int, 0, n)
   queue = append(queue, 1)
   for i := 0; i < len(queue); i++ {
       u := queue[i]
       for _, v := range children[u] {
           depth[v] = depth[u] + 1
           queue = append(queue, v)
       }
   }
   // postorder nodes
   stack := make([]int, 0, n)
   stack = append(stack, 1)
   order := make([]int, 0, n)
   for len(stack) > 0 {
       u := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       order = append(order, u)
       for _, v := range children[u] {
           stack = append(stack, v)
       }
   }
   // prepare dp arrays
   need1 := make([]int, n+1)
   need2 := make([]int, n+1)
   // count leaves
   m := 0
   // traverse in postorder
   for i := len(order) - 1; i >= 0; i-- {
       u := order[i]
       if len(children[u]) == 0 {
           need1[u] = 1
           need2[u] = 1
           m++
       } else if depth[u]%2 == 0 {
           // Max node
           mn1 := n + 5
           sum2 := 0
           for _, v := range children[u] {
               if need1[v] < mn1 {
                   mn1 = need1[v]
               }
               sum2 += need2[v]
           }
           need1[u] = mn1
           need2[u] = sum2
       } else {
           // Min node
           sum1 := 0
           mn2 := n + 5
           for _, v := range children[u] {
               sum1 += need1[v]
               if need2[v] < mn2 {
                   mn2 = need2[v]
               }
           }
           need1[u] = sum1
           need2[u] = mn2
       }
   }
   // result for Shambambukli and Mazukta
   maxRes := m - need1[1] + 1
   minRes := need2[1]
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintf(w, "%d %d\n", maxRes, minRes)
}

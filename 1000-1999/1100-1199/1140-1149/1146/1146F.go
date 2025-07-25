package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 998244353

func add(a, b int) int {
   a += b
   if a >= mod {
       a -= mod
   }
   return a
}

func mul(a, b int) int {
   return int((int64(a) * int64(b)) % mod)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // build tree
   children := make([][]int, n+1)
   for i := 2; i <= n; i++ {
       var p int
       fmt.Fscan(reader, &p)
       children[p] = append(children[p], i)
   }
   // prepare arrays
   leafCount := make([]int, n+1)
   dp := make([]int, n+1)
   // post-order traversal using stack
   type frame struct{ u, state int }
   stack := make([]frame, 0, n*2)
   stack = append(stack, frame{1, 0})
   for len(stack) > 0 {
       fr := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       u, state := fr.u, fr.state
       if state == 0 {
           // enter
           stack = append(stack, frame{u, 1})
           for _, v := range children[u] {
               stack = append(stack, frame{v, 0})
           }
       } else {
           // exit: compute for u
           if len(children[u]) == 0 {
               // leaf
               leafCount[u] = 1
               dp[u] = 1
           } else {
               totalLeaves := 0
               prod := 1
               for _, v := range children[u] {
                   totalLeaves += leafCount[v]
                   if totalLeaves >= mod {
                       // totalLeaves not used mod
                   }
                   prod = mul(prod, dp[v])
               }
               leafCount[u] = totalLeaves
               if totalLeaves >= 2 {
                   dp[u] = add(prod, 1)
               } else {
                   dp[u] = prod
               }
           }
       }
   }
   fmt.Fprintln(writer, dp[1])
}

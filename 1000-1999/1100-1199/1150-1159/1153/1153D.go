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
   op := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &op[i])
   }
   children := make([][]int, n+1)
   for v := 2; v <= n; v++ {
       var p int
       fmt.Fscan(in, &p)
       children[p] = append(children[p], v)
   }
   // f[v]: DP value as described
   f := make([]int, n+1)
   // Count leaves and compute f bottom-up: nodes have parent < child, so reverse order
   leafCount := 0
   for v := n; v >= 1; v-- {
       if len(children[v]) == 0 {
           f[v] = 1
           leafCount++
       } else if op[v] == 0 {
           // min node: sum of children
           sum := 0
           for _, c := range children[v] {
               sum += f[c]
           }
           f[v] = sum
       } else {
           // max node: minimum of children
           mn := f[children[v][0]]
           for _, c := range children[v] {
               if f[c] < mn {
                   mn = f[c]
               }
           }
           f[v] = mn
       }
   }
   // Maximum possible value at root
   res := leafCount - f[1] + 1
   fmt.Fprintln(out, res)
}

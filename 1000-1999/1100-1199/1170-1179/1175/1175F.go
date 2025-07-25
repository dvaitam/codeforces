package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // last occurrence
   lastOcc := make([]int, n+1)
   for i := 1; i <= n; i++ {
       lastOcc[i] = -1
   }
   L := 0
   ans := 0
   // monotonic stack of indices with decreasing a values
   stack := make([]int, 0, n)
   for r := 0; r < n; r++ {
       v := a[r]
       if lastOcc[v] >= 0 && lastOcc[v] >= L {
           L = lastOcc[v] + 1
       }
       lastOcc[v] = r
       // maintain monotonic decreasing stack
       for len(stack) > 0 && a[stack[len(stack)-1]] < v {
           stack = stack[:len(stack)-1]
       }
       stack = append(stack, r)
       // scan stack for valid l's
       for j, idx := range stack {
           M := a[idx]
           l := r - M + 1
           var low int
           if j == 0 {
               low = L
           } else {
               low = stack[j-1] + 1
           }
           if l < low || l > idx {
               continue
           }
           ans++
       }
   }
   fmt.Fprintln(writer, ans)
}

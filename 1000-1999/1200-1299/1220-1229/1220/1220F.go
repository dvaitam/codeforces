package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // find position of value 1
   root := 0
   for i, v := range a {
       if v == 1 {
           root = i
           break
       }
   }
   // prepare rotated arrays
   b := make([][]int, 2)
   b[0] = make([]int, n)
   b[1] = make([]int, n)
   for i := 0; i < n; i++ {
       b[0][i] = a[(root+i)%n]
       b[1][i] = a[(root+n-i)%n]
   }
   // dp arrays
   dpl := make([][]int, 2)
   dpr := make([][]int, 2)
   ans := make([][]int, 2)
   for t := 0; t < 2; t++ {
       dpl[t] = make([]int, n)
       dpr[t] = make([]int, n)
       ans[t] = make([]int, n)
   }
   // compute for each direction
   for t := 0; t < 2; t++ {
       stack := make([]int, 0, n)
       dpl[t][0] = 1
       ans[t][0] = 1
       stack = append(stack, 0)
       for i := 1; i < n; i++ {
           nex := -1
           for len(stack) > 0 && b[t][stack[len(stack)-1]] > b[t][i] {
               j := stack[len(stack)-1]
               stack = stack[:len(stack)-1]
               if nex == -1 {
                   dpr[t][j] = dpl[t][j]
                   nex = j
               } else {
                   dpr[t][j] = max(dpr[t][nex]+1, dpl[t][j])
                   nex = j
               }
           }
           if nex == -1 {
               dpl[t][i] = 1
           } else {
               dpl[t][i] = dpr[t][nex] + 1
           }
           stack = append(stack, i)
           // ans = max(previous, current stack size + dpl - 1)
           cur := len(stack) + dpl[t][i] - 1
           ans[t][i] = max(cur, ans[t][i-1])
       }
   }
   // find minimal result and index
   ret := n + n // large
   k := 0
   for i := 0; i < n; i++ {
       val := max(ans[0][i], ans[1][n-1-i])
       if val < ret {
           ret = val
           k = (root + 1 + i) % n
       }
   }
   // output
   fmt.Println(ret, k)
}

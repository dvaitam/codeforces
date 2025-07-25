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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   p := make([]int, n+2)
   pos := make([]int, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &p[i])
       pos[p[i]] = i
   }
   // previous greater element index
   L := make([]int, n+2)
   stack := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       for len(stack) > 0 && p[stack[len(stack)-1]] < p[i] {
           stack = stack[:len(stack)-1]
       }
       if len(stack) == 0 {
           L[i] = 0
       } else {
           L[i] = stack[len(stack)-1]
       }
       stack = append(stack, i)
   }
   // next greater element index
   R := make([]int, n+2)
   stack = stack[:0]
   for i := n; i >= 1; i-- {
       for len(stack) > 0 && p[stack[len(stack)-1]] < p[i] {
           stack = stack[:len(stack)-1]
       }
       if len(stack) == 0 {
           R[i] = n + 1
       } else {
           R[i] = stack[len(stack)-1]
       }
       stack = append(stack, i)
   }
   var ans int64
   for k := 1; k <= n; k++ {
       leftCount := k - 1 - L[k]
       rightCount := R[k] - k - 1
       if leftCount <= 0 || rightCount <= 0 {
           continue
       }
       if leftCount < rightCount {
           // iterate left
           for i := L[k] + 1; i < k; i++ {
               need := p[k] - p[i]
               if need <= 0 || need > n {
                   continue
               }
               j := pos[need]
               if j > k && j < R[k] {
                   ans++
               }
           }
       } else {
           // iterate right
           for j := k + 1; j < R[k]; j++ {
               need := p[k] - p[j]
               if need <= 0 || need > n {
                   continue
               }
               i := pos[need]
               if i > L[k] && i < k {
                   ans++
               }
           }
       }
   }
   fmt.Fprintln(writer, ans)
}

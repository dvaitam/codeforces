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

   prev := make([]int, n)
   next := make([]int, n)
   stack := make([]int, 0, n)
   // previous smaller element
   for i := 0; i < n; i++ {
       for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
           stack = stack[:len(stack)-1]
       }
       if len(stack) == 0 {
           prev[i] = -1
       } else {
           prev[i] = stack[len(stack)-1]
       }
       stack = append(stack, i)
   }
   // next smaller element
   stack = stack[:0]
   for i := n - 1; i >= 0; i-- {
       for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
           stack = stack[:len(stack)-1]
       }
       if len(stack) == 0 {
           next[i] = n
       } else {
           next[i] = stack[len(stack)-1]
       }
       stack = append(stack, i)
   }

   // ans[len] = maximum of minimums for windows of size len
   ans := make([]int, n+2)
   for i := 0; i < n; i++ {
       length := next[i] - prev[i] - 1
       if a[i] > ans[length] {
           ans[length] = a[i]
       }
   }
   // fill empty entries
   for i := n - 1; i >= 1; i-- {
       if ans[i+1] > ans[i] {
           ans[i] = ans[i+1]
       }
   }

   // output
   for i := 1; i <= n; i++ {
       if i > 1 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, ans[i])
   }
   writer.WriteByte('\n')
}

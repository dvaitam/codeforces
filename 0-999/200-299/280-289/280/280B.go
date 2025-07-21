package main

import (
   "bufio"
   "fmt"
   "os"
)

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
   // Monotonic stack to find pairs of max and second max
   stack := make([]int, 0, n)
   var ans int
   for _, v := range a {
       // Compare with elements in stack
       for len(stack) > 0 {
           top := stack[len(stack)-1]
           xor := top ^ v
           if xor > ans {
               ans = xor
           }
           if top > v {
               break
           }
           // pop top
           stack = stack[:len(stack)-1]
       }
       // push current
       stack = append(stack, v)
   }
   fmt.Println(ans)
}

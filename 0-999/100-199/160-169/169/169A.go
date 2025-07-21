package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, a, b int
   if _, err := fmt.Fscan(reader, &n, &a, &b); err != nil {
       return
   }
   h := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &h[i])
   }
   sort.Ints(h)
   // Vasya does b chores with complexity <= x
   // Petya does a chores with complexity > x, and a + b = n
   // Choose x so that exactly b chores have h_i <= x => x in [h[b-1], h[b]-1]
   L := h[b-1]
   R := h[b]
   ans := R - L
   if ans < 0 {
       ans = 0
   }
   fmt.Println(ans)
}

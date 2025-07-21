package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, X int
   if _, err := fmt.Fscan(in, &n, &X); err != nil {
       return
   }
   c := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &c[i])
   }
   // Compute total consumption weight for each animal
   w := make([]int, n)
   for i := 0; i < n; i++ {
       // animal at day i+1 consumes c[i] per day from day i+1 to n inclusive: total days = n - i
       w[i] = c[i] * (n - i)
   }
   sort.Ints(w)
   cnt, sum := 0, 0
   for _, wi := range w {
       if sum+wi <= X {
           sum += wi
           cnt++
       } else {
           break
       }
   }
   fmt.Println(cnt)
}

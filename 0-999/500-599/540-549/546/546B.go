package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   sort.Ints(a)
   var ans int
   for i := 1; i < n; i++ {
       if a[i] <= a[i-1] {
           diff := a[i-1] + 1 - a[i]
           ans += diff
           a[i] = a[i-1] + 1
       }
   }
   fmt.Println(ans)
}

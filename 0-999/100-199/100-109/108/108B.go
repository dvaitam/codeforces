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
   for i := 0; i < n-1; i++ {
       if a[i] != a[i+1] && a[i]*2 > a[i+1] {
           fmt.Println("YES")
           return
       }
   }
   fmt.Println("NO")
}

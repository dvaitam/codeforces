package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Ints(a)
   res := 0
   for i := 0; i < n && i < m; i++ {
       if a[i] < 0 {
           res += -a[i]
       } else {
           break
       }
   }
   fmt.Println(res)
}

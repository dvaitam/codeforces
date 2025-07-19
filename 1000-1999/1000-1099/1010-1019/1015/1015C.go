package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var m int64
   fmt.Fscan(reader, &n, &m)
   savings := make([]int64, n)
   var totalNew int64
   var totalCurr int64
   for i := 0; i < n; i++ {
       var a, b int64
       fmt.Fscan(reader, &a, &b)
       totalNew += b
       totalCurr += a
       savings[i] = a - b
   }
   if totalNew > m {
       fmt.Println(-1)
       return
   }
   sort.Slice(savings, func(i, j int) bool {
       return savings[i] > savings[j]
   })
   ans := 0
   for totalCurr > m {
       totalCurr -= savings[ans]
       ans++
   }
   fmt.Println(ans)
}

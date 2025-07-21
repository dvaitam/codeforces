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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
   var sum int64
   for i, v := range a {
       idx := int64(i + 1)
       if v > idx {
           sum += v - idx
       } else {
           sum += idx - v
       }
   }
   fmt.Println(sum)
}

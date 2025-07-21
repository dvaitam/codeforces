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
   for _, v := range a {
       sum += v
   }
   var ans int64 = sum
   for i, v := range a {
       c := i + 1
       if c > n-1 {
           c = n - 1
       }
       ans += v * int64(c)
   }
   fmt.Println(ans)
}

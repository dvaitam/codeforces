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
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   sort.Ints(a)
   const mod = 10007
   var sum int64
   for i := 0; i < n; i++ {
       sum = (sum + int64(a[i])*int64(a[n-1-i])) % mod
   }
   fmt.Println(sum)
}

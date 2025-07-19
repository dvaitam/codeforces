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
   const MAXN = 200050
   cnt := make([]int, MAXN)
   ans := 0
   for i := 0; i < n; i++ {
       for j := i + 1; j < n; j++ {
           s := a[i] + a[j]
           cnt[s]++
       }
   }
   for _, c := range cnt {
       if c > ans {
           ans = c
       }
   }
   fmt.Println(ans)
}

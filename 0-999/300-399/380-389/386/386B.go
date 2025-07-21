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
   times := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &times[i])
   }
   var T int
   fmt.Fscan(reader, &T)

   sort.Ints(times)
   ans := 0
   r := 0
   for l := 0; l < n; l++ {
       for r < n && times[r]-times[l] <= T {
           r++
       }
       // window is [l, r-1]
       cnt := r - l
       if cnt > ans {
           ans = cnt
       }
   }
   fmt.Println(ans)
}

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
   // Use max int as infinity
   const INF = int(^uint(0) >> 1)
   ans := INF
   for i := n - 1; i >= 0; i-- {
       u := a[i] - 1
       x := 1
       if a[i] >= ans {
           break
       }
       for j := i + 1; j < n; j++ {
           v1 := u / a[j]
           x += v1
           u -= v1 * a[j]
           v2 := a[i] + a[j] - u - 1
           if v2 >= ans {
               continue
           }
           y := 0
           vv := v2
           k := 0
           for vv != 0 && y <= x {
               y += vv / a[k]
               vv %= a[k]
               k++
           }
           if y > x {
               ans = v2
           }
       }
   }
   if ans == INF {
       ans = -1
   }
   fmt.Println(ans)
}

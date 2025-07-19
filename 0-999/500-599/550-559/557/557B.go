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
   a := make([]int, 2*n)
   for i := 0; i < 2*n; i++ {
       if _, err := fmt.Fscan(reader, &a[i]); err != nil {
           return
       }
   }
   sort.Ints(a)
   a1 := float64(a[0])
   mid := float64(a[n])
   mm := float64(m)
   var ans float64
   if a1*2 <= mid {
       ans = a1 * float64(n) * 3
   } else {
       ans = mid * 1.5 * float64(n)
   }
   if ans > mm {
       ans = mm
   }
   fmt.Printf("%.7f", ans)
}

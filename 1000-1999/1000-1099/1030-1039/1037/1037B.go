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
   var s int
   if _, err := fmt.Fscan(in, &n, &s); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   sort.Ints(a)
   // mid is the index of median (0-based)
   mid := n / 2
   // find last index where a[i] <= s
   pos := sort.Search(n, func(i int) bool { return a[i] > s }) - 1
   var ans int64
   if pos < mid {
       // increase elements from pos+1 to mid to s
       for i := pos + 1; i <= mid; i++ {
           if a[i] < s {
               ans += int64(s - a[i])
           }
       }
   } else {
       // decrease elements from mid to pos to s
       for i := mid; i <= pos; i++ {
           if a[i] > s {
               ans += int64(a[i] - s)
           }
       }
   }
   fmt.Printf("%d", ans)
}

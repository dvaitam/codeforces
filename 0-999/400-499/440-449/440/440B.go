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
   a := make([]int64, n)
   var sum int64
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       sum += a[i]
   }
   avg := sum / int64(n)
   var ans int64
   var pref int64
   for i := 0; i < n-1; i++ {
       pref += a[i] - avg
       if pref < 0 {
           ans -= pref
       } else {
           ans += pref
       }
   }
   fmt.Println(ans)
}

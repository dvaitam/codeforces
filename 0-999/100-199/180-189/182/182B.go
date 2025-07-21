package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var d int
   if _, err := fmt.Fscan(reader, &d); err != nil {
       return
   }
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   total := 0
   for i := 0; i+1 < n; i++ {
       // automatic increment after month i leads to display = (a[i] % d) + 1
       // manual steps to reach day 1: (1 - display + d) % d = (d - (a[i] % d)) % d
       // since a[i] <= d, a[i] % d equals a[i] or 0 when a[i]==d
       if a[i] < d {
           total += d - a[i]
       }
       // if a[i] == d, (d - a[i]) == 0, no need to add
   }
   fmt.Println(total)
}

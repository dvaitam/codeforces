package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   cnt := [3]int{}
   for i := 0; i < n; i++ {
       var a int
       fmt.Fscan(reader, &a)
       cnt[a%3]++
   }
   ans := cnt[0] / 2
   if cnt[1] < cnt[2] {
       ans += cnt[1]
   } else {
       ans += cnt[2]
   }
   fmt.Fprintln(writer, ans)
}

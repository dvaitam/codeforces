package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var k int64
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   var maxJoy int64 = -1 << 60
   for i := 0; i < n; i++ {
       var f, t int64
       fmt.Fscan(reader, &f, &t)
       var joy int64
       if t > k {
           joy = f - (t - k)
       } else {
           joy = f
       }
       if joy > maxJoy {
           maxJoy = joy
       }
   }
   fmt.Println(maxJoy)
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var a int64
   fmt.Fscan(in, &a)
   for i := int64(1); ; i++ {
       x := a + i
       if has8(x) {
           fmt.Fprintln(out, i)
           return
       }
   }
}

// has8 returns true if the decimal representation of x contains digit '8'
func has8(x int64) bool {
   if x < 0 {
       x = -x
   }
   if x == 0 {
       return false
   }
   for x > 0 {
       if x%10 == 8 {
           return true
       }
       x /= 10
   }
   return false
}

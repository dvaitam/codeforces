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
   var str string
   fmt.Fscan(reader, &str)
   s := []byte(str)
   f := make([]int, 10)
   for i := 1; i <= 9; i++ {
       fmt.Fscan(reader, &f[i])
   }
   started := false
   for i := 0; i < n; i++ {
       d := int(s[i] - '0')
       if !started {
           if f[d] > d {
               started = true
               s[i] = byte('0' + f[d])
           }
       } else {
           if f[d] >= d {
               s[i] = byte('0' + f[d])
           } else {
               break
           }
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   writer.Write(s)
}

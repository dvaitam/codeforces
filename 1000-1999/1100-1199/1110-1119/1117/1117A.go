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

   var n int64
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var a int64
   if _, err := fmt.Fscan(reader, &a); err != nil {
       return
   }
   now := a
   savenow := a
   l := int64(1)
   savel := int64(1)
   for i := int64(1); i < n; i++ {
       if _, err := fmt.Fscan(reader, &a); err != nil {
           break
       }
       if a > now/l {
           now = a
           l = 1
           if a > savenow/savel {
               savenow = a
               savel = 1
           }
       } else if (now+a)/(l+1) == now/l {
           now += a
           l++
           if now/l >= savenow/savel && l >= savel {
               savel = l
               savenow = now
           }
       } else {
           now = a
           l = 1
       }
   }
   fmt.Fprint(writer, savel)
}

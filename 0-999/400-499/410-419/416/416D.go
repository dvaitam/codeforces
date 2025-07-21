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
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   ans := 0
   i := 0
   for i < n {
       ans++
       firstIdx := -1
       var firstVal, d int64
       haveD := false
       j := i
       for ; j < n; j++ {
           v := a[j]
           if v != -1 {
               if firstIdx < 0 {
                   firstIdx = j
                   firstVal = v
               } else if !haveD {
                   dist := int64(j - firstIdx)
                   delta := v - firstVal
                   if delta%dist != 0 {
                       break
                   }
                   d = delta / dist
                   haveD = true
               } else {
                   exp := firstVal + int64(j-firstIdx)*d
                   if exp != v {
                       break
                   }
               }
           }
       }
       if j == i {
           j++
       }
       i = j
   }
   fmt.Fprintln(writer, ans)
}

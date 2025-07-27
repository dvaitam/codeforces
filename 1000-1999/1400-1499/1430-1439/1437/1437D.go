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

   var t int
   fmt.Fscan(reader, &t)
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       // BFS levels reconstruction for minimal height
       head, tail := 0, 0
       idx := 1
       levels := 0
       for idx < n {
           levels++
           newTail := idx
           // assign children for each parent in current layer
           for i := head; i <= tail && idx < n; i++ {
               j := idx
               // group sorted children of parent a[i]
               for j+1 < n && a[j+1] > a[j] {
                   j++
               }
               // consume positions idx..j as children
               idx = j + 1
               if j > newTail {
                   newTail = j
               }
           }
           // move to next layer
           head = tail + 1
           tail = newTail
       }
       // print minimal height
       fmt.Fprintln(writer, levels)
   }
}

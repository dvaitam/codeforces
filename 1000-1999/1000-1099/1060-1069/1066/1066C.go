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

   var q int
   if _, err := fmt.Fscan(reader, &q); err != nil {
       return
   }
   pos := make(map[int]int)
   l, r := 0, 0
   size := 0
   for i := 0; i < q; i++ {
       var op string
       var id int
       fmt.Fscan(reader, &op, &id)
       switch op {
       case "L":
           if size == 0 {
               pos[id] = 0
               l, r = 0, 0
               size = 1
           } else {
               l--
               pos[id] = l
               size++
           }
       case "R":
           if size == 0 {
               pos[id] = 0
               l, r = 0, 0
               size = 1
           } else {
               r++
               pos[id] = r
               size++
           }
       case "?":
           p := pos[id]
           left := p - l
           right := r - p
           if left < right {
               fmt.Fprintln(writer, left)
           } else {
               fmt.Fprintln(writer, right)
           }
       }
   }
}

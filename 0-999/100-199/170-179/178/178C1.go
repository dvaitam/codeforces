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

   var h, m, n int
   fmt.Fscan(reader, &h, &m, &n)
   table := make([]bool, h)
   idPos := make(map[int]int, n)
   var total int64
   for i := 0; i < n; i++ {
       var op string
       fmt.Fscan(reader, &op)
       if op == "+" {
           var id, hash int
           fmt.Fscan(reader, &id, &hash)
           pos := hash
           var cnt int64
           for table[pos] {
               cnt++
               pos += m
               if pos >= h {
                   pos %= h
               }
           }
           table[pos] = true
           idPos[id] = pos
           total += cnt
       } else {
           var id int
           fmt.Fscan(reader, &id)
           pos := idPos[id]
           table[pos] = false
           delete(idPos, id)
       }
   }
   fmt.Fprintln(writer, total)
}

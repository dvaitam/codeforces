package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   var n int
   if _, err := fmt.Fscan(os.Stdin, &n); err != nil {
       return
   }
   // Initialize with sequence for 1
   old := []string{"1#"}
   // Build sequences for 2..n
   for c := 2; c <= n; c++ {
       ch := rune('0' + c)
       dir := 1
       newSeq := make([]string, 0)
       // Process each existing sequence
       for len(old) > 0 {
           sold := old[0]
           old = old[1:]
           var start int
           if dir == 1 {
               start = 0
           } else {
               seq := sold + string(ch) + "#"
               newSeq = append(newSeq, seq)
               start = len(sold) - 1
           }
           // Insert at each '#' position
           pos := start
           for pos >= 0 && pos < len(sold) {
               if sold[pos] == '#' {
                   seq := sold[:pos] + string(ch) + sold[pos:]
                   newSeq = append(newSeq, seq)
               }
               pos += dir
           }
           if dir == 1 {
               seq := sold + string(ch) + "#"
               newSeq = append(newSeq, seq)
           }
           dir = -dir
       }
       old = newSeq
   }
   // Output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, len(old))
   for _, s := range old {
       isFirst := true
       for i := 0; i < len(s); i++ {
           if s[i] == '#' {
               writer.WriteByte('}')
               isFirst = true
           } else {
               if !isFirst {
                   writer.WriteByte(',')
               } else {
                   if i != 0 {
                       writer.WriteByte(',')
                   }
                   writer.WriteByte('{')
               }
               isFirst = false
               if s[i] <= '9' {
                   writer.WriteByte(s[i])
               } else {
                   writer.WriteString("10")
               }
           }
       }
       writer.WriteByte('\n')
   }
}

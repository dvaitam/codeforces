package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, q int
   fmt.Fscan(reader, &n, &q)
   var s string
   fmt.Fscan(reader, &s)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < q; i++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       seq := []byte(s[l-1 : r])
       counts := make([]int, 10)
       cp := 0
       dp := 1
       for cp >= 0 && cp < len(seq) {
           c := seq[cp]
           if c >= '0' && c <= '9' {
               d := int(c - '0')
               counts[d]++
               pos := cp
               cp += dp
               if d > 0 {
                   // decrease digit in sequence at pos
                   // note: pos may be out of bounds after cp move, but seq unchanged length
                   seq[pos] = byte('0' + d - 1)
               } else {
                   // delete digit at pos
                   seq = append(seq[:pos], seq[pos+1:]...)
                   // adjust cp if needed
                   if cp > pos {
                       cp--
                   }
               }
           } else {
               // bracket
               pos := cp
               if c == '<' {
                   dp = -1
               } else {
                   dp = 1
               }
               cp += dp
               // after move, if new char is bracket, delete previous char at pos
               if cp >= 0 && cp < len(seq) {
                   nc := seq[cp]
                   if nc == '<' || nc == '>' {
                       seq = append(seq[:pos], seq[pos+1:]...)
                       if cp > pos {
                           cp--
                       }
                   }
               }
           }
       }
       // output counts
       for j := 0; j < 10; j++ {
           if j > 0 {
               fmt.Fprint(writer, " ")
           }
           fmt.Fprint(writer, counts[j])
       }
       fmt.Fprintln(writer)
   }
}

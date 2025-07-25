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
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   freq := make([]int, 26)
   for i := 0; i < len(s); i++ {
       c := s[i] - 'a'
       if c >= 0 && c < 26 {
           freq[c]++
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // Arrange characters in contiguous blocks by letter
   for i := 0; i < 26; i++ {
       for cnt := 0; cnt < freq[i]; cnt++ {
           writer.WriteByte(byte('a' + i))
       }
   }
}

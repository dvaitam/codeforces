package main

import (
   "bufio"
   "fmt"
   "os"
)

func maxLen(s string, k int, ch byte) int {
   left := 0
   count := 0
   best := 0
   for right := 0; right < len(s); right++ {
      if s[right] != ch {
         count++
      }
      for count > k {
         if s[left] != ch {
            count--
         }
         left++
      }
      if cur := right - left + 1; cur > best {
         best = cur
      }
   }
   return best
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
      return
   }
   var s string
   fmt.Fscan(reader, &s)

   resA := maxLen(s, k, 'a')
   resB := maxLen(s, k, 'b')
   if resA > resB {
      fmt.Fprintln(writer, resA)
   } else {
      fmt.Fprintln(writer, resB)
   }
}

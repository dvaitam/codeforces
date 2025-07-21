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

   var k int
   var s string
   fmt.Fscan(reader, &k, &s)

   n := len(s)
   if n%k != 0 {
       fmt.Fprintln(writer, -1)
       return
   }

   // count frequency of each letter
   freq := [26]int{}
   for i := 0; i < n; i++ {
       freq[s[i]-'a']++
   }
   // check divisibility
   for j := 0; j < 26; j++ {
       if freq[j]%k != 0 {
           fmt.Fprintln(writer, -1)
           return
       }
   }

   // build base part of length n/k
   part := make([]byte, 0, n/k)
   for j := 0; j < 26; j++ {
       cnt := freq[j] / k
       for i := 0; i < cnt; i++ {
           part = append(part, byte('a'+j))
       }
   }

   // output k repeats of part
   for i := 0; i < k; i++ {
       writer.Write(part)
   }
   writer.WriteByte('\n')
}

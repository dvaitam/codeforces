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
   for i := 0; i < n; i++ {
       var s, t string
       fmt.Fscan(reader, &s, &t)
       if s < t {
           fmt.Fprintln(writer, s)
           continue
       }
       var b [26]int
       for j := 0; j < len(s); j++ {
           b[s[j]-'A']++
       }
       id := -1
       var cc int
       for l := 0; l < len(s); l++ {
           b[s[l]-'A']--
           for j := 0; j < 26; j++ {
               if b[j] > 0 && int(s[l]-'A') > j {
                   id = l
                   cc = j
                   break
               }
           }
           if id >= 0 {
               break
           }
       }
       if id == -1 {
           fmt.Fprintln(writer, "---")
           continue
       }
       ss := []byte(s)
       for j := len(ss) - 1; j > id; j-- {
           if int(ss[j]-'A') == cc {
               ss[j] = ss[id]
               break
           }
       }
       ss[id] = byte(cc) + 'A'
       s2 := string(ss)
       if s2 < t {
           fmt.Fprintln(writer, s2)
       } else {
           fmt.Fprintln(writer, "---")
       }
   }
}

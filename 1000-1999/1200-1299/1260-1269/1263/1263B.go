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
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       codes := make([]string, n)
       mp := make(map[string]int)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &codes[i])
           mp[codes[i]]++
       }
       cnt := 0
       for i := 0; i < n; i++ {
           if mp[codes[i]] > 1 {
               cnt++
               orig := codes[i]
               tmp := []byte(orig)
               var ok bool
               // try changing fourth digit
               for j := byte(0); j <= 9; j++ {
                   tmp[3] = '0' + j
                   s := string(tmp)
                   if mp[s] == 0 {
                       ok = true
                       codes[i] = s
                       break
                   }
               }
               // if not found, try changing third digit
               if !ok {
                   tmp = []byte(orig)
                   for j := byte(0); j <= 9; j++ {
                       tmp[2] = '0' + j
                       s := string(tmp)
                       if mp[s] == 0 {
                           ok = true
                           codes[i] = s
                           break
                       }
                   }
               }
               mp[orig]--
               mp[codes[i]]++
           }
       }
       fmt.Fprintln(writer, cnt)
       for _, s := range codes {
           fmt.Fprintln(writer, s)
       }
   }
}

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

   var n, m int
   fmt.Fscan(reader, &n, &m)
   mp := make(map[string]int)
   left := make([]string, 0)
   right := make([]string, 0)
   var q string
   ans := 0
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       t := reverseString(s)
       if mp[t] > 0 {
           mp[t]--
           // add to front of left
           left = append([]string{t}, left...)
           // add to end of right
           right = append(right, s)
           ans += 2 * len(s)
       } else if s == t {
           q = s
       } else {
           mp[s]++
       }
   }
   totalLen := ans + len(q)
   fmt.Fprintln(writer, totalLen)
   for _, str := range left {
       fmt.Fprint(writer, str)
   }
   fmt.Fprint(writer, q)
   for _, str := range right {
       fmt.Fprint(writer, str)
   }
}

func reverseString(s string) string {
   bs := []rune(s)
   for i, j := 0, len(bs)-1; i < j; i, j = i+1, j-1 {
       bs[i], bs[j] = bs[j], bs[i]
   }
   return string(bs)
}

package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func sortString(s string) string {
   b := []byte(s)
   sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
   return string(b)
}

func reverseString(s string) string {
   b := []byte(s)
   for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
       b[i], b[j] = b[j], b[i]
   }
   return string(b)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var s string
       fmt.Fscan(reader, &s)
       a := sortString(s)
       b := reverseString(a)
       if a == s && b == s {
           fmt.Fprintln(writer, "NO")
       } else {
           fmt.Fprintln(writer, "YES")
           if a == s {
               fmt.Fprintln(writer, b)
           } else {
               fmt.Fprintln(writer, a)
           }
       }
   }
}

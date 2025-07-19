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

   var n, k int
   var str string
   fmt.Fscan(reader, &n, &k, &str)
   s := []byte(str)

   if k == 2 {
       // Only two characters: choose best of ABAB... or BABA...
       d := 0
       for i := 0; i < n; i++ {
           if s[i] == byte('A'+byte(i&1)) {
               d++
           }
       }
       var changes, offset int
       if d <= n-d {
           changes = d
           offset = 1
       } else {
           changes = n - d
           offset = 0
       }
       for i := 0; i < n; i++ {
           s[i] = byte('A' + byte((offset+i)&1))
       }
       fmt.Fprintln(writer, changes)
       writer.Write(s)
       return
   }

   // k > 2: greedy change when two consecutive are equal
   changes := 0
   for i := 0; i+1 < n; i++ {
       if s[i] == s[i+1] {
           changes++
           // choose a different character not equal to neighbors
           for c := byte('A'); ; c++ {
               // next next char (if exists)
               if c != s[i] && (i+2 >= n || c != s[i+2]) {
                   s[i+1] = c
                   break
               }
           }
       }
   }
   fmt.Fprintln(writer, changes)
   writer.Write(s)
}

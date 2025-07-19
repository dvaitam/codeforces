package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)

   v := make([][]string, 10)
   v1 := make([][]string, 10)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       l := len(s)
       if l < len(v) {
           v[l] = append(v[l], s)
       }
   }
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       l := len(s)
       if l < len(v1) {
           v1[l] = append(v1[l], s)
       }
   }

   var ans, ss int64
   for i := 1; i < 5; i++ {
       sort.Strings(v[i])
       sort.Strings(v1[i])
       var ar [26]int64
       for j := 0; j < len(v[i]); j++ {
           for k := 0; k < len(v[i][j]); k++ {
               idx := int(v[i][j][k]) - int('a') + 32
               ar[idx]++
               idx1 := int(v1[i][j][k]) - int('a') + 32
               ar[idx1]--
           }
       }
       for k := 0; k < 26; k++ {
           if ar[k] < 0 {
               ans += -ar[k]
           } else {
               ans += ar[k]
           }
       }
       ans /= 2
       ss += ans
       ans = 0
   }
   fmt.Fprint(writer, ss)
}

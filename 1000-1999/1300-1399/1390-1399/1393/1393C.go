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

   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   for T > 0 {
       T--
       var n int
       fmt.Fscan(reader, &n)
       freq := make(map[int]int)
       for i := 0; i < n; i++ {
           var v int
           fmt.Fscan(reader, &v)
           freq[v]++
       }
       // find maximum frequency and its count
       mX := 0
       for _, f := range freq {
           if f > mX {
               mX = f
           }
       }
       c := 0
       for _, f := range freq {
           if f == mX {
               c++
           }
       }
       // using scheduling formula: minimal required length = (mX-1)*(d+1) + c <= n
       // solve for maximum d: (mX-1)*(d+1) + c <= n => d <= (n-c)/(mX-1) - 1
       ans := (n - c) / (mX - 1) - 1
       if ans < 0 {
           ans = 0
       }
       fmt.Fprintln(writer, ans)
   }
}

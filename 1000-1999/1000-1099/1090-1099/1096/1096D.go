package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   var s string
   fmt.Fscan(reader, &s)

   target := "hard"
   inf := int64(1) << 60
   f := make([]int64, 5)
   for i := 1; i <= 4; i++ {
       f[i] = inf
   }

   for i := 0; i < n; i++ {
       var a int64
       fmt.Fscan(reader, &a)
       for j := 3; j >= 0; j-- {
           if s[i] == target[j] {
               f[j+1] = min(f[j+1], f[j])
               f[j] += a
           }
       }
   }

   res := f[0]
   for j := 1; j <= 3; j++ {
       if f[j] < res {
           res = f[j]
       }
   }
   fmt.Fprint(writer, res)
}

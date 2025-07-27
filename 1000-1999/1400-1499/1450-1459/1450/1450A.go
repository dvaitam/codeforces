package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var t int
   fmt.Fscan(in, &t)
   for ; t > 0; t-- {
       var n int
       var s string
       fmt.Fscan(in, &n, &s)
       bs := []byte(s)
       sort.Slice(bs, func(i, j int) bool { return bs[i] < bs[j] })
       out.Write(bs)
       out.WriteByte('\n')
   }
}

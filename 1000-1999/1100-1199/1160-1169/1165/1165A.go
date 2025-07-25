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

   var n, x, y int
   if _, err := fmt.Fscan(reader, &n, &x, &y); err != nil {
       return
   }
   var s string
   fmt.Fscan(reader, &s)

   ans := 0
   // Only last x digits affect modulo 10^x
   // We need remainder 10^y, i.e., digit at position y (0-indexed from least significant) is 1, others 0
   for i := 0; i < x; i++ {
       target := byte('0')
       if i == y {
           target = '1'
       }
       if s[n-1-i] != target {
           ans++
       }
   }
   fmt.Fprintln(writer, ans)

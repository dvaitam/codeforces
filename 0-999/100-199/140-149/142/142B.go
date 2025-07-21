package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m int
   fmt.Fscan(reader, &n, &m)
   // Ensure small <= big
   small, big := n, m
   if small > big {
       small, big = big, small
   }
   var ans int
   switch small {
   case 1:
       // All cells can be occupied
       ans = big
   case 2:
       // Special handling: place in blocks of 4 columns
       blocks := big / 4
       rem := big % 4
       ans = blocks*4 + min(rem*2, 4)
   default:
       // General case: fill half (ceil)
       ans = (n*m + 1) / 2
   }
   fmt.Fprintln(writer, ans)
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

func isSorted(a []int) bool {
   for i := 1; i < len(a); i++ {
       if a[i-1] > a[i] {
           return false
       }
   }
   return true
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for tc := 0; tc < t; tc++ {
       var n, x int
       fmt.Fscan(reader, &n, &x)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       ops := 0
       // simulate operations
       for !isSorted(a) {
           swapped := false
           for i := 0; i < n; i++ {
               if a[i] > x {
                   // swap a[i] and x
                   a[i], x = x, a[i]
                   ops++
                   swapped = true
                   break
               }
           }
           if !swapped {
               ops = -1
               break
           }
       }
       fmt.Fprintln(writer, ops)
   }
}

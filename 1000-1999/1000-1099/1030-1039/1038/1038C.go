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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   b := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &b[i])
   }
   sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
   sort.Slice(b, func(i, j int) bool { return b[i] > b[j] })

   var sumA, sumB int64
   i, j := 0, 0
   turn := 0 // 0: A's turn, 1: B's turn
   for i < n || j < n {
       if turn == 0 {
           if i < n && (j >= n || a[i] > b[j]) {
               sumA += int64(a[i])
               i++
           } else {
               j++
           }
           turn = 1
       } else {
           if j < n && (i >= n || b[j] > a[i]) {
               sumB += int64(b[j])
               j++
           } else {
               i++
           }
           turn = 0
       }
   }
   fmt.Fprint(writer, sumA-sumB)
}

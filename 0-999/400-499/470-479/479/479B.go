package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, k int
   fmt.Fscan(in, &n, &k)
   a := make([]int, n)
   var sum int64
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
       sum += int64(a[i])
   }

   minDiff := 0
   if sum%int64(n) != 0 {
       minDiff = 1
   }

   minH, maxH := 0, 0
   for i := 0; i < n; i++ {
       if a[i] < a[minH] {
           minH = i
       }
       if a[i] > a[maxH] {
           maxH = i
       }
   }
   currDiff := a[maxH] - a[minH]
   if currDiff == minDiff {
       fmt.Fprintln(out, currDiff, 0)
       return
   }

   ops := make([][2]int, 0, k)
   for t := 0; t < k; t++ {
       a[maxH]--
       a[minH]++
       ops = append(ops, [2]int{maxH + 1, minH + 1})
       // recompute minH and maxH
       minH, maxH = 0, 0
       for i := 0; i < n; i++ {
           if a[i] < a[minH] {
               minH = i
           }
           if a[i] > a[maxH] {
               maxH = i
           }
       }
       currDiff = a[maxH] - a[minH]
       if currDiff == minDiff {
           break
       }
   }

   t := len(ops)
   fmt.Fprintln(out, currDiff, t)
   for i := 0; i < t; i++ {
       fmt.Fprintln(out, ops[i][0], ops[i][1])
   }
}

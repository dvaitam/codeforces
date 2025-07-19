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
   for t := 0; t < T; t++ {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       b := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &b[i])
       }
       if isSorted(a) && isSorted(b) {
           fmt.Fprintln(writer, 0)
           continue
       }
       var ops [][2]int
       for i := 0; i < n-1; i++ {
           minA, minB := a[i], b[i]
           for j := i + 1; j < n; j++ {
               if a[j] < minA {
                   minA = a[j]
               }
               if b[j] < minB {
                   minB = b[j]
               }
           }
           index := -1
           for j := i + 1; j < n; j++ {
               if a[j] == minA && b[j] == minB {
                   index = j
                   break
               }
           }
           if index != -1 {
               a[i], a[index] = a[index], a[i]
               b[i], b[index] = b[index], b[i]
               ops = append(ops, [2]int{i + 1, index + 1})
           }
       }
       if isSorted(a) && isSorted(b) {
           fmt.Fprintln(writer, len(ops))
           for _, p := range ops {
               fmt.Fprintln(writer, p[0], p[1])
           }
       } else {
           fmt.Fprintln(writer, -1)
       }
   }
}

func isSorted(a []int) bool {
   for i := 1; i < len(a); i++ {
       if a[i] < a[i-1] {
           return false
       }
   }
   return true
}

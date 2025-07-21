package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   a := make([]int, n)
   sign := make([]byte, n)
   countPlus := make([]int, n+1)
   countMinus := make([]int, n+1)
   totalMinus := 0

   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       sign[i] = s[0]
       ai, _ := strconv.Atoi(s[1:])
       a[i] = ai
       if sign[i] == '+' {
           countPlus[ai]++
       } else {
           countMinus[ai]++
           totalMinus++
       }
   }
   // Determine possible criminals
   inS := make([]bool, n+1)
   var Ssize int
   for c := 1; c <= n; c++ {
       // truths when criminal is c
       t := countPlus[c] + (totalMinus - countMinus[c])
       if t == m {
           inS[c] = true
           Ssize++
       }
   }
   // For each statement, decide
   for i := 0; i < n; i++ {
       x := a[i]
       if sign[i] == '+' {
           if inS[x] && Ssize == 1 {
               fmt.Fprintln(writer, "Truth")
           } else if !inS[x] {
               fmt.Fprintln(writer, "Lie")
           } else {
               fmt.Fprintln(writer, "Not defined")
           }
       } else {
           // '-'
           if inS[x] && Ssize == 1 {
               fmt.Fprintln(writer, "Lie")
           } else if !inS[x] {
               fmt.Fprintln(writer, "Truth")
           } else {
               fmt.Fprintln(writer, "Not defined")
           }
       }
   }
}

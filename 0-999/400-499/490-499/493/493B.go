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
   var n int
   fmt.Fscan(reader, &n)
   var firstSum, secondSum int64
   firstSeq := make([]int64, 0, n)
   secondSeq := make([]int64, 0, n)
   last := 0
   for i := 0; i < n; i++ {
       var ai int64
       fmt.Fscan(reader, &ai)
       if ai > 0 {
           firstSum += ai
           firstSeq = append(firstSeq, ai)
           last = 1
       } else {
           val := -ai
           secondSum += val
           secondSeq = append(secondSeq, val)
           last = 2
       }
   }
   if firstSum > secondSum {
       fmt.Fprintln(writer, "first")
       return
   } else if firstSum < secondSum {
       fmt.Fprintln(writer, "second")
       return
   }
   // sums equal, compare lexicographically
   minLen := len(firstSeq)
   if len(secondSeq) < minLen {
       minLen = len(secondSeq)
   }
   for i := 0; i < minLen; i++ {
       if firstSeq[i] > secondSeq[i] {
           fmt.Fprintln(writer, "first")
           return
       } else if firstSeq[i] < secondSeq[i] {
           fmt.Fprintln(writer, "second")
           return
       }
   }
   if len(firstSeq) > len(secondSeq) {
       fmt.Fprintln(writer, "first")
       return
   } else if len(firstSeq) < len(secondSeq) {
       fmt.Fprintln(writer, "second")
       return
   }
   // sequences equal, last move wins
   if last == 1 {
       fmt.Fprintln(writer, "first")
   } else {
       fmt.Fprintln(writer, "second")
   }
}

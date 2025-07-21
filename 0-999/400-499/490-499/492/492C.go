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
   var r, avg int64
   fmt.Fscan(reader, &n, &r, &avg)
   exams := make([]exam, n)
   var sumGrades int64
   for i := 0; i < n; i++ {
       var ai, bi int64
       fmt.Fscan(reader, &ai, &bi)
       exams[i] = exam{cost: bi, grade: ai}
       sumGrades += ai
   }
   // Calculate total required increase
   need := int64(n)*avg - sumGrades
   if need <= 0 {
       fmt.Fprintln(writer, 0)
       return
   }
   // Sort by cost ascending
   sort.Slice(exams, func(i, j int) bool {
       return exams[i].cost < exams[j].cost
   })
   var totalCost int64
   for i := 0; i < n && need > 0; i++ {
       // Maximum points we can add for this exam
       canAdd := r - exams[i].grade
       if canAdd <= 0 {
           continue
       }
       var delta int64
       if need < canAdd {
           delta = need
       } else {
           delta = canAdd
       }
       totalCost += delta * exams[i].cost
       need -= delta
   }
   fmt.Fprintln(writer, totalCost)
}

type exam struct {
   cost  int64
   grade int64
}

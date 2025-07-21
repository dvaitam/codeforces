package main

import (
   "fmt"
   "sort"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   exams := make([]struct{ a, b int }, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&exams[i].a, &exams[i].b)
   }
   // Sort by scheduled date a, then by early date b
   sort.Slice(exams, func(i, j int) bool {
       if exams[i].a != exams[j].a {
           return exams[i].a < exams[j].a
       }
       return exams[i].b < exams[j].b
   })
   currentDay := 0
   for _, ex := range exams {
       if ex.b >= currentDay {
           currentDay = ex.b
       } else {
           currentDay = ex.a
       }
   }
   fmt.Println(currentDay)
}

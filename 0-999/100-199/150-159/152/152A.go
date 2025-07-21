package main

import (
   "fmt"
)

func main() {
   var n, m int
   if _, err := fmt.Scan(&n, &m); err != nil {
       return
   }
   grades := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&grades[i])
   }
   // find max in each subject (column)
   maxGrades := make([]byte, m)
   for j := 0; j < m; j++ {
       maxGrades[j] = '0'
       for i := 0; i < n; i++ {
           if grades[i][j] > maxGrades[j] {
               maxGrades[j] = grades[i][j]
           }
       }
   }
   // count successful students
   count := 0
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grades[i][j] == maxGrades[j] {
               count++
               break
           }
       }
   }
   fmt.Println(count)
}

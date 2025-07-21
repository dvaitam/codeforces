package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var year, day int
   if _, err := fmt.Fscan(reader, &year, &day); err != nil {
       return
   }
   // Determine if it's a leap year
   isLeap := (year%400 == 0) || (year%4 == 0 && year%100 != 0)
   // Days in each month
   months := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
   if isLeap {
       months[1] = 29
   }
   // Find month and day
   m := 0
   d := day
   for i, dim := range months {
       if d > dim {
           d -= dim
       } else {
           m = i + 1
           break
       }
   }
   fmt.Printf("%d %d\n", d, m)
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var day1, day2 string
   if _, err := fmt.Fscan(reader, &day1, &day2); err != nil {
       return
   }
   // Map weekday names to integers 0-6
   daysMap := map[string]int{
       "monday":    0,
       "tuesday":   1,
       "wednesday": 2,
       "thursday":  3,
       "friday":    4,
       "saturday":  5,
       "sunday":    6,
   }
   start, ok1 := daysMap[day1]
   target, ok2 := daysMap[day2]
   if !ok1 || !ok2 {
       fmt.Println("NO")
       return
   }
   // Month lengths in a non-leap year
   monthDays := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
   for _, d := range monthDays {
       if (start+d)%7 == target {
           fmt.Println("YES")
           return
       }
   }
   fmt.Println("NO")
}

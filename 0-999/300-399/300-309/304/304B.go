package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
)

func isLeap(year int) bool {
   if year%400 == 0 {
       return true
   }
   if year%100 == 0 {
       return false
   }
   return year%4 == 0
}

// countDays returns the number of days from year 1-01-01 to the given date
func countDays(year, month, day int) int64 {
   // days in full years before this year
   y := int64(year - 1)
   leaps := y/4 - y/100 + y/400
   days := y*365 + leaps
   // days in months before this month in this year
   for m := 1; m < month; m++ {
       switch m {
       case 1, 3, 5, 7, 8, 10, 12:
           days += 31
       case 4, 6, 9, 11:
           days += 30
       case 2:
           if isLeap(year) {
               days += 29
           } else {
               days += 28
           }
       }
   }
   // days in current month
   days += int64(day)
   return days
}

func abs64(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   var dates [2]string
   for i := 0; i < 2; i++ {
       if !scanner.Scan() {
           return
       }
       dates[i] = strings.TrimSpace(scanner.Text())
   }
   var days [2]int64
   for i, ds := range dates {
       parts := strings.Split(ds, ":")
       y, _ := strconv.Atoi(parts[0])
       m, _ := strconv.Atoi(parts[1])
       d, _ := strconv.Atoi(parts[2])
       days[i] = countDays(y, m, d)
   }
   diff := abs64(days[1] - days[0]) + 1
   fmt.Println(diff)
}

package main

import (
   "fmt"
)

func main() {
   var s int
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   ones := map[int]string{
       0: "zero", 1: "one", 2: "two", 3: "three", 4: "four", 5: "five",
       6: "six", 7: "seven", 8: "eight", 9: "nine", 10: "ten",
       11: "eleven", 12: "twelve", 13: "thirteen", 14: "fourteen", 15: "fifteen",
       16: "sixteen", 17: "seventeen", 18: "eighteen", 19: "nineteen",
   }
   tens := map[int]string{
       20: "twenty", 30: "thirty", 40: "forty", 50: "fifty",
       60: "sixty", 70: "seventy", 80: "eighty", 90: "ninety",
   }
   var ans string
   switch {
   case s < 20:
       ans = ones[s]
   case s%10 == 0:
       ans = tens[s]
   default:
       t := s - s%10
       o := s % 10
       ans = tens[t] + "-" + ones[o]
   }
   fmt.Println(ans)
}

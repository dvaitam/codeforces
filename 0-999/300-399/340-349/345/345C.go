package main

import (
   "fmt"
   "time"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   count := 0
   for i := 0; i < n; i++ {
       var s string
       fmt.Scan(&s)
       // check if day is 13
       if len(s) >= 10 && s[8:10] == "13" {
           t, err := time.Parse("2006-01-02", s)
           if err != nil {
               continue
           }
           if t.Weekday() == time.Friday {
               count++
           }
       }
   }
   fmt.Println(count)
}

package main

import (
   "fmt"
)

func main() {
   var da, db int64
   if _, err := fmt.Scan(&da, &db); err != nil {
       return
   }
   if db-da == 1 {
       fmt.Printf("%d %d", db*10-1, db*10)
   } else if da == db {
       fmt.Printf("%d %d", da*10, db*10+1)
   } else if da == 9 && db == 1 {
       fmt.Printf("99 100")
   } else {
       fmt.Printf("-1")
   }
}

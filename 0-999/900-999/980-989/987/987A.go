package main

import "fmt"

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   stones := map[string]string{
       "purple": "Power",
       "green":  "Time",
       "blue":   "Space",
       "orange": "Soul",
       "red":    "Reality",
       "yellow": "Mind",
   }
   for i := 0; i < n; i++ {
       var color string
       fmt.Scan(&color)
       delete(stones, color)
   }
   fmt.Println(len(stones))
   for _, name := range stones {
       fmt.Println(name)
   }
}

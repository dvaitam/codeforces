package main

import (
   "fmt"
)

func main() {
   // Print commands for increment algorithm
   for i := 0; i < 9; i++ {
       fmt.Printf("%d??<>%d\n", i, i+1)
   }
   // Handle rollover from 9 to 0 and start termination condition
   fmt.Println("9??>>??0")
   fmt.Println("??<>1")
   // Preserve other digits unchanged
   for i := 0; i < 10; i++ {
       fmt.Printf("?%d>>%d?\n", i, i)
   }
   // Final termination commands
   fmt.Println("?>>??")
   fmt.Println(">>?")
}

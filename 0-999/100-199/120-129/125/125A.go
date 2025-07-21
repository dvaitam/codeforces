package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   // Convert centimeters to inches (1 inch = 3 cm) with rounding
   inches := n / 3
   if n%3 == 2 {
       inches++
   }
   // Convert to feet and remaining inches (1 foot = 12 inches)
   feet := inches / 12
   inches = inches % 12
   fmt.Fprintln(writer, feet, inches)
}

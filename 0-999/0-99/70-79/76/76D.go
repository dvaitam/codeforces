package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   // read A
   if !scanner.Scan() {
       return
   }
   aStr := scanner.Text()
   // read B
   if !scanner.Scan() {
       return
   }
   bStr := scanner.Text()
   
   a, err := strconv.ParseUint(aStr, 10, 64)
   if err != nil {
       fmt.Println(-1)
       return
   }
   b, err := strconv.ParseUint(bStr, 10, 64)
   if err != nil {
       fmt.Println(-1)
       return
   }
   // Check feasibility: A >= B and (A-B) even
   if a < b {
       fmt.Println(-1)
       return
   }
   diff := a - b
   if diff&1 != 0 {
       fmt.Println(-1)
       return
   }
   d := diff >> 1
   // X & Y = d, X ^ Y = b => d & b must be zero
   if d&b != 0 {
       fmt.Println(-1)
       return
   }
   x := d
   y := d + b
   fmt.Printf("%d %d", x, y)
}

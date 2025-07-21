package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var x1, y1, z1 int
   var x2, y2, z2 int
   if _, err := fmt.Fscan(reader, &x1, &y1, &z1); err != nil {
       return
   }
   if _, err := fmt.Fscan(reader, &x2, &y2, &z2); err != nil {
       return
   }
   // Flies see each other if they lie on the same face of the cube,
   // i.e., share a coordinate (x, y or z) equal to 0 or 1
   if x1 == x2 || y1 == y2 || z1 == z2 {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}

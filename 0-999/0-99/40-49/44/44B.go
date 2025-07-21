package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   scanner.Split(bufio.ScanWords)
   if !scanner.Scan() {
       return
   }
   n, _ := strconv.Atoi(scanner.Text())
   if !scanner.Scan() {
       return
   }
   a, _ := strconv.Atoi(scanner.Text())
   if !scanner.Scan() {
       return
   }
   b, _ := strconv.Atoi(scanner.Text())
   if !scanner.Scan() {
       return
   }
   c, _ := strconv.Atoi(scanner.Text())

   // Count solutions to x + 2y + 4z = 2*n with 0<=x<=a, 0<=y<=b, 0<=z<=c
   total := int64(0)
   zmax := n / 2
   if c < zmax {
       zmax = c
   }
   for z := 0; z <= zmax; z++ {
       S := 2*n - 4*z
       // Solve x + 2y = S
       // y bounds: L = max(0, ceil((S-a)/2)), U = min(b, floor(S/2))
       t := S - a
       L := (t + 1) / 2
       if L < 0 {
           L = 0
       }
       U := S / 2
       if U > b {
           U = b
       }
       if L <= U {
           total += int64(U - L + 1)
       }
   }
   fmt.Println(total)
}

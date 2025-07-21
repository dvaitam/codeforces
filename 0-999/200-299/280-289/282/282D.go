package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   switch n {
   case 1:
       if a[0] > 0 {
           fmt.Println("BitLGM")
       } else {
           fmt.Println("BitAryo")
       }
   case 2:
       sort.Ints(a)
       // Wythoff Nim losing positions: a = floor(k*phi), b = a + k
       phi := (1 + math.Sqrt(5)) / 2
       k := a[1] - a[0]
       t := int(math.Floor(float64(k) * phi))
       if t == a[0] {
           fmt.Println("BitAryo")
       } else {
           fmt.Println("BitLGM")
       }
   case 3:
       x := a[0] ^ a[1] ^ a[2]
       if x != 0 {
           fmt.Println("BitLGM")
       } else {
           fmt.Println("BitAryo")
       }
   default:
       // Should not happen for n<=3
       fmt.Println("BitLGM")
   }
}

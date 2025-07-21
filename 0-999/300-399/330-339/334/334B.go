package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var x, y int
   points := make(map[[2]int]bool)
   xs := make(map[int]bool)
   ys := make(map[int]bool)
   for i := 0; i < 8; i++ {
       if _, err := fmt.Fscan(reader, &x, &y); err != nil {
           fmt.Println("ugly")
           return
       }
       p := [2]int{x, y}
       points[p] = true
       xs[x] = true
       ys[y] = true
   }
   // Must have 8 unique points
   if len(points) != 8 || len(xs) != 3 || len(ys) != 3 {
       fmt.Println("ugly")
       return
   }
   // Sort unique x and y values
   sx := make([]int, 0, 3)
   sy := make([]int, 0, 3)
   for xv := range xs {
       sx = append(sx, xv)
   }
   for yv := range ys {
       sy = append(sy, yv)
   }
   sort.Ints(sx)
   sort.Ints(sy)
   // Check grid minus center
   for i := 0; i < 3; i++ {
       for j := 0; j < 3; j++ {
           if i == 1 && j == 1 {
               continue
           }
           if !points[[2]int{sx[i], sy[j]}] {
               fmt.Println("ugly")
               return
           }
       }
   }
   fmt.Println("respectable")
}

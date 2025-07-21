package main

import (
   "fmt"
   "os"
   "sort"
)

func main() {
   var n int
   if _, err := fmt.Fscan(os.Stdin, &n); err != nil {
       return
   }
   xs := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(os.Stdin, &xs[i])
   }
   sort.Ints(xs)
   // piles holds current heights of each pile
   piles := []int{}
   for _, x := range xs {
       // find pile with maximum height h such that h <= x
       bestIdx := -1
       bestH := -1
       for i, h := range piles {
           if h <= x && h > bestH {
               bestH = h
               bestIdx = i
           }
       }
       if bestIdx == -1 {
           // start new pile
           piles = append(piles, 1)
       } else {
           // place on existing pile
           piles[bestIdx]++
       }
   }
   fmt.Println(len(piles))
}

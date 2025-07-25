package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   a := make([]int, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(in, &a[i])
   }
   // sort descending for better block allocation
   sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
   // count of 4-seat blocks (middle) and 2-seat blocks (ends)
   count4 := n
   count2 := 2 * n
   for _, ai := range a {
       // try to allocate ai seats using x four-blocks and y two-blocks
       // minimize x+y (blocks used), tie-breaking by larger x (use more 4-blocks)
       bestCost := int(1e18)
       bestX, bestY := 0, 0
       // maximum x we need to try is ceil(ai/4)
       maxX := ai/4 + 1
       if maxX > count4 {
           maxX = count4
       }
       for x := 0; x <= maxX; x++ {
           rem := ai - 4*x
           var y int
           if rem > 0 {
               y = (rem + 1) / 2
           } else {
               y = 0
           }
           if y > count2 {
               continue
           }
           cost := x + y
           if cost < bestCost || (cost == bestCost && x > bestX) {
               bestCost = cost
               bestX = x
               bestY = y
           }
       }
       if bestCost > count4+count2 { // no valid allocation
           fmt.Println("NO")
           return
       }
       // use blocks
       count4 -= bestX
       count2 -= bestY
   }
   fmt.Println("YES")
}

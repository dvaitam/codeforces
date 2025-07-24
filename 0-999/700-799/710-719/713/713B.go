package main

import (
   "bufio"
   "fmt"
   "os"
)

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

func query(x1, y1, x2, y2 int) int {
   fmt.Fprintf(writer, "? %d %d %d %d\n", x1, y1, x2, y2)
   writer.Flush()
   var res int
   fmt.Fscan(reader, &res)
   return res
}

// findRect finds the bounds of a single rectangle fully inside [x1,y1,x2,y2]
func findRect(x1, y1, x2, y2 int) (rx1, ry1, rx2, ry2 int) {
   // find rx1: maximal low bound such that [rx1, x2] contains rectangle
   low, high := x1, x2
   for low <= high {
       mid := (low + high) / 2
       if query(mid, y1, x2, y2) >= 1 {
           rx1 = mid
           low = mid + 1
       } else {
           high = mid - 1
       }
   }
   // find rx2: minimal high bound such that [x1, rx2] contains rectangle
   low, high = x1, x2
   var ans2 int
   for low <= high {
       mid := (low + high) / 2
       if query(x1, y1, mid, y2) >= 1 {
           ans2 = mid
           high = mid - 1
       } else {
           low = mid + 1
       }
   }
   rx2 = ans2
   // find ry1: maximal low bound for y
   low, high = y1, y2
   for low <= high {
       mid := (low + high) / 2
       if query(x1, mid, x2, y2) >= 1 {
           ry1 = mid
           low = mid + 1
       } else {
           high = mid - 1
       }
   }
   // find ry2: minimal high bound for y
   low, high = y1, y2
   var ansY2 int
   for low <= high {
       mid := (low + high) / 2
       if query(x1, y1, x2, mid) >= 1 {
           ansY2 = mid
           high = mid - 1
       } else {
           low = mid + 1
       }
   }
   ry2 = ansY2
   return
}

func main() {
   defer writer.Flush()
   var n int
   fmt.Fscan(reader, &n)
   // try vertical split
   split := -1
   for low, high := 1, n-1; low <= high; {
       mid := (low + high) / 2
       cnt := query(1, 1, mid, n)
       if cnt == 1 {
           split = mid
           break
       } else if cnt == 0 {
           low = mid + 1
       } else {
           high = mid - 1
       }
   }
   var r1, r2 [4]int
   if split != -1 {
       r1[0], r1[1], r1[2], r1[3] = findRect(1, 1, split, n)
       r2[0], r2[1], r2[2], r2[3] = findRect(split+1, 1, n, n)
   } else {
       // horizontal split
       splitY := -1
       for low, high := 1, n-1; low <= high; {
           mid := (low + high) / 2
           cnt := query(1, 1, n, mid)
           if cnt == 1 {
               splitY = mid
               break
           } else if cnt == 0 {
               low = mid + 1
           } else {
               high = mid - 1
           }
       }
       r1[0], r1[1], r1[2], r1[3] = findRect(1, 1, n, splitY)
       r2[0], r2[1], r2[2], r2[3] = findRect(1, splitY+1, n, n)
   }
   // output answer
   fmt.Fprintf(writer, "! %d %d %d %d %d %d %d %d\n",
       r1[0], r1[1], r1[2], r1[3], r2[0], r2[1], r2[2], r2[3])
   writer.Flush()
}

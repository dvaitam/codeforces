package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

func query(x1, y1, x2, y2 int) int {
   fmt.Fprintf(writer, "? %d %d %d %d\n", x1, y1, x2, y2)
   writer.Flush()
   var res int
   fmt.Fscan(reader, &res)
   return res
}

// findSegment finds the two endpoints' sorted coordinates (low, high) in one dimension
// when querying rectangles from origin to mid along that dimension.
// If no odd-parity segment, returns (-1, -1).
func findSegment(n int, horizontal bool) (int, int) {
   // find lowest mid with odd parity
   lo, hi := 1, n
   for lo < hi {
       mid := (lo + hi) / 2
       var ans int
       if horizontal {
           ans = query(1, 1, mid, n)
       } else {
           ans = query(1, 1, n, mid)
       }
       if ans%2 == 1 {
           hi = mid
       } else {
           lo = mid + 1
       }
   }
   // check if exists
   var check int
   if horizontal {
       check = query(1, 1, lo, n)
   } else {
       check = query(1, 1, n, lo)
   }
   if check%2 == 0 {
       return -1, -1
   }
   low := lo
   // find highest mid with odd parity
   lo2, hi2 := 1, n
   for lo2 < hi2 {
       mid := (lo2 + hi2 + 1) / 2
       var ans int
       if horizontal {
           ans = query(1, 1, mid, n)
       } else {
           ans = query(1, 1, n, mid)
       }
       if ans%2 == 1 {
           lo2 = mid
       } else {
           hi2 = mid - 1
       }
   }
   // high endpoint is last odd + 1
   high := lo2 + 1
   return low, high
}

// find in a fixed row the column of the endpoint (parity flips at that column)
func findColInRow(row, n int) int {
   lo, hi := 1, n
   for lo < hi {
       mid := (lo + hi) / 2
       ans := query(row, 1, row, mid)
       if ans%2 == 1 {
           hi = mid
       } else {
           lo = mid + 1
       }
   }
   return lo
}

// find in a fixed column the row of the endpoint
func findRowInCol(col, n int) int {
   lo, hi := 1, n
   for lo < hi {
       mid := (lo + hi) / 2
       ans := query(1, col, mid, col)
       if ans%2 == 1 {
           hi = mid
       } else {
           lo = mid + 1
       }
   }
   return lo
}

// check if a single cell is an endpoint (degree 1)
func isEndpoint(r, c int) bool {
   ans := query(r, c, r, c)
   return ans%2 == 1
}

func main() {
   var n int
   fmt.Fscan(reader, &n)
   // find row and column segments
   r1, r2 := findSegment(n, true)
   c1, c2 := findSegment(n, false)
   var x1, y1, x2, y2 int
   if r1 != -1 {
       // different rows
       x1, x2 = r1, r2
       if c1 != -1 {
           // different columns
           y1, y2 = c1, c2
       } else {
           // same column, find per row
           y1 = findColInRow(x1, n)
           y2 = findColInRow(x2, n)
       }
       // ensure correct pairing
       if !isEndpoint(x1, y1) {
           y1, y2 = y2, y1
       }
   } else {
       // same row
       y1, y2 = c1, c2
       x1 = findRowInCol(y1, n)
       x2 = findRowInCol(y2, n)
       if !isEndpoint(x1, y1) {
           x1, x2 = x2, x1
       }
   }
   fmt.Fprintf(writer, "! %d %d %d %d\n", x1, y1, x2, y2)
   writer.Flush()
}

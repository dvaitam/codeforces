package main

import (
   "bufio"
   "fmt"
   "os"
)

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

// query sends a query (x, y) and returns true if response is "x", false if "y".
func query(x, y int64) (bool, error) {
   fmt.Fprintf(writer, "? %d %d\n", x, y)
   writer.Flush()
   var resp string
   if _, err := fmt.Fscan(reader, &resp); err != nil {
       return false, err
   }
   if len(resp) > 0 && resp[0] == 'x' {
       return true, nil
   }
   return false, nil
}

// solve one game: interactive find a
func solve() (int64, error) {
   // special check for a == 1
   ok, err := query(0, 1)
   if err != nil {
       return 0, err
   }
   if ok {
       // a divides 1 => a == 1
       return 1, nil
   }
   // exponential search: find k such that k < a <= 2*k
   var k int64 = 1
   for {
       ok, err = query(k, 2*k)
       if err != nil {
           return 0, err
       }
       if ok {
           break
       }
       k <<= 1
   }
   // binary search in (k, 2*k]
   low, high := k+1, 2*k
   for low < high {
       mid := (low + high) / 2
       ok, err = query(mid, mid+k)
       if err != nil {
           return 0, err
       }
       if ok {
           high = mid
       } else {
           low = mid + 1
       }
   }
   return low, nil
}

func main() {
   defer writer.Flush()
   for {
       a, err := solve()
       if err != nil {
           // no more games
           break
       }
       fmt.Fprintf(writer, "! %d\n", a)
       writer.Flush()
   }
}

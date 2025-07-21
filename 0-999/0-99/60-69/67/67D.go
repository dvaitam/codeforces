package main

import (
   "bufio"
   "fmt"
   "os"
)

var reader = bufio.NewReader(os.Stdin)

// readInt reads next integer from stdin
func readInt() (int, error) {
   var neg bool
   var b byte
   var err error
   // skip non-digit
   for {
       b, err = reader.ReadByte()
       if err != nil {
           return 0, err
       }
       if (b >= '0' && b <= '9') || b == '-' {
           break
       }
   }
   if b == '-' {
       neg = true
       b, err = reader.ReadByte()
       if err != nil {
           return 0, err
       }
   }
   var x int
   for ; b >= '0' && b <= '9'; b, err = reader.ReadByte() {
       if err != nil {
           break
       }
       x = x*10 + int(b - '0')
   }
   if neg {
       x = -x
   }
   return x, nil
}

func main() {
   // read n
   n, err := readInt()
   if err != nil {
       return
   }
   x := make([]int, n)
   for i := 0; i < n; i++ {
       xi, _ := readInt()
       x[i] = xi
   }
   // exit positions
   exitPos := make([]int, n+1)
   for pos := 1; pos <= n; pos++ {
       yi, _ := readInt()
       if yi >= 1 && yi <= n {
           exitPos[yi] = pos
       }
   }
   // compute longest decreasing subsequence of exitPos[x[i]]
   // map to increasing by B = (n+1 - exit)
   tails := make([]int, 0, 16)
   for i := 0; i < n; i++ {
       e := exitPos[x[i]]
       b := n + 1 - e
       // binary search lower_bound in tails for b
       lo, hi := 0, len(tails)
       for lo < hi {
           mid := (lo + hi) >> 1
           if tails[mid] < b {
               lo = mid + 1
           } else {
               hi = mid
           }
       }
       if lo == len(tails) {
           tails = append(tails, b)
       } else {
           tails[lo] = b
       }
   }
   // answer
   fmt.Fprint(os.Stdout, len(tails))
}

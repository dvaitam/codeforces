package main

import (
   "bufio"
   "fmt"
   "os"
)

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

// query compares a[i] and a[j] and returns '<' if a[i]<a[j],
// '>' if a[i]>a[j], '=' if equal.
func query(i, j int) byte {
   fmt.Fprintf(writer, "? %d %d\n", i, j)
   writer.Flush()
   var resp string
   if _, err := fmt.Fscan(reader, &resp); err != nil {
       os.Exit(1)
   }
   return resp[0]
}

func main() {
   // read n
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // trivial case
   if n == 1 {
       fmt.Fprintf(writer, "! 1 1\n")
       writer.Flush()
       return
   }
   // prepare candidate lists
   minCands := make([]int, 0, n/2+1)
   maxCands := make([]int, 0, n/2+1)
   // initial pairing
   for i := 1; i+1 <= n; i += 2 {
       j := i + 1
       r := query(i, j)
       if r == '<' {
           minCands = append(minCands, i)
           maxCands = append(maxCands, j)
       } else {
           minCands = append(minCands, j)
           maxCands = append(maxCands, i)
       }
   }
   // if odd, include last element in both
   if n%2 == 1 {
       k := n
       minCands = append(minCands, k)
       maxCands = append(maxCands, k)
   }
   // find global minimum
   minIdx := minCands[0]
   for _, idx := range minCands[1:] {
       if query(idx, minIdx) == '<' {
           minIdx = idx
       }
   }
   // find global maximum
   maxIdx := maxCands[0]
   for _, idx := range maxCands[1:] {
       if query(idx, maxIdx) == '>' {
           maxIdx = idx
       }
   }
   // output result
   fmt.Fprintf(writer, "! %d %d\n", minIdx, maxIdx)
   writer.Flush()
}

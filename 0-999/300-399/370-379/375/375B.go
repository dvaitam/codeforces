package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       fmt.Fprintln(os.Stderr, "failed to read n and m:", err)
       return
   }
   colCount := make([]int, m)
   // Read each row and count ones per column
   row := make([]byte, m)
   for i := 0; i < n; i++ {
       // read exactly m non-newline bytes into row
       var read int
       for read < m {
           b, err := reader.ReadByte()
           if err != nil {
               fmt.Fprintln(os.Stderr, "failed to read row data:", err)
               return
           }
           if b == '\n' || b == '\r' {
               continue
           }
           row[read] = b
           read++
       }
       for j := 0; j < m; j++ {
           if row[j] == '1' {
               colCount[j]++
           }
       }
   }
   // freq[k] = number of columns with exactly k ones
   freq := make([]int, n+1)
   for _, cnt := range colCount {
       if cnt >= 0 && cnt <= n {
           freq[cnt]++
       }
   }
   // compute max area: for each possible height h, count eligible columns (suffix sum)
   maxArea := 0
   eligible := 0
   for h := n; h >= 1; h-- {
       eligible += freq[h]
       area := h * eligible
       if area > maxArea {
           maxArea = area
       }
   }
   fmt.Println(maxArea)
}

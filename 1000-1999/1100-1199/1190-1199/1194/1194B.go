package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var q int
   if _, err := fmt.Fscan(reader, &q); err != nil {
       return
   }
   for qi := 0; qi < q; qi++ {
       var n, m int
       fmt.Fscan(reader, &n, &m)
       grid := make([][]byte, n)
       for i := 0; i < n; i++ {
           line, _ := reader.ReadBytes('\n')
           // remove newline
           if len(line) > 0 && (line[len(line)-1] == '\n' || line[len(line)-1] == '\r') {
               // trim newline and possible carriage return
               end := len(line) - 1
               if end > 0 && line[end-1] == '\r' {
                   end--
               }
               line = line[:end]
           }
           grid[i] = line
       }
       rowWhite := make([]int, n)
       colWhite := make([]int, m)
       for i := 0; i < n; i++ {
           for j := 0; j < m; j++ {
               if grid[i][j] == '.' {
                   rowWhite[i]++
                   colWhite[j]++
               }
           }
       }
       ans := n*m + 5
       for i := 0; i < n; i++ {
           for j := 0; j < m; j++ {
               need := rowWhite[i] + colWhite[j]
               if grid[i][j] == '.' {
                   need--
               }
               if need < ans {
                   ans = need
               }
           }
       }
       fmt.Fprintln(writer, ans)
   }
}

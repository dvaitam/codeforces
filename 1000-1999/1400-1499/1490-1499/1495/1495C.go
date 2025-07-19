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
   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       solve(reader, writer)
   }
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
   var n, m int
   fmt.Fscan(reader, &n, &m)
   grid := make([][]rune, n)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       grid[i] = []rune(s)
   }

   // Determine starting column in 0-based index
   start := m%3 - 1
   if m%3 == 0 {
       start = 1
   }
   for i := start; i < m; i += 3 {
       // Fill entire column i with 'X'
       for j := 0; j < n; j++ {
           grid[j][i] = 'X'
       }
       // Connect with previous block if not the first stripe
       if i >= 2 {
           ok := false
           for j := 0; j < n; j++ {
               if grid[j][i-1] == 'X' {
                   grid[j][i-2] = 'X'
                   ok = true
                   break
               }
               if grid[j][i-2] == 'X' {
                   grid[j][i-1] = 'X'
                   ok = true
                   break
               }
           }
           if !ok {
               grid[0][i-1] = 'X'
               grid[0][i-2] = 'X'
           }
       }
   }

   // Output the grid
   for i := 0; i < n; i++ {
       writer.WriteString(string(grid[i]))
       writer.WriteByte('\n')
   }
}

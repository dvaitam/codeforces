package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   // read number of test cases
   line, _ := reader.ReadString('\n')
   t, _ := strconv.Atoi(line[:len(line)-1])
   for i := 0; i < t; i++ {
       // read next line
       line, _ = reader.ReadString('\n')
       // split fields
       fields := make([]int64, 6)
       var idx, start int
       for j := 0; j < 6; j++ {
           // skip spaces
           for start < len(line) && (line[start] == ' ' || line[start] == '\n' || line[start] == '\r' || line[start] == '\t') {
               start++
           }
           idx = start
           for start < len(line) && line[start] >= '0' && line[start] <= '9' {
               start++
           }
           val, _ := strconv.ParseInt(line[idx:start], 10, 64)
           fields[j] = val
       }
       n, m := fields[0], fields[1]
       x1, y1 := fields[2], fields[3]
       x2, y2 := fields[4], fields[5]

       // compute absolute differences
       dx := x1 - x2
       if dx < 0 {
           dx = -dx
       }
       dy := y1 - y2
       if dy < 0 {
           dy = -dy
       }
       // number of positions for the pair in x and y
       W := n - dx
       H := m - dy
       // intersection overlap
       overlapW := n - 2*dx
       if overlapW < 0 {
           overlapW = 0
       }
       overlapH := m - 2*dy
       if overlapH < 0 {
           overlapH = 0
       }
       // total melted cells
       // melted = 2*W*H - overlapW*overlapH
       melted := 2*W*H - overlapW*overlapH
       total := n * m
       unmelted := total - melted
       fmt.Fprintln(writer, unmelted)
   }
}

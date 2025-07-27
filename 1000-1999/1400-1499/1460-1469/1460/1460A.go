package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   line, err := reader.ReadString('\n')
   if err != nil {
       return
   }
   parts := strings.Fields(line)
   if len(parts) < 2 {
       return
   }
   n, _ := strconv.Atoi(parts[0])
   m, _ := strconv.Atoi(parts[1])
   // skip server specs
   for i := 0; i < n; i++ {
       _, _ = reader.ReadString('\n')
   }
   // skip vm specs
   for j := 0; j < m; j++ {
       _, _ = reader.ReadString('\n')
   }
   oldSrv := make([]int, m)
   newSrv := make([]int, m)
   for j := 0; j < m; j++ {
       ln, err := reader.ReadString('\n')
       if err != nil {
           break
       }
       f := strings.Fields(ln)
       if len(f) >= 2 {
           oldSrv[j], _ = strconv.Atoi(f[0])
           newSrv[j], _ = strconv.Atoi(f[1])
       }
   }
   // collect moves
   moves := make([][3]int, 0, m)
   for j := 0; j < m; j++ {
       if oldSrv[j] != newSrv[j] {
           moves = append(moves, [3]int{oldSrv[j], newSrv[j], j})
       }
   }
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   S := len(moves)
   fmt.Fprintln(w, S)
   for _, mv := range moves {
       // one move per step
       fmt.Fprintln(w, 1)
       fmt.Fprintf(w, "%d %d %d\n", mv[0], mv[1], mv[2])
   }
}

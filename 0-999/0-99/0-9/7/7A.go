package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   board := make([]string, 8)
   for i := 0; i < 8; i++ {
       line, err := reader.ReadString('\n')
       if err != nil && len(line) == 0 {
           fmt.Fprintln(os.Stderr, "error reading input:", err)
           return
       }
       if len(line) > 0 && line[len(line)-1] == '\n' {
           line = line[:len(line)-1]
       }
       board[i] = line
   }
   // Determine rows and columns that are fully black (no W)
   candidateRow := make([]bool, 8)
   candidateCol := make([]bool, 8)
   for i := 0; i < 8; i++ {
       ok := true
       for j := 0; j < 8; j++ {
           if board[i][j] == 'W' {
               ok = false
               break
           }
       }
       candidateRow[i] = ok
   }
   for j := 0; j < 8; j++ {
       ok := true
       for i := 0; i < 8; i++ {
           if board[i][j] == 'W' {
               ok = false
               break
           }
       }
       candidateCol[j] = ok
   }
   // Forced selections
   forcedRow := make([]bool, 8)
   forcedCol := make([]bool, 8)
   // Remaining edges between candidate rows and cols
   adj := make([][]int, 8)
   for i := 0; i < 8; i++ {
       for j := 0; j < 8; j++ {
           if board[i][j] != 'B' {
               continue
           }
           rOk := candidateRow[i]
           cOk := candidateCol[j]
           if !rOk && cOk {
               forcedCol[j] = true
           } else if rOk && !cOk {
               forcedRow[i] = true
           } else if rOk && cOk {
               // candidate both sides: possible edge
               adj[i] = append(adj[i], j)
           }
           // if neither is candidate, impossible per problem guarantee
       }
   }
   // Build graph excluding forced nodes
   // Prepare matching for columns
   matchCol := make([]int, 8)
   for j := 0; j < 8; j++ {
       matchCol[j] = -1
   }
   var dfs func(int, []bool) bool
   dfs = func(i int, seen []bool) bool {
       for _, j := range adj[i] {
           if forcedCol[j] {
               continue
           }
           if seen[j] {
               continue
           }
           seen[j] = true
           if matchCol[j] == -1 || dfs(matchCol[j], seen) {
               matchCol[j] = i
               return true
           }
       }
       return false
   }
   // Count matching size on non-forced candidate rows
   matchSize := 0
   for i := 0; i < 8; i++ {
       if !candidateRow[i] || forcedRow[i] {
           continue
       }
       seen := make([]bool, 8)
       if dfs(i, seen) {
           matchSize++
       }
   }
   // Count forced selections
   total := matchSize
   for i := 0; i < 8; i++ {
       if forcedRow[i] {
           total++
       }
       if forcedCol[i] {
           total++
       }
   }
   fmt.Println(total)
}

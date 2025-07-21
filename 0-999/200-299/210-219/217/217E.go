package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // Read initial DNA string
   s, _ := reader.ReadString('\n')
   // strip newline
   if len(s) > 0 && (s[len(s)-1] == '\n' || s[len(s)-1] == '\r') {
       s = s[:len(s)-1]
   }
   // Read k
   line, _ := reader.ReadString('\n')
   k, _ := strconv.Atoi(string(line[:len(line)-1]))
   // Read n
   line, _ = reader.ReadString('\n')
   n, _ := strconv.Atoi(string(line[:len(line)-1]))
   L := make([]int, n)
   R := make([]int, n)
   segLen := make([]int, n)
   half := make([]int, n)
   // Read operations
   for i := 0; i < n; i++ {
       line, _ = reader.ReadString('\n')
       // parse li ri
       parts := []byte(line)
       // find space
       sp := 0
       for sp < len(parts) && parts[sp] != ' ' {
           sp++
       }
       li, _ := strconv.Atoi(string(parts[:sp]))
       ri, _ := strconv.Atoi(string(parts[sp+1 : len(parts)-1]))
       L[i] = li
       R[i] = ri
       ln := ri - li + 1
       segLen[i] = ln
       half[i] = ln / 2
   }
   // Prepare output buffer
   out := make([]byte, k)
   // For each position in 1..k, map back
   for pos := 1; pos <= k; pos++ {
       p := pos
       // apply operations in reverse
       for i := n - 1; i >= 0; i-- {
           li := L[i]
           ri := R[i]
           ln := segLen[i]
           // inserted range [ri+1, ri+ln]
           if p > ri && p <= ri+ln {
               idx := p - ri
               h := half[i]
               if idx <= h {
                   // even part
                   p = li + idx*2 - 1
               } else {
                   // odd part
                   p = li + (idx-h-1)*2
               }
           } else if p > ri+ln {
               p -= ln
           }
       }
       // now p in original string
       out[pos-1] = s[p-1]
   }
   // write output
   writer := bufio.NewWriter(os.Stdout)
   writer.Write(out)
   writer.WriteByte('\n')
   writer.Flush()
}

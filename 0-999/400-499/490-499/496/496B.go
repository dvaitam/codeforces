package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // Read the next line and extract digits
   line, _ := reader.ReadString('\n')
   line = strings.TrimSpace(line)
   // Remove spaces if any
   line = strings.ReplaceAll(line, " ", "")
   s := line
   if len(s) != n {
       // maybe no newline consumed, try reading again tokens
       s = ""
       for len(s) < n {
           var d string
           if _, err := fmt.Fscan(reader, &d); err != nil {
               break
           }
           s += d
       }
   }
   best := strings.Repeat("9", n)
   // compare best with default initial highest
   // function to compare numeric strings a and b
   cmp := func(a, b string) int {
       // strip leading zeros
       ia := 0
       for ia < len(a) && a[ia] == '0' {
           ia++
       }
       ib := 0
       for ib < len(b) && b[ib] == '0' {
           ib++
       }
       ra := a[ia:]
       rb := b[ib:]
       // if both empty, equal
       if len(ra) == 0 && len(rb) == 0 {
           return 0
       }
       if len(ra) != len(rb) {
           if len(ra) < len(rb) {
               return -1
           }
           return 1
       }
       if ra < rb {
           return -1
       }
       if ra > rb {
           return 1
       }
       return 0
   }
   nrun := n
   // Precompute digits as ints 0-9
   digits := make([]int, n)
   for i := 0; i < n; i++ {
       digits[i] = int(s[i] - '0')
   }
   // Try all rotations
   for rot := 0; rot < nrun; rot++ {
       // Build rotated slice where pos i maps to digits[(i-rot+n)%n]
       rotDigits := make([]int, nrun)
       for i := 0; i < nrun; i++ {
           idx := (i - rot + nrun) % nrun
           rotDigits[i] = digits[idx]
       }
       // Try all increments x = 0..9
       for x := 0; x < 10; x++ {
           // build candidate string
           b := make([]byte, nrun)
           for i := 0; i < nrun; i++ {
               v := (rotDigits[i] + x) % 10
               b[i] = byte(v) + '0'
           }
           cand := string(b)
           if cmp(cand, best) < 0 {
               best = cand
           }
       }
   }
   // output best
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, best)
}

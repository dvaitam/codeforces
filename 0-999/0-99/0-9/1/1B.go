package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func isRxCy(s string) bool {
   // check pattern R<digits>C<digits>
   if len(s) < 4 || s[0] != 'R' {
       return false
   }
   i := 1
   // at least one digit after R
   for i < len(s) && s[i] >= '0' && s[i] <= '9' {
       i++
   }
   // must have consumed at least one digit, and next char is 'C'
   if i == 1 || i >= len(s) || s[i] != 'C' {
       return false
   }
   // after C, at least one digit
   j := i + 1
   if j >= len(s) {
       return false
   }
   for j < len(s) && s[j] >= '0' && s[j] <= '9' {
       j++
   }
   return j == len(s)
}

func colToLetters(col int) string {
   // convert 1-based column number to letters
   res := make([]byte, 0)
   for col > 0 {
       col--
       r := col % 26
       res = append(res, byte('A'+r))
       col /= 26
   }
   // reverse
   for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
       res[i], res[j] = res[j], res[i]
   }
   return string(res)
}

func lettersToCol(s string) int {
   col := 0
   for i := 0; i < len(s); i++ {
       col = col*26 + int(s[i]-'A'+1)
   }
   return col
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   // read n
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   for ; n > 0; n-- {
       var s string
       fmt.Fscan(reader, &s)
       if isRxCy(s) {
           // R<row>C<col> to letters+row
           // find 'C'
           // strip 'R'
           i := 1
           for s[i] >= '0' && s[i] <= '9' {
               i++
           }
           row := s[1:i]
           colNum, _ := strconv.Atoi(s[i+1:])
           letters := colToLetters(colNum)
           fmt.Fprintf(writer, "%s%s", letters, row)
       } else {
           // letters+row to R<row>C<col>
           // split letters and digits
           i := 0
           for i < len(s) && s[i] >= 'A' && s[i] <= 'Z' {
               i++
           }
           letters := s[:i]
           row := s[i:]
           colNum := lettersToCol(letters)
           fmt.Fprintf(writer, "R%sC%d", row, colNum)
       }
       if n > 1 {
           writer.WriteByte('\n')
       }
   }
}

package main

import (
   "bufio"
   "fmt"
   "os"
   "time"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var dateStr string
   var shift int
   // Read date string and shift value
   if _, err := fmt.Fscan(reader, &dateStr); err != nil {
       return
   }
   if _, err := fmt.Fscan(reader, &shift); err != nil {
       return
   }
   // Parse date in DD.MM.YYYY format
   t, err := time.Parse("02.01.2006", dateStr)
   if err != nil {
       return
   }
   // Add shift days
   t = t.AddDate(0, 0, shift)
   // Output result in same format
   fmt.Println(t.Format("02.01.2006"))
}

package main
import (
   "bufio"
   "fmt"
   "os"
   "math/bits"
)
func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   allowed := make([]int, n)
   fullMask := (1 << 6) - 1
   for i := 0; i < n; i++ {
       allowed[i] = fullMask
   }
   var m int
   fmt.Fscan(reader, &m)
   for i := 0; i < m; i++ {
       var pos int
       var str string
       fmt.Fscan(reader, &pos, &str)
       pos--
       mask := 0
       for j := 0; j < len(str); j++ {
           mask |= 1 << int(str[j]-'a')
       }
       if pos >= 0 && pos < n {
           allowed[pos] = mask
       }
   }

   // count letters in s
   cou := make([]int, 6)
   for i := 0; i < n; i++ {
       cou[s[i]-'a']++
   }

   // compute left[mask]: number of letters available in mask
   size := 1 << 6
   left := make([]int, size)
   for mask := 0; mask < size; mask++ {
       cnt := 0
       for bit := 0; bit < 6; bit++ {
           if mask>>bit&1 == 1 {
               cnt += cou[bit]
           }
       }
       left[mask] = cnt
   }

   // compute matching[mask]: positions that can accept letters in mask
   matching := make([]int, size)
   for i := 0; i < n; i++ {
       a := allowed[i]
       for mask := 1; mask < size; mask++ {
           if a&mask != 0 {
               matching[mask]++
           }
       }
   }

   // check feasibility
   for mask := 0; mask < size; mask++ {
       if left[mask] > matching[mask] {
           fmt.Fprintln(writer, "Impossible")
           return
       }
   }

   ans := make([]byte, n)
   // assign greedily
   for i := 0; i < n; i++ {
       poss := allowed[i]
       for mask := 1; mask < size; mask++ {
           if left[mask] == matching[mask] && (mask&allowed[i]) != 0 {
               poss &= mask
           }
       }
       if poss == 0 {
           fmt.Fprintln(writer, "Impossible")
           return
       }
       take := bits.TrailingZeros(uint(poss))
       ans[i] = byte('a' + take)
       // update counts
       for mask := 0; mask < size; mask++ {
           if mask>>take&1 == 1 {
               left[mask]--
           }
           if mask&allowed[i] != 0 {
               matching[mask]--
           }
       }
   }

   fmt.Fprintln(writer, string(ans))
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int64
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   var x, y string
   fmt.Fscan(reader, &x)
   fmt.Fscan(reader, &y)
   px, py := len(x), len(y)
   g := gcd(px, py)
   // Number of small-block repetitions for y per cycle
   ty := int64(py / g)
   // Number of full cycles over which each pair (i,j) aligns
   fullCycles := n / ty

   var totalPairsPerCycle int64
   // For each residue class modulo g
   for r := 0; r < g; r++ {
       var cntX [26]int64
       for i := r; i < px; i += g {
           cntX[x[i]-'a']++
       }
       var cntY [26]int64
       for j := r; j < py; j += g {
           cntY[y[j]-'a']++
       }
       for c := 0; c < 26; c++ {
           totalPairsPerCycle += cntX[c] * cntY[c]
       }
   }

   // Total matches over all cycles
   totalMatches := totalPairsPerCycle * fullCycles
   // Total positions in the concatenated string a
   totalPositions := n * int64(px)
   // Hamming distance is mismatches
   result := totalPositions - totalMatches
   fmt.Println(result)
}

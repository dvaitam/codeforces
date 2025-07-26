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

    var s string
    if _, err := fmt.Fscan(reader, &s); err != nil {
        return
    }
    // digits are 1-indexed
    digits := make([]int, 101)
    for i := 0; i < 100 && i < len(s); i++ {
        digits[i+1] = int(s[i] - '0')
    }

    // Apply operations
    if digits[39]%2 == 1 { digits[39] = (digits[39] + 9) % 10 } else { digits[37] = (digits[37] + 1) % 10 }
    if digits[24]%2 == 1 { digits[24] = (digits[24] + 1) % 10 } else { digits[76] = (digits[76] + 3) % 10 }
    if digits[13] + digits[91] > 10 { digits[14] = (digits[14] + 6) % 10 } else { digits[34] = (digits[34] + 8) % 10 }
    if digits[87]%2 == 1 { digits[87] = (digits[87] + 7) % 10 } else { digits[22] = (digits[22] + 9) % 10 }
    if digits[79] > digits[15] { digits[74] = (digits[74] + 7) % 10 } else { digits[84] = (digits[84] + 6) % 10 }
    if digits[26] + digits[66] > 9 { digits[31] = (digits[31] + 7) % 10 } else { digits[95] = (digits[95] + 4) % 10 }
    if digits[53] + digits[1] > 8 { digits[66] = (digits[66] + 1) % 10 } else { digits[94] = (digits[94] + 6) % 10 }
    if digits[41] > digits[29] { digits[67] = (digits[67] + 5) % 10 } else { digits[41] = (digits[41] + 9) % 10 }
    if digits[79] + digits[20] > 10 { digits[18] = (digits[18] + 2) % 10 } else { digits[72] = (digits[72] + 9) % 10 }
    if digits[14] + digits[24] > 10 { digits[64] = (digits[64] + 2) % 10 } else { digits[84] = (digits[84] + 2) % 10 }
    if digits[16] > digits[34] { digits[81] = (digits[81] + 5) % 10 } else { digits[15] = (digits[15] + 2) % 10 }
    if digits[48] + digits[65] > 9 { digits[57] = (digits[57] + 2) % 10 } else { digits[28] = (digits[28] + 5) % 10 }
    if digits[81]%2 == 1 { digits[81] = (digits[81] + 5) % 10 } else { digits[25] = (digits[25] + 4) % 10 }
    if digits[70]%2 == 1 { digits[70] = (digits[70] + 9) % 10 } else { digits[93] = (digits[93] + 3) % 10 }
    if digits[92] + digits[49] > 9 { digits[81] = (digits[81] + 2) % 10 } else { digits[42] = (digits[42] + 3) % 10 }
    if digits[96] > digits[20] { digits[45] = (digits[45] + 4) % 10 } else { digits[45] = (digits[45] + 1) % 10 }
    if digits[91] > digits[21] { digits[60] = (digits[60] + 3) % 10 } else { digits[72] = (digits[72] + 1) % 10 }
    if digits[89] > digits[7] { digits[98] = (digits[98] + 9) % 10 } else { digits[52] = (digits[52] + 7) % 10 }
    if digits[38] > digits[97] { digits[92] = (digits[92] + 6) % 10 } else { digits[35] = (digits[35] + 4) % 10 }
    if digits[96] > digits[99] { digits[42] = (digits[42] + 4) % 10 } else { digits[40] = (digits[40] + 9) % 10 }
    if digits[86]%2 == 1 { digits[86] = (digits[86] + 1) % 10 } else { digits[14] = (digits[14] + 3) % 10 }
    if digits[23]%2 == 1 { digits[23] = (digits[23] + 5) % 10 } else { digits[55] = (digits[55] + 9) % 10 }
    if digits[79]%2 == 1 { digits[79] = (digits[79] + 1) % 10 } else { digits[29] = (digits[29] + 8) % 10 }
    if digits[4] > digits[91] { digits[98] = (digits[98] + 8) % 10 } else { digits[69] = (digits[69] + 4) % 10 }
    if digits[93] > digits[24] { digits[75] = (digits[75] + 9) % 10 } else { digits[95] = (digits[95] + 3) % 10 }
    if digits[32] + digits[50] > 10 { digits[91] = (digits[91] + 3) % 10 } else { digits[1] = (digits[1] + 5) % 10 }
    if digits[81] > digits[31] { digits[86] = (digits[86] + 7) % 10 } else { digits[67] = (digits[67] + 5) % 10 }
    if digits[83] > digits[86] { digits[48] = (digits[48] + 7) % 10 } else { digits[2] = (digits[2] + 6) % 10 }
    if digits[20] > digits[88] { digits[9] = (digits[9] + 2) % 10 } else { digits[99] = (digits[99] + 4) % 10 }
    if digits[14]%2 == 1 { digits[14] = (digits[14] + 5) % 10 } else { digits[97] = (digits[97] + 7) % 10 }
    if digits[38] > digits[14] { digits[48] = (digits[48] + 2) % 10 } else { digits[81] = (digits[81] + 5) % 10 }
    if digits[92] > digits[74] { digits[92] = (digits[92] + 1) % 10 } else { digits[50] = (digits[50] + 9) % 10 }
    if digits[76] > digits[89] { digits[68] = (digits[68] + 6) % 10 } else { digits[69] = (digits[69] + 5) % 10 }
    if digits[2] > digits[28] { digits[75] = (digits[75] + 1) % 10 } else { digits[89] = (digits[89] + 1) % 10 }
    if digits[67]%2 == 1 { digits[67] = (digits[67] + 9) % 10 } else { digits[49] = (digits[49] + 1) % 10 }
    if digits[23]%2 == 1 { digits[23] = (digits[23] + 1) % 10 } else { digits[59] = (digits[59] + 3) % 10 }
    if digits[81]%2 == 1 { digits[81] = (digits[81] + 9) % 10 } else { digits[9] = (digits[9] + 4) % 10 }
    if digits[92] + digits[82] > 9 { digits[81] = (digits[81] + 2) % 10 } else { digits[91] = (digits[91] + 5) % 10 }
    if digits[42] + digits[48] > 9 { digits[35] = (digits[35] + 8) % 10 } else { digits[59] = (digits[59] + 6) % 10 }
    if digits[55]%2 == 1 { digits[55] = (digits[55] + 9) % 10 } else { digits[61] = (digits[61] + 6) % 10 }
    if digits[83]%2 == 1 { digits[83] = (digits[83] + 5) % 10 } else { digits[85] = (digits[85] + 4) % 10 }
    if digits[96]%2 == 1 { digits[96] = (digits[96] + 1) % 10 } else { digits[72] = (digits[72] + 4) % 10 }
    if digits[17]%2 == 1 { digits[17] = (digits[17] + 1) % 10 } else { digits[28] = (digits[28] + 3) % 10 }
    if digits[85] > digits[74] { digits[37] = (digits[37] + 3) % 10 } else { digits[10] = (digits[10] + 3) % 10 }
    if digits[50] + digits[67] > 9 { digits[85] = (digits[85] + 9) % 10 } else { digits[42] = (digits[42] + 4) % 10 }
    if digits[11] + digits[43] > 10 { digits[56] = (digits[56] + 7) % 10 } else { digits[50] = (digits[50] + 7) % 10 }
    if digits[95] + digits[64] > 9 { digits[95] = (digits[95] + 4) % 10 } else { digits[95] = (digits[95] + 9) % 10 }
    if digits[21] + digits[16] > 9 { digits[87] = (digits[87] + 3) % 10 } else { digits[30] = (digits[30] + 1) % 10 }
    if digits[91]%2 == 1 { digits[91] = (digits[91] + 1) % 10 } else { digits[77] = (digits[77] + 1) % 10 }
    if digits[95] > digits[82] { digits[53] = (digits[53] + 2) % 10 } else { digits[100] = (digits[100] + 5) % 10 }
    if digits[88] + digits[66] > 10 { digits[34] = (digits[34] + 4) % 10 } else { digits[57] = (digits[57] + 4) % 10 }
    if digits[73] > digits[84] { digits[52] = (digits[52] + 3) % 10 } else { digits[42] = (digits[42] + 9) % 10 }
    if digits[66] > digits[38] { digits[94] = (digits[94] + 7) % 10 } else { digits[78] = (digits[78] + 7) % 10 }
    if digits[23] > digits[12] { digits[78] = (digits[78] + 2) % 10 } else { digits[62] = (digits[62] + 8) % 10 }
    if digits[13] > digits[9] { digits[42] = (digits[42] + 7) % 10 } else { digits[1] = (digits[1] + 9) % 10 }
    if digits[43] > digits[29] { digits[20] = (digits[20] + 2) % 10 } else { digits[47] = (digits[47] + 2) % 10 }
    if digits[100] + digits[51] > 8 { digits[10] = (digits[10] + 6) % 10 } else { digits[89] = (digits[89] + 1) % 10 }
    if digits[19] > digits[37] { digits[26] = (digits[26] + 7) % 10 } else { digits[30] = (digits[30] + 8) % 10 }
    if digits[73] > digits[25] { digits[77] = (digits[77] + 3) % 10 } else { digits[41] = (digits[41] + 1) % 10 }
    if digits[67] + digits[96] > 10 { digits[47] = (digits[47] + 6) % 10 } else { digits[33] = (digits[33] + 5) % 10 }
    if digits[11] > digits[10] { digits[33] = (digits[33] + 3) % 10 } else { digits[4] = (digits[4] + 3) % 10 }
    if digits[85]%2 == 1 { digits[85] = (digits[85] + 7) % 10 } else { digits[37] = (digits[37] + 9) % 10 }
    if digits[14]%2 == 1 { digits[14] = (digits[14] + 1) % 10 } else { digits[28] = (digits[28] + 4) % 10 }
    if digits[30] + digits[18] > 8 { digits[93] = (digits[93] + 5) % 10 } else { digits[68] = (digits[68] + 1) % 10 }
    if digits[54] + digits[72] > 8 { digits[88] = (digits[88] + 8) % 10 } else { digits[25] = (digits[25] + 8) % 10 }
    if digits[72]%2 == 1 { digits[72] = (digits[72] + 5) % 10 } else { digits[10] = (digits[10] + 3) % 10 }
    if digits[15]%2 == 1 { digits[15] = (digits[15] + 3) % 10 } else { digits[68] = (digits[68] + 1) % 10 }
    if digits[81] + digits[31] > 9 { digits[2] = (digits[2] + 5) % 10 } else { digits[35] = (digits[35] + 1) % 10 }
    if digits[57]%2 == 1 { digits[57] = (digits[57] + 1) % 10 } else { digits[25] = (digits[25] + 9) % 10 }
    if digits[75] + digits[51] > 9 { digits[73] = (digits[73] + 8) % 10 } else { digits[49] = (digits[49] + 1) % 10 }
    if digits[81] + digits[61] > 10 { digits[61] = (digits[61] + 3) % 10 } else { digits[88] = (digits[88] + 1) % 10 }
    if digits[60]%2 == 1 { digits[60] = (digits[60] + 1) % 10 } else { digits[31] = (digits[31] + 2) % 10 }
    if digits[93]%2 == 1 { digits[93] = (digits[93] + 5) % 10 } else { digits[50] = (digits[50] + 1) % 10 }
    if digits[19] + digits[82] > 9 { digits[48] = (digits[48] + 7) % 10 } else { digits[88] = (digits[88] + 8) % 10 }
    if digits[45]%2 == 1 { digits[45] = (digits[45] + 7) % 10 } else { digits[100] = (digits[100] + 1) % 10 }
    if digits[46] > digits[71] { digits[28] = (digits[28] + 8) % 10 } else { digits[37] = (digits[37] + 6) % 10 }
    if digits[79]%2 == 1 { digits[79] = (digits[79] + 5) % 10 } else { digits[10] = (digits[10] + 1) % 10 }
    if digits[19] > digits[95] { digits[76] = (digits[76] + 9) % 10 } else { digits[95] = (digits[95] + 8) % 10 }
    if digits[49]%2 == 1 { digits[49] = (digits[49] + 5) % 10 } else { digits[66] = (digits[66] + 3) % 10 }
    if digits[62]%2 == 1 { digits[62] = (digits[62] + 1) % 10 } else { digits[26] = (digits[26] + 8) % 10 }
    if digits[67] > digits[33] { digits[27] = (digits[27] + 8) % 10 } else { digits[96] = (digits[96] + 2) % 10 }
    if digits[73] + digits[15] > 8 { digits[98] = (digits[98] + 6) % 10 } else { digits[11] = (digits[11] + 6) % 10 }
    if digits[63] > digits[42] { digits[66] = (digits[66] + 1) % 10 } else { digits[58] = (digits[58] + 2) % 10 }
    if digits[41]%2 == 1 { digits[41] = (digits[41] + 9) % 10 } else { digits[99] = (digits[99] + 5) % 10 }
    if digits[93]%2 == 1 { digits[93] = (digits[93] + 5) % 10 } else { digits[53] = (digits[53] + 1) % 10 }
    if digits[46]%2 == 1 { digits[46] = (digits[46] + 3) % 10 } else { digits[64] = (digits[64] + 4) % 10 }
    if digits[99] + digits[64] > 10 { digits[72] = (digits[72] + 9) % 10 } else { digits[51] = (digits[51] + 5) % 10 }
    if digits[75] > digits[23] { digits[89] = (digits[89] + 2) % 10 } else { digits[76] = (digits[76] + 7) % 10 }
    if digits[6]%2 == 1 { digits[6] = (digits[6] + 1) % 10 } else { digits[44] = (digits[44] + 6) % 10 }
    if digits[58]%2 == 1 { digits[58] = (digits[58] + 3) % 10 } else { digits[49] = (digits[49] + 9) % 10 }
    if digits[5] > digits[13] { digits[46] = (digits[46] + 9) % 10 } else { digits[21] = (digits[21] + 7) % 10 }
    if digits[44] + digits[94] > 9 { digits[36] = (digits[36] + 4) % 10 } else { digits[15] = (digits[15] + 3) % 10 }
    if digits[52] + digits[43] > 8 { digits[29] = (digits[29] + 8) % 10 } else { digits[72] = (digits[72] + 6) % 10 }
    if digits[87] + digits[48] > 9 { digits[61] = (digits[61] + 8) % 10 } else { digits[14] = (digits[14] + 3) % 10 }
    if digits[81]%2 == 1 { digits[81] = (digits[81] + 7) % 10 } else { digits[64] = (digits[64] + 2) % 10 }
    if digits[88]%2 == 1 { digits[88] = (digits[88] + 7) % 10 } else { digits[53] = (digits[53] + 9) % 10 }
    if digits[86] + digits[78] > 10 { digits[96] = (digits[96] + 7) % 10 } else { digits[79] = (digits[79] + 1) % 10 }
    if digits[20]%2 == 1 { digits[20] = (digits[20] + 7) % 10 } else { digits[2] = (digits[2] + 7) % 10 }
    if digits[77] > digits[80] { digits[60] = (digits[60] + 5) % 10 } else { digits[38] = (digits[38] + 8) % 10 }
    if digits[65]%2 == 1 { digits[65] = (digits[65] + 1) % 10 } else { digits[85] = (digits[85] + 3) % 10 }

    // Output result
    for i := 1; i <= 100; i++ {
        fmt.Fprint(writer, digits[i])
    }
}

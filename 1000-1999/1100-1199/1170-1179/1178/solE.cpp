#include <bits/stdc++.h>

using namespace std;

const int LEN = 1000000;
char s[LEN + 1], ans[(LEN >> 1) + 1];

int main() {
    gets(s);
    int len = strlen(s), tmp = len >> 2;
    for (int i = 0; i < tmp; ++i) {
        int k1 = i << 1, k2 = k1 | 1, k3 = len - ((i + 1) << 1), k4 = k3 + 1;
        if (s[k1] == s[k3] || s[k1] == s[k4])
            ans[i] = s[k1];
        else
            ans[i] = s[k2];
    }
    printf("%s", ans);
    if (len % 4) putchar(s[tmp << 1]);
    for (int i = tmp - 1; ~i; --i) putchar(ans[i]);
    return 0;
}
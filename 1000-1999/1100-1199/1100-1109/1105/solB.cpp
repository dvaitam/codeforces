#include <bits/stdc++.h>

using namespace std;
const int N = 2e5 + 5;
int n, k;
int res[26];
char s[N];
int main() {
#ifndef ONLINE_JUDGE
    freopen("input.in", "r", stdin);
#endif
    scanf("%d %d %s", &n, &k, s);
    for (int i = 0; i < n; ) {
        int c = 0;
        char x = s[i];
        while (s[i] == x) {
            ++c;
            ++i;
        }
        res[x - 'a'] += c / k;
    }
    printf("%d\n", *max_element(res, res + 26));
    return 0;
}
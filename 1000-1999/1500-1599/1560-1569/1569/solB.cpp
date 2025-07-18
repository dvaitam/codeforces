#include <bits/stdc++.h>
using namespace std;

int main() {
    const int N = 100;
    int t;
    scanf("%d", &t);
    while (t--) {
        int n;
        char s[N];
        scanf("%d", &n);
        scanf("%s", s);
        int one = 0, two = 0, first = -1, prev = -1, defeats[n];
        for (int i = 0; i < n; i++) {
            if (s[i] == '1') {
                one++;
            } else {
                two++;
                if (first < 0) {
                    first = i;
                    prev = i;
                } else {
                    defeats[i] = prev;
                    defeats[first] = i;
                    prev = i;
                }
            }
        }
        if (two == 1 || two == 2) {
            puts("NO");
            continue;
        }
        puts("YES");
        vector<string> ans(n, string(n, '?'));
        for (int i = 0; i < n; i++) {
            for (int j = i; j < n; j++) {
                if (i == j)
                    ans[i][j] = 'X';
                else if (s[i] == '1' || s[j] == '1')
                    ans[i][j] = ans[j][i] = '=';
                else {
                    ans[i][j] = '+';
                    ans[j][i] = '-';
                    if (defeats[j] == i)
                        swap(ans[i][j], ans[j][i]);
                }
            }
        }
        for (auto &x : ans) {
            printf("%s\n", x.c_str());
        }
    }
    return 0;
}
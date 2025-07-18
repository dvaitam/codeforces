#include<bits/stdc++.h>
using namespace std;

#define MX 507

char s[MX][MX];
int n;

bool check(int i, int j) {
    if(i < 0 || j < 0 || i >= n || j >= n) return false;
    if(s[i][j] != 'X') return false;
    return true;
}

int main() {

    scanf("%d", &n);
    int ans = 0;
    for(int i = 0; i < n; i ++) scanf("%s", s[i]);
    for(int i = 0; i < n; i ++) {
        for(int j = 0; j < n; j ++) {
            if(s[i][j] == 'X') {
                if(check(i - 1, j - 1) && check(i - 1, j + 1) && check(i + 1, j - 1) && check(i + 1, j + 1)) ans ++;
            }
        }
    }
    printf("%d\n", ans);

    return 0;
}
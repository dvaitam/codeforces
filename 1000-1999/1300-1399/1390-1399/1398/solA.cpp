// Problem : A. Bad Triangle
// Contest : Codeforces - Educational Codeforces Round 93 (Rated for Div. 2)
// URL : https://codeforces.ml/contest/1398/problem/0
// Memory Limit : 256 MB
// Time Limit : 1000 ms
// Powered by CP Editor (https://github.com/cpeditor/cpeditor)

#include <algorithm>
#include <cstdio>
#include <iostream>

using namespace std;

template <typename T> inline void read(T& t) {
    t = 0;
    char c = getchar();
    int f = 1;
    while (c < '0' || c > '9') {
        if (c == '-')
            f = -f;
        c = getchar();
    }
    while (c >= '0' && c <= '9') {
        t = (t << 3) + (t << 1) + c - '0';
        c = getchar();
    }
    t *= f;
}
template <typename T, typename... Args> inline void read(T& t, Args&... args) {
    read(t);
    read(args...);
}

const int MAXN = 5e4;

int t, n;

int a[MAXN];

int main() {
    read(t);
    while (t--) {
        read(n);
        for (int i = 0; i < n; i++) {
            read(a[i]);
        }
        bool flg = false;
        for (int i = 2; i < n; i++) {
            if (a[0] + a[1] <= a[i]) {
                printf("%d %d %d\n", 1, 2, i + 1), flg = true;
                break;
            }
        }
        if (!flg)
            printf("-1\n");
    }
    return 0;
}
#pragma GCC optimize(2)
#pragma GCC optimize(3)
#include <bits/stdc++.h>
#define Fast_cin ios::sync_with_stdio(false), cin.tie(0);
#define rep(i, a, b) for(register int i = a; i <= b; i++)
#define per(i, a, b) for(register int i = a; i >= b; i--)
using namespace std;

typedef unsigned long long ull;
typedef pair <int, int> pii;
typedef long long ll;

template <typename _T>
inline void read(_T &f) {
    f = 0; _T fu = 1; char c = getchar();
    while(c < '0' || c > '9') { if(c == '-') fu = -1; c = getchar(); }
    while(c >= '0' && c <= '9') { f = (f << 3) + (f << 1) + (c & 15); c = getchar(); }
    f *= fu;
}

template <typename T>
void print(T x) {
    if(x < 0) putchar('-'), x = -x;
    if(x < 10) putchar(x + 48);
    else print(x / 10), putchar(x % 10 + 48);
}

template <typename T>
void print(T x, char t) {
    print(x); putchar(t);
}

const int N = 1e6 + 5;

int a[N], ans[N], used[N];
ll k, res;
int n, pos, now;

int main() {
    read(n); read(k);
    if(k < 1ll * n * (n + 1) / 2) { print(-1, '\n'); return 0; }
    res = 1ll * n * (n + 1) / 2;
    while(res < k && pos < n) {
        ++pos;
        int det = (int)min(k - res, (ll)n - pos - pos + 1);
        if(det <= 0) { --pos; break; }
        a[pos] = pos + det; res += det;
    }
    for(register int i = 1; i <= n; i++) used[a[i]] = 1;
    now = pos + 1;
    for(register int i = 1; i <= n; i++) {
        if(!used[i]) {
            a[now] = i; ++now;
        }
    }
    now = pos;
    for(register int i = 1; i <= pos; i++) {
        ans[now] = i;
        --now;
    }
    for(register int i = pos + 1; i <= n; i++) ans[i] = i;
    print(res, '\n');
    for(register int i = 1; i <= n; i++) print(a[i], i == n ? '\n' : ' ');
    for(register int i = 1; i <= n; i++) print(ans[i], i == n ? '\n' : ' ');
    return 0;
}
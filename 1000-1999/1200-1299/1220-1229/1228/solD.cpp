#include <bits/stdc++.h>
#include <unordered_map>
using namespace std;
#define int long long
#define lowbit(x)  ((x) & - (x))
using PII = pair<int, int>;
using PIII = tuple<int, int, int>;
const int MOD = 5e18;
const int N = 3e5 + 10;
const long long INF = LLONG_MAX;

inline int gcd(int a, int b)
{
    if (b) while (b ^= a ^= b ^= a %= b);
    return a;
}

inline int mypow(int n, int k, int p = MOD) {
    int r = 1;
    for (; k; k >>= 1, n = n * n % p) {
        if (k & 1) r = r * n % p;
    }
    return r;
}

int inv(int x) { return mypow(x, MOD - 2, MOD); }

int fa[N];
struct cc {
    int id;
    basic_string<int> d;
    friend bool operator<(cc x, cc y) {
        return x.d < y.d;
    }
}s[N];

void solve(){
    int n, m;
    cin >> n >> m;
    for (int i = 1; i <= n; ++i) {
        s[i].id = i;
    }
    for (int i = 1; i <= m; ++i) {
        int x, y;
        cin >> x >> y;
        s[x].d += y;
        s[y].d += x;
    }
    for (int i = 1; i <= n; ++i) {
        sort(s[i].d.begin(), s[i].d.end());
    }
    sort(s + 1, s + n + 1);
    if (!s[1].d.length()) {
        cout << -1 << endl;
        return;
    }
    int cnt = 1;
    fa[s[1].id] = 1;
    for (int i = 2; i <= n; ++i) {
        if (s[i].d == s[i - 1].d) fa[s[i].id] = fa[s[i - 1].id];
        else fa[s[i].id] = ++cnt;
    }
    if (cnt != 3) cout << -1 << endl;
    else {
        for (int i = 1; i <= n; ++i) {
            cout << fa[i] << ' ';
        }
    }
}

signed main(){
    ios::sync_with_stdio(0);
    cin.tie(0);
    cout.tie(0);
    int t;
    t = 1;
    while (t--) {
        solve();
    }
    return 0;
}
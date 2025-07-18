#include <bits/stdc++.h>

#define F first
#define S second
#define x1 privet1
#define x2 privet2
#define y1 privet3
#define y2 privet4
#define left privet6
#define prev ppppp

using namespace std;
typedef long long ll;
typedef long double ld;

const int max_n = 66, log_n = 32, max_m = 111, mod = 1000000007, inf = 1000111222;

int n, a[max_n][max_n], used[7][max_n];
pair<ld, ld> dp[7][max_n];
pair<int, int> v[7][max_n];

pair<ld, ld> get_dp(int x, int b) {
    if (used[x][b]) return dp[x][b];
    if (x == 0) return make_pair(0, 1);
    pair<ld, ld> ans = make_pair(0, 0);
    pair<ld, ld> p = get_dp(x - 1, b);
    for (int i = v[x][b].F; i <= v[x][b].S; ++i) {
        if (v[x - 1][b].F <= i && v[x - 1][b].S >= i) continue;
        // cout << "      " << x << " " << b << " " << i << " " << v[x][b].F << " " << v[x][b].S << endl;
        pair<ld, ld> tmp = get_dp(x - 1, i);
        ans.S += tmp.S * a[b][i] / 100;
        ans.F = max(ans.F, tmp.F + p.F);
    }
    ans.S *= p.S;
    ans.F += ans.S * (1 << (x - 1));
    used[x][b] = 1;
    // cout << x << " " << b << " " << ans.F << " " << ans.S << endl;
    return dp[x][b] = ans;
}

void gett(int a, int l, int r) {
    if (a == -1) return;
    for (int i = l; i <= r; ++i) {
        v[a][i] = make_pair(l, r);
    }
    int mid = (l + r) / 2;
    gett(a - 1, l, mid);
    gett(a - 1, mid + 1, r);
}

int main() {
    //freopen("input.txt", "r", stdin);
    cin >> n;
    for (int i = 0; i < (1 << n); ++i) {
        for (int q = 0; q < (1 << n); ++q) {
            cin >> a[i][q];
        }
    }
    gett(n, 0, (1 << n) - 1);
    /*for (int i =0 ; i <= n; ++i) {
        for (int q = 0; q < (1 << n); ++q) {
            cout << i << " " << q << " " << v[i][q].F << " " << v[i][q].S << endl;
        }
    }*/
    ld ans = 0;
    for (int i =0 ; i < (1 << n); ++i) {
        ans = max(ans, get_dp(n, i).F);
    }
    //for (int i = 0; i < (1 << n); ++i) {
    //    cout << get_dp(1, i).F << " " << get_dp(1, i).S << endl;
    //}
    printf("%.13f", (double)ans);
    return 0;
}
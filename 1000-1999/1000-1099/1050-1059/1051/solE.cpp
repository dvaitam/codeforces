#include <bits/stdc++.h>
#pragma GCC optimize("Ofast,unroll-loops,no-stack-protector,fast-math")
#pragma GCC target("sse,sse2,sse3,ssse3,sse4,popcnt,abm,mmx,avx,tune=native")
#pragma comment(linker, "/STACK:16777216")

#include <bits/stdc++.h>
#include <ext/pb_ds/assoc_container.hpp>
#include <ext/pb_ds/tree_policy.hpp>

using namespace std;
using namespace __gnu_pbds;
using matrix = vector<vector<long long>>;

typedef long long ll;
typedef long double ld;
typedef pair<int, int> pii;
typedef pair<ll, ll> pll;
typedef pair<double, double> pd;
typedef tree<int, null_type, less<int>, rb_tree_tag, tree_order_statistics_node_update> ordered_set;

#define all(x) (x).begin(), (x).end()
//#define int ll

void FAST_IO() {
    ios_base::sync_with_stdio(0);
    cin.tie(0);
    cout.tie(0);
    //cout.setf(ios::fixed);
    //cout.precision(20);
    #ifdef _offline
    freopen("input.txt", "r", stdin);
    freopen("output.txt", "w", stdout);
    #endif // _offline
}

const ll mod = 998244353;
const int MAXN = 2e6 + 100;
int z1[MAXN], z2[MAXN];
ll dp[MAXN], sum[MAXN];

signed main() {
    FAST_IO();
    string a, L, R;
    cin >> a >> L >> R;
    string t = L + '#' + a;
    int l = 0, r = 0;
    for (int i = 1; i < t.length(); ++i) {
        z1[i] = max(0, min(z1[i - l], r - i));
        while (t[z1[i]] == t[i + z1[i]]) {
            z1[i]++;
        }
        if (i + z1[i] > r) {
            l = i;
            r = i + z1[i];
        }
    }
    t = R + '#' + a;
    l = 0, r = 0;
    for (int i = 1; i < t.length(); ++i) {
        z2[i] = max(0, min(z2[i - l], r - i));
        while (t[z2[i]] == t[i + z2[i]]) {
            z2[i]++;
        }
        if (i + z2[i] > r) {
            l = i;
            r = i + z2[i];
        }
    }
    int n = a.length();
    dp[n] = 1;
    sum[n] = 1;
    for (int i = n - 1; i >= 0; --i) {
        if (a[i] == '0') {
            if (L.length() == 1 && L[0] == '0') {
                dp[i] = dp[i + 1];
                sum[i] = (sum[i + 1] + dp[i]);
                if (sum[i] >= mod) {
                    sum[i] -= mod;
                }
                continue;
            }
            dp[i] = 0;
            sum[i] = sum[i + 1];
            continue;
        }
        if (n - i < L.length()) {
            dp[i] = 0;
            sum[i] = sum[i + 1];
            continue;
        }
        int indL = i + L.length() + 1;
        int indR = i + R.length() + 1;
        int lf = i + L.length();
        if (z1[indL] != int(L.length()) && L[z1[indL]] > a[i + z1[indL]]) {
            lf++;
        }
        int rg;
        if (n - i < R.length()) {
            rg = n;
        }
        else if (z2[indR] == int(R.length()) || R[z2[indR]] > a[i + z2[indR]]) {
            rg = i + R.length();
        }
        else {
            rg = i + R.length() - 1;
        }
        if (rg >= lf) {
            ll add = sum[lf];
            add -= sum[rg + 1];
            if (add < 0) {
                add += mod;
            }
            dp[i] = add;
            sum[i] = sum[i + 1] + add;
            if (sum[i] >= mod) {
                sum[i] -= mod;
            }
        }
        else {
            dp[i] = 0;
            sum[i] = sum[i + 1];
        }
    }
    cout << dp[0];
}
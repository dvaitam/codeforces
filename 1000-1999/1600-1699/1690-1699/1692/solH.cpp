// Made by TheQuantiX
#include <bits/stdc++.h>
//#pragma GCC optimize("O3,unroll-loops")
//#pragma GCC target("avx2,sse4.2,bmi,bmi2,lzcnt,popcnt")

using namespace std;

using ll = long long;

constexpr ll MAXN = 1e5;
constexpr ll MOD = 1000000007;
constexpr ll INF = 1000000000000000000LL;

ll tt, n, m, k, x, y, a, b, c, d;

void solve() {
    cin >> n;
    vector<ll> v(n);
    for (int i = 0; i < n; i++) {
        cin >> v[i];
    }
    map<ll, ll> mp;
    map<ll, ll> occ;
    ll mx = 0;
    ll a = -1, l = -1, r = -1;
    for (int i = 0; i < n; i++) {
        set<ll> st;
        if (mp.count(v[i])) {
            mp[v[i]]++;
        }
        else {
            occ[v[i]] = i;
            mp[v[i]] = 1;
        }
        if (mp[v[i]] > mx) {
            a = v[i];
            l = occ[v[i]] + 1;
            r = i + 1;
        }
        mx = max(mx, mp[v[i]]);
        for (auto j : mp) {
            if (j.first != v[i]) {
                mp[j.first]--;
                if (mp[j.first] == 0) {
                    st.insert(j.first);
                }
            }
        }
        for (int j : st) {
            mp.erase(j);
            occ.erase(j);
        }
    }
    cout << a << ' ' << l << ' ' << r << '\n';
}

int main() {
    ios_base::sync_with_stdio(0);
    cin.tie(0);
    cout.tie(0);
    cout << fixed;
    cout << setprecision(10);
    tt = 1;
    cin >> tt;
    while (tt --> 0) {
        solve();
    }
}
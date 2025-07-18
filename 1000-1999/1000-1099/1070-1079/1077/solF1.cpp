#include <bits/stdc++.h>
using namespace std;
typedef long long ll;

const ll maxn = 2e2 + 1;

ll d[maxn][maxn];

int main() {
    ios_base::sync_with_stdio(0);
    ll n, k, x;
    cin >> n >> k >> x;
    for (ll i = 0; i <= n; i++) {
        for (ll j = 0; j <= x; j++) {
            d[i][j] = INT_MIN;
        }
    }
    ll a[n + 1];
    for (ll i = 0; i++ < n;) cin >> a[i];
    
    d[0][0] = 0;
    for (ll i = 0; i < n; i++) {
        for (ll j = 0; j < x; j++) {
            if(d[i][j] != INT_MIN) {
                for (ll z = i + 1; z <= i + k && z <= n; z++) {
                    d[z][j + 1] = max(d[z][j + 1], d[i][j] + a[z]);
                }
            }
        }
    }
    
    /*for (ll i = 0; i <= n; i++) {
        for (ll j = 0; j <= x; j++) {
            if(d[i][j] == INT_MIN) cout << -1 << ' ';
            else cout << d[i][j] << ' ';
        }
        cout << endl;
    }*/
    
    ll answ = INT_MIN;
    for (ll i = n - k + 1; i <= n; i++) {
        answ = max(answ, d[i][x]);
    }
    if(answ != INT_MIN) cout << answ;
    else cout << -1;
}
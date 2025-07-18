#include <bits/stdc++.h>
#include <ext/pb_ds/assoc_container.hpp>
#include <ext/pb_ds/tree_policy.hpp>
#define ll long long
#define endl "\n"
 
using namespace std;
using namespace __gnu_pbds;
mt19937 rng(chrono::steady_clock::now().time_since_epoch().count());
template<class T> using ordered_set = tree<T, null_type, less_equal<T>, rb_tree_tag, tree_order_statistics_node_update>;
const ll N = 2e5 + 5;
    ll a[3][N];
void solve()
{
    ll n;
    cin >> n;
    for (ll i = 1; i <= n; i++) cin >> a[0][i];
    for (ll i = 1; i <= n; i++) cin >> a[1][i];
    for (ll i = 1; i <= n; i++) cin >> a[2][i];
    ll p[] = {0, 1, 2}, need = 0;
    for (ll i = 1; i <= n; i++) need += a[0][i];
    need = (need + 2) / 3;
    do
    {
        pair<ll, ll> seg[3];
        ll ptr = 1;
        bool ok = true;
        for (ll x : p)
        {
            seg[x].first = ptr;
            ll cur = 0;
            while (ptr <= n and cur < need) cur += a[x][ptr++];
            if (cur < need) ok = false;
            seg[x].second = ptr - 1;
        }
        if (ok)
        {
            for (ll i = 0; i < 3; i++) cout << seg[i].first << " " << seg[i].second << " ";
            cout << endl;
            return;
        }
    } while (next_permutation(p, p + 3));
    cout << "-1\n";
}

int main()
{
    ios_base::sync_with_stdio(0);
    cin.tie(0);
    ll t = 1;
    // precomp();
    cin >> t;
    for (ll cs = 1; cs <= t; cs++)
        solve();
    // cerr << "\nTime elapsed: " << clock() * 1000.0 / CLOCKS_PER_SEC << " ms\n";
}
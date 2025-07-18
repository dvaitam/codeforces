#include<bits/stdc++.h>
#include <ext/pb_ds/assoc_container.hpp>
using namespace std;
typedef long long int ll;
typedef long double ld;
#define sz(x) ((int)(x).size())
#define pb push_back
#define mp make_pair
#define int long long
#define pii pair<int, int>
#define pip pair<ll, pii>
//#define x first
//#define y second

signed main() {
    ios::sync_with_stdio(0);
    cin.tie(0);
    cout.tie(0);
    int n;
    cin >> n;
    int na = 0, nb = 0;
    int ans = 0;
    for (int i = 0; i < n; ++i) {
        int a,b;
        cin >> a >> b;
        int x = max(na, nb);
        int y = min(a, b);
        ans += max(0LL, y - x + 1);
        na = a, nb = b;
        if (na < nb)
            ++na;
        else
            ++nb;
    }
    cout << ans;
    return 0;
}
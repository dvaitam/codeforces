// include
#include<bits/stdc++.h>

using namespace std;

// random
mt19937 rng(chrono::steady_clock::now().time_since_epoch().count());

// templates
template<class X, class Y>
    bool maximize(X &x, const Y &y) {
        if (y > x) {
            x = y;
            return (true);
        } else return (false);
    }

template<class X, class Y>
    bool minimize(X &x, const Y &y) {
        if (y < x) {
            x = y;
            return (true);
        } else return (false);
    }

// define
#define fi    first
#define se    second

#define pub   push_back
#define pob   pop_back
#define puf   push_front
#define pof   pop_front
#define eb    emplace_back
#define upb   upper_bound
#define lwb   lower_bound

#define left  VAN
#define right TAN

#define all(a) (a).begin(),(a).end()
#define rall(a) (a).begin(),(a).end()
#define sort_and_unique(a) sort(all(a));(a).resize(unique(all(a))-a.begin())

// another define
using ll  = long long;
using ld  = long double;
using pii = pair<int, int>;
using pil = pair<int, ll>;
using pli = pair<ll, int>;
using pll = pair<ll, ll>;

// limit
const int oo = 2e9;

void solve() {
    ll X;
    cin >> X;
    ll k;
    for (ll i = 60; i >= 0; i--)
        if (X >> i & (1ll)) {
            k = i;
            break;
        }
    ll p = 200 - k + 1  ;
    vector<int> ans;
    for (ll i = 200 - k + 1; i <= 200; i++) {
        ans.pub(i);

        if (X >> (200ll - 1ll * i) & (1ll)) {
            --p;
            ans.pub(p);
        }
    }
    cout << (int)ans.size() << '\n';
    for (int x : ans) cout << x << ' ';
    cout << '\n';
}

int main() {
    ios_base::sync_with_stdio(false);
    cin.tie(0); cout.tie(0);
    int t = 1;
    cin >> t;

    while (t--) solve();
    return 0;
}
#include <bits/stdc++.h>

#define nl "\n"
#define no "NO"
#define yes "YES"
#define fi first
#define se second
#define vec vector
#define task "main"
#define _mp make_pair
#define ii pair<int, int>
#define sz(x) (int)x.size()
#define all(x) x.begin(), x.end()
#define evoid(val) return void(std::cout << val)
#define FOR(i, a, b) for(int i = (a); i <= (b); ++i)
#define FOD(i, b, a) for(int i = (b); i >= (a); --i)
#define unq(x) sort(all(x)); x.resize(unique(all(x)) - x.begin())

using namespace std;

template<typename U, typename V> bool maxi(U &a, V b) {
    if (a < b) { a = b; return 1; } return 0;
}
template<typename U, typename V> bool mini(U &a, V b) {
    if (a > b) { a = b; return 1; } return 0;
}

const int N = (int)2e5 + 9;
const int mod = (int)1e9 + 7;

void prepare(); void main_code();

int main() {
    ios::sync_with_stdio(0); cin.tie(0); cout.tie(0);
    if (fopen(task".inp", "r")) {
        freopen(task".inp", "r", stdin);
        freopen(task".out", "w", stdout);
    }
    const bool MULTITEST = 1; prepare();
    int num_test = 1; if (MULTITEST) cin >> num_test;
    while (num_test--) { main_code(); cout << "\n"; }
}

void prepare() {};

int a[N];

void main_code() {
    int n; cin >> n;
    FOR(i, 1, n) cin >> a[i];
    if (n % 2 == 0) {
        cout << 2 << nl;
        cout << 1 << ' ' << n << nl;
        cout << 1 << ' ' << n << nl;
        return ;
    }
    cout << 4 << nl;
    cout << 1 << ' ' << n - 1 << nl;
    cout << 1 << ' ' << n - 1 << nl;
    cout << 2 << ' ' << n << nl;
    cout << 2 << ' ' << n << nl;
}


/*     Let the river flows naturally     */
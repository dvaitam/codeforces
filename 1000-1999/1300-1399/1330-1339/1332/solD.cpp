#include <bits/stdc++.h>



using namespace std;



typedef long long ll;

typedef long double ld;

typedef pair<int, int> pii;

typedef pair<ll, ll> pll;

typedef vector<int> vi;



#define fi first

#define se second

#define pp push_back

#define all(x) (x).begin(), (x).end()

#define Ones(n) __builtin_popcount(n)

#define endl '\n'

#define fill(arrr,xx) memset(arrr,xx,sizeof arrr)

#define rep(aa, bb, cc) for(int aa = bb; aa < cc;aa++)

#define PI acos(-1)

//#define int long long



void Gamal() {

    ios_base::sync_with_stdio(false);

    cin.tie(nullptr);

    cout.tie(nullptr);

#ifdef Clion

    freopen("input.txt", "r", stdin), freopen("output.txt", "w", stdout);

#endif

}



int dx[] = {+0, +0, -1, +1, +1, +1, -1, -1};

int dy[] = {-1, +1, +0, +0, +1, -1, +1, -1};



const double EPS = 1e-9;

const ll N = 2e5 + 5, INF = INT_MAX, MOD = 1e9 + 7, OO = 0X3F3F3F3F3F3F3F3F, LOG = 20;



void solve() {

    int k;cin >> k;

    cout << 3 << ' ' << 3 << endl;

    int add = 1 << 17;

    cout << (k + add) << ' ' << (k + add) << ' ' << k << endl;

    cout << add << ' ' << add << ' ' << (k + add) << endl;

    cout << 0<< ' ' << 0 << ' ' << k << endl;

}





signed main() {

    Gamal();

    int t = 1;

//    cin >> t;

    while (t--) {

        solve();

    }

}
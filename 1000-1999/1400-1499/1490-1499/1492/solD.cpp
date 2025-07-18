#include <bits/stdc++.h>



using namespace std;



#define mst(a, b) memset(a, b, sizeof(a))

#define endl '\n'

#define rep(i, b, e) for (int i = b; i < (int)e; i++)

#define repn(i, e) for (int i = 0; i < (int)e; i++)

#define YES cout << "YES" << endl

#define NO cout << "NO" << endl

#define all(_) _.begin(), _.end()

#define sz(_) (int)_.size()

#define fst first

#define snd second

#define Odd(_x) ((_x) & 1)

typedef long long ll;

typedef pair<int, int> pii;

const int N = 100;

const int maxn = 100000 + 10;

const int inf = 1e9 + 10;

const ll mod = 1e9 + 7;



void solve() {

    int a, b, k;

    cin >> a >> b >> k;

    if (a == 0 || b == 1) {

        if (k) NO;

        else {

            YES;

            cout << string(b, '1') << string(a, '0') << endl;

            cout << string(b, '1') << string(a, '0') << endl;

        }

        return;

    }

    if (k > a + b - 2) {

        NO;

        return;

    }

    string x(a + b, '1'), y(a + b, '1');

    x[0] = '0';

    a--;

    for (int i = 1; i < k && a > 0; i++) {

        x[i] = '0';

        y[i] = '0';

        a--;

    }

    y[k] = '0';

    while (a) {

        x[++k] = '0';

        y[k] = '0';

        a--;

    }

    YES;

    std::reverse(x.begin(), x.end());

    std::reverse(y.begin(), y.end());

    cout << x << endl;

    cout << y << endl;

}

//#define MULTI_INPUT



int main() {

#ifndef ONLINE_JUDGE

    freopen(R"(D:\source files\source file2\input.txt)", "r", stdin);

    freopen(R"(D:\source files\source file2\output.txt)", "w", stdout);

#endif

    std::ios::sync_with_stdio(false);

    cin.tie(0);

#ifdef MULTI_INPUT

    int T;

    cin >> T;

    while (T--) {

        solve();

    }

#else

    solve();

#endif

    return 0;

}
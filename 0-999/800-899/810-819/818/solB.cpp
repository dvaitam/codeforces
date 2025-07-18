#include <bits/stdc++.h>
#define pb push_back
#define mp make_pair
#define pp pair<int, int>
#define ppp pair<pp, int>
#define fi first
#define se second
#define N 111
#define M 777
#define mod 1000000007
#define inf 1000000001
#define esp 1e-9
typedef long long ll;
using namespace std;
int n, m;
int a[N];
int d[N];
bool f[N];
int mask[N];

int main() {
	//freopen("in.in", "r", stdin);
	//freopen("ou.ou", "w", stdsout);
    ios::sync_with_stdio(false);
    cin >> n >> m;
    bool ok = true;
    for (int i = 1; i <= m; i++) cin >> a[i];
    for (int i = 1; i < m; i++) {
        int dd = a[i + 1] - a[i];
        if (dd <= 0) dd += n;
        if (mask[a[i]] == 0) mask[a[i]] = dd;
        else
            if (mask[a[i]] != dd) ok = false;
    }

    if (!ok) {
        cout << -1;
    }
    else {
        for (int i = 1; i <= n; i++)
        if (mask[i] > 0) {
            if (f[mask[i]]) {
                cout << -1;
                return 0;
            }
            else
            f[mask[i]] = true;
        }

        vector <int> tmp;
        for (int i = 1; i <= n; i++)
        if (f[i] == false) {
            tmp.pb(i);
        }
        for (int i = 1; i <= n; i++)
        if (mask[i] == 0) {
            mask[i] = tmp[tmp.size() - 1];
            tmp.pop_back();
        }

        for (int i = 1; i <= n; i++)
            cout << mask[i] << " ";
    }
    /**/return 0;
}
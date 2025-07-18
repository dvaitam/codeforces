/// endless ?



#pragma GCC optimize("O3")

#pragma GCC optimize("unroll-loops")

#include <bits/stdc++.h>

using namespace std;



typedef long double ld;

#define ll                       long long

#define F                        first

#define S                        second

#define pii                      pair<int, int>

#define all(x)                   x.begin(), x.end()

#define vi                       vector<int>

#define vii                      vector<pii>

#define pb                       push_back

#define pf                       push_front

#define wall                     cout <<'\n'<< "-------------------------------------" <<'\n';

#define fast                     ios_base::sync_with_stdio(false);cin.tie(0);cout.tie(0);



#define int ll

const ll MAXN = 3e5 + 43;

const ll MOD  = 1e9  + 7; ///998244353;

const ll INF  = 1e18 + 19763;

const ll LG   = 19;

ll pw(ll a, ll b){return b == 0 ? 1LL : (pw(a * a%MOD , b / 2)%MOD * (b % 2 == 0 ? 1LL : a))%MOD;}

int a[MAXN][10], n, m, bt[MAXN], mrk[MAXN];

pii ans = {-1, -1};



bool check(int x)

{

    ans = {-1, -1};

    fill(mrk, mrk + 600, 0);

    for (int i = 0; i < n; i++)

    {

        int num = 0;

        for (int j = 0; j < m; j++)

        {

            if (a[i][j] >= x)

                num += (1 << j);

        }

        //cout << num << '\n';

        mrk[num] = i + 1;

    }



    for (int i = 0; i < 300; i++)

    {

        ///cout << i << ' ' << mrk[i] << '\n';

        for (int j = 0; j < 300; j++)

        {

            int ord = i|j;

            if (ord == ((1 << m) - 1) && mrk[i] != 0 && mrk[j] != 0){ans = {mrk[i], mrk[j]}; return 1;}

        }

    }

    return 0;

}



void solve()

{

    fast

    cin >> n >> m;

    for (int i = 0; i < n; i++)

    {

        for (int j = 0; j < m; j++)

        {

            cin >> a[i][j];

        }

    }



    int l = 0, r = 1e9 + 3;

    while (r - l > 1)

    {

        int mid = (l + r)/2;

        ///cout << mid << ' ' << check(mid) << ' ' << ans.F << ' ' << ans.S << '\n';

        if (check(mid))l = mid;

        else r = mid;

    }

    bool bl = check(l);

    int mid = 2;

    //cout << mid << ' ' << check(mid) << ' ' << ans.F << ' ' << ans.S << '\n';

    cout  << ans.F << ' ' << ans.S << '\n';

}



int32_t main ()

{

    fast

    int t = 1;// cin >> t;

    while (t --)

    {

        solve();

    }

}

/// Thanks GOD :)
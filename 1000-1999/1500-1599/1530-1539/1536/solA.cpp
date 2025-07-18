#include<bits/stdc++.h>

using namespace std;

#define int long long

int mod = 1000000007;

#define PII pair<long long,long long>

#define x first

#define y second

#define pi 3.14159265359

int qpow(int a, int b)

{

    int res = 1;

    while (b != 0)

    {

        if (b & 1)res = res * a % mod;

        a = a * a % mod;

        b >>= 1;

    }

    return res;

}

int inv(int a)

{

    return qpow(a, mod - 2);

}

int ai, bi, xi, yi;

int exgcd(int ai, int bi, int& xi, int& yi)

{

    if (bi == 0)

    {

        xi = 1, yi = 0;

        return ai;

    }

    int d = exgcd(bi, ai % bi, yi, xi);

    yi -= ai / bi * xi;

    return d;

}

int gcd(int a, int b)

{

    return b == 0 ? a : gcd(b, a % b);

}

/*----------------------------------------------------------------------*/







void solve()

{

    int n;

    cin >> n;

    vector<int> a(n);

    int minv = INT_MAX;

    int maxv = 0;

    set<int> s;

    for (int i = 0; i < n; i++)

    {

        cin >> a[i];

        minv = min(minv, a[i]);

        maxv = max(maxv, a[i]);

        s.insert(a[i]);

    }

    if (minv < 0)

    {

        cout << "NO\n";

        return;

    }

    cout << "YES\n";

    sort(a.begin(), a.end());

    for (int i = 0; i < a.size(); i++)

    {

        for (int j = i + 1; j < a.size(); j++)

        {

            int tmp = abs(a[j] - a[i]);

            if (!s.count(tmp))

            {

                a.emplace_back(tmp);

                s.insert(tmp);

                sort(a.begin(), a.end());

                i = -1;

                break;

            }

        }

    }

    cout << s.size() << "\n";

    for (int it : s)cout << it << " ";

    cout << "\n";

}



signed main()

{

    ios::sync_with_stdio(0);

    cin.tie(0);

    cout.tie(0);

    int tt = 1;

    cin >> tt;

    while (tt--)solve();

    return 0;

}
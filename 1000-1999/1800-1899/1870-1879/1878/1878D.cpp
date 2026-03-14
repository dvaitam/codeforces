#include <bits/stdc++.h>
using namespace std;

typedef long long ll;
typedef vector<ll> vec;

const ll MOD = 1000000007;
const int N = 2000005;

ll fact[N];
ll invFact[N];

ll power(ll a, ll b)
{
    ll res = 1;
    while (b)
    {
        if (b & 1)
            res = (res * a) % MOD;
        a = (a * a) % MOD;
        b >>= 1;
    }
    return res;
}

void precompute()
{
    fact[0] = 1;

    for (int i = 1; i < N; i++)
        fact[i] = (fact[i - 1] * i) % MOD;

    invFact[N - 1] = power(fact[N - 1], MOD - 2);

    for (int i = N - 2; i >= 0; i--)
        invFact[i] = (invFact[i + 1] * (i + 1)) % MOD;
}

ll nCr(ll n, ll r)
{
    if (r > n || r < 0)
        return 0;
    return (((fact[n] * invFact[r]) % MOD) * invFact[n - r]) % MOD;
}

int main()
{
    ios::sync_with_stdio(false);
    cin.tie(NULL);

    precompute();

    ll t;
    cin >> t;

    while (t--)
    {
        ll n, k;
        cin >> n >> k;

        string s;
        cin >> s;

        vector<char> v(n + 1);
        for (int i = 1; i <= n; i++)
            v[i] = s[i - 1];

        vec l(k+1), r(k+1);
        for(int i=1;i<=k;i++) cin>>l[i];
        for(int i=1;i<=k;i++) cin>>r[i];

        ll q;
        cin >> q;

        vec x(q + 1);
        for (int i = 1; i <= q; i++)
            cin >> x[i];

        vec pre(n + 2, 0);
        for (auto val : x)
        {
            ll idx = lower_bound(r.begin(), r.end(), val) - r.begin();

            ll fi = min(l[idx] + r[idx] - val, val);
            ll se = max(l[idx] + r[idx] - val, val);
            pre[fi]++;
            pre[se + 1]--;
        }
        for (int i = 1; i < pre.size(); i++)
            pre[i] += pre[i - 1];

        ll i=1;
        while(i<=k)
        {
            ll first = l[i];
            ll last = r[i];
            for (int j = first; j<=(first+(last-first)/2); j++)
            {
                if ((pre[j] % 2) == 0)
                    continue;
                else
                {
                    ll diff = j - first;
                    swap(v[j], v[last - diff]);
                }
            }
            i++;
        }
        for (int i = 1; i <= n; i++)
            cout << v[i];
        cout << endl;
    }
}

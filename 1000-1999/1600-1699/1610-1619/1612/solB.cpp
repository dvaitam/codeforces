#include <bits/stdc++.h>
using namespace std;
#define int long long
#define nl "\n"
#define fastio                        \
    ios_base::sync_with_stdio(false); \
    cin.tie(NULL);                    \
    cout.tie(NULL);
#define all(x) x.begin(), x.end()
int gcd(int a, int b)
{
    if (b == 0)
        return a;
    return gcd(b, a % b);
}

int lcm(int a, int b)
{
    return (a / gcd(a, b)) * b;
}
int CEIL(int a, int b)
{
    return (a + b - 1) / b;
}
int power_mod(int x, int y, int mod)
{

    int res = 1;
    x = x % mod;
    while (y > 0)
    {
        if (y & 1)
            res = (res * x) % mod;
        y = y >> 1;
        x = (x * x) % mod;
    }
    return res;
}

void solve()
{
    int a, b, n;
    cin >> n >> a >> b;
    if (a <=n / 2 + 1 && b >=n / 2)
    {
        if(a == n / 2 + 1 && b == n / 2){
            cout << a << " ";
            for (int i = n / 2 + 2; i <= n; i++)
            {
                    cout << i << " ";
            }
            for (int i = 1; i < n / 2; i++)
            {
               cout << i << " ";
            }
            cout << b << nl;
          }
        else if (a == n / 2 + 1 || b == n / 2)
        {
            if (b > a )
                cout << -1<<nl;
        }
        else
        {
            cout << a << " ";
            for (int i = n / 2 + 1; i <= n; i++)
            {
                if (i == b || i == a)
                    continue;
                else
                    cout << i << " ";
            }
            for (int i = 1; i <= n / 2; i++)
            {
                if (i == a || i == b)
                    continue;
                else
                    cout << i << " ";
            }
            cout << b << nl;
        }
    }
    else
    cout<<-1<<nl;

}

signed main()
{
    fastio
        // freopen("input.txt", "r", stdin);
    int t;
    cin >> t;
    while (t--)
        solve();
    return 0;
}
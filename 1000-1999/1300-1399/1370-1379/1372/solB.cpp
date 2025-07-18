#include <bits/stdc++.h>

using namespace std;

#define ll long long int

#define mod 1000000007



void solve()

{

    ll n, temp;

    cin >> n;

    temp = n;

    map<ll, ll> m;



    for (ll i = 2; i * i <= n; i++)

    {

        while (n % i == 0)

        {

            m[i]++;

            n /= i;

        }

    }

    if (n >= 2)

        m[n]++;

    ll mn;

    for (auto it : m)

    {

        mn = it.first;

        break;

    }

    cout << temp - (temp / mn) << " " << temp / mn << endl;

}

int main()

{

    ios_base::sync_with_stdio(false);

    cin.tie(NULL);

    cout.tie(NULL);

    ll test;

    cin >> test;

    while (test--)

    {

        solve();

    }

    return 0;

}
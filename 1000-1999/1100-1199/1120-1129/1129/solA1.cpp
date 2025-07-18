#include <bits/stdc++.h>

#define ll long long

using namespace std;

int main()
{
    cin.sync_with_stdio(0);
    ll n, m;
    cin >> n >> m;
    ll* noc = new ll[n];
    ll* bd = new ll[n];
    for(ll i = 0; i < n; i++)
    {
        noc[i] = 0;
        bd[i] = n;
    }
    ll mc = 0;
    for(ll i = 0; i < m; i++)
    {
        ll a, b;
        cin >> a >> b;
        a--;
        b--;
        noc[a]++;
        ll dst = b - a;
        if(dst < 0)
        {
            dst += n;
        }
        bd[a] = min(bd[a], dst);
        mc = max(noc[a], mc);
    }
    ll vas = (mc - 1) * n;
    ll lv = 0;
    for(ll i = 0; i < n; i++)
    {
        if(noc[i] == mc)
        {
            lv = max(lv, i + bd[i]);
        }
        if(noc[i] == mc - 1 && mc > 1)
        {
            lv = max(lv, bd[i] + i - n);
        }
    }
    cout << (vas + lv);
    for(ll i = 1; i < n; i++)
    {
        cout << " ";
        lv--;
        if(noc[i - 1] == mc)
        {
            lv = max(lv, n - 1 + bd[i - 1]);
        }
        if(noc[i - 1] == mc - 1 && mc > 1)
        {
            lv = max(lv, bd[i - 1] - 1);
        }
        cout << (vas + lv);
    }
    cout << endl;
    cin >> n;
}
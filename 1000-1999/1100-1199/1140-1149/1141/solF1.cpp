#include <bits/stdc++.h>

#define ll long long

using namespace std;

const ll mod = 1e9+7;

bool cmp(pair<ll,ll>a,pair<ll,ll>b)

{

    return a.second < b.second;

}

int main()

{

    ios::sync_with_stdio(false);

    cin.tie(NULL);

    cout.tie(NULL);

    ll n;

    cin >> n;

    vector<ll>a(n+1),pref(n+1,0);

    map<ll,vector<pair<ll,ll>>>M;

    vector<pair<ll,ll>>R;

    for(ll i=1;i<=n;i++)

    {

        cin >> a[i];

        if(i) pref[i]+=pref[i-1];

        pref[i]+=a[i];

    }

    for(ll i=1;i<=n;i++)

    {

        for(ll j=i;j<=n;j++)

        {

            M[pref[j]-pref[i-1]].push_back({i,j});

        }

    }

    for(auto &x:M)

    {

        sort(x.second.begin(),x.second.end(),cmp);

        ll r=0;

        vector<pair<ll,ll>>tmp;

        for(ll i=0;i<x.second.size();i++)

        {

           // cout << x.second[i].first << ' ' << x.second[i].second << endl;

            if(x.second[i].first>r)

            {

             //   cout << "usao\n";

                tmp.push_back(x.second[i]);

                r=x.second[i].second;

            }

        }

        if(tmp.size()>R.size()) R=tmp;

    }

    cout << R.size() << endl;

    for(auto x:R) cout << x.first << ' ' << x.second << endl;

    return 0;

}
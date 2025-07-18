#pragma GCC optimize("Ofast,unroll-loops,no-stack-protector")

#include <bits/stdc++.h>

#include <iostream>

#include<set>

using namespace std;

#define ll long long

#define endl '\n'

#define endl2(j,x) (j==x-1?"\n":"")

#define pb push_back

#define ff first

#define ss second









#define si(x) int(x.size())

#define sum_all(a)     ( accumulate ((a).begin(), (a).end(), 0ll))

#define mine(a)    (*min_element((a).begin(), (a).end()))

#define maxe(a)    (*max_element((a).begin(), (a).end()))

#define mini(a)    ( min_element((a).begin(), (a).end()) - (a).begin())

#define maxi(a)    ( max_element((a).begin(), (a).end()) - (a).begin())





int t ;

int i0;



void solve()

{

    int n;

    cin >> n;



    vector<int> a(n);



    for (int i = 0; i < n; i++)

    {

        cin >> a[i];

    }



    vector<int> freq(n);

    int cur = 0;

    for (int i = 0; i < n; i++)

    {

        a[i]--;

        freq[a[i]]++;

        if (freq[a[i]] > freq[cur])

        {

            cur = a[i];

        }

    }



    vector<vector<int>> pos(n);

    for (int i = 0; i < n; i++)

    {

        pos[a[i]].push_back(i);

    }

    auto comp = [&](const vector<int>& L, const vector<int>& R) -> bool

    {

        return si(R) < si(L);

    };

    sort(pos.begin(), pos.end(), comp);



    vector<int> res(n);

    for (int i = 0; i < freq[cur]; i++)

    {

        vector<pair<int,int>> cyc;

        for (auto& v : pos)

        {

            if (si(v) > 0)

            {

                cyc.push_back({a[v.back()],v.back()});

                v.pop_back();

            }

            else

            {

                break;

            }

        }

        sort(cyc.begin(), cyc.end());

        for (int i = 0; i < si(cyc); i++)

        {

            res[cyc[i].ss] = cyc[(i + 1) % si(cyc)].ss;

        }

    }

    // vector<int> ans(n);

    // for(int i=0;i<n;i++) ans[i] = a[res[i]]+1;





    for (auto r : res)

        cout << a[r] + 1 << " ";

    cout << endl;

}







signed main()

{

    ios_base::sync_with_stdio(false);

    cin.tie(NULL);

    #ifndef ONLINE_JUDGE

        freopen("input.txt", "r", stdin);

        freopen("output.txt", "w", stdout);

    #endif

    cin >> t;

    for(i0=1; i0<=t;i0++) 

        solve();



}
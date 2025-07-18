#include<bits/stdc++.h>

using namespace std;

using ll = long long;

#define endl '\n' 



void solve()

{

    ll n, m; cin >> n >> m;

    for(int i=1;i<=n;i++){

        for(int j=1;j<=m;j++){

            cout<<((i%4<=1) == (j%4<=1))<<" ";

        }

        cout<<'\n';

    }

    

}

int main()

{

    ios_base::sync_with_stdio(false);

    cin.tie(NULL);

    int t;

    cin>>t;

    while(t--)

    {

        solve();

    }

}
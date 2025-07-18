#include <bits/stdc++.h>

using namespace std;

#define ll long long

void solve()

{

    int n;

    cin>>n;

    int i=3;

    for(int j=0;j<n;j++)

    {

        cout<<i<<" ";

        i+=2;

    }cout<<'\n';

}

int main()

{

    ios::sync_with_stdio;

    cin.tie(0);

    int t;

    cin >> t;

    while (t--)

    {

        solve();

    }

}
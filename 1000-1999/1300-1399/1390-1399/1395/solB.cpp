#include<bits/stdc++.h>

#define ll long long

#define endl "\n"

#define fast ios::sync_with_stdio(0),cin.tie(0),cout.tie(0);

using namespace std;

ll q,t,p,n,y,m,x;

vector<ll>v;

int main()

{

    fast

     cin>>n>>m>>x>>y;

      v.push_back(y);

       for(ll i=1;i<=m;i++){

        if(i!=y){

           v.push_back(i);

        }

       }

       for(ll i=x;i<=n;i++){

        for(ll j=1;j<=m;j++){

            cout<<i<<" "<<v[j-1]<<endl;

        }

        reverse(v.begin(),v.end());

       }







       for(ll i=1;i<x;i++){

        for(ll j=1;j<=m;j++){

            cout<<i<<" "<<v[j-1]<<endl;

        }

        reverse(v.begin(),v.end());

       }



    return 0;

}
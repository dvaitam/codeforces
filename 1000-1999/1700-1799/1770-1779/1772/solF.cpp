#include<bits/stdc++.h>

using namespace std;

#define ll long long

int main(){

    ios::sync_with_stdio(false);

    cin.tie(nullptr);

    ll n,m,k;cin>>n>>m>>k;

    char a[k+2][n][m];

    for(ll i=1;i<=k+1;i++){

        for(ll j=0;j<n;j++){

            for(ll l = 0;l<m;l++){

                cin>>a[i][j][l];

            }

        }

    }

    vector<pair<ll,ll>> v;

    set<ll> s;

    for(ll i=1;i<=k+1;i++){

        ll x = 0;

        for(ll j=1;j<n-1;j++){

            for(ll l=1;l<m-1;l++){

                if((a[i][j-1][l]==a[i][j+1][l]) && (a[i][j][l+1]==a[i][j][l-1]) && (a[i][j-1][l]==a[i][j][l+1]) && (a[i][j][l]!=a[i][j+1][l])){

                    x++;

                }

            }

        }

        v.push_back({x,i});

        s.insert(x);

    }

    sort(v.begin(),v.end());

    cout<<v.back().second<<"\n";

    vector<string> ans;

    ll moves = k;

    ll current = v.back().second;

    ll i = v.size();

    i-=2;

    for(i;i>=0;i--){

        ll r = v[i].second;

        for(ll j=1;j<n-1;j++){

            for(ll l=1;l<m-1;l++){

                if(a[r][j][l]!=a[current][j][l]){

                    string temp = "1 ";

                    temp+=to_string(j+1);

                    temp.push_back(' ');

                    temp+=to_string(l+1);

                    ans.push_back(temp);

                    moves++;

                }

            }

        }

        string t = "2 ";

        t+=to_string(r);

        ans.push_back(t);

        current = r;

    }

    cout<<moves<<"\n";

    for(auto &x:ans){

        cout<<x<<"\n";

    }

    return 0;

}
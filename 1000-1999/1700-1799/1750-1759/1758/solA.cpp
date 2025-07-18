#include<bits/stdc++.h>

using namespace std;



void solve(){

    string st,t;

    cin>>st;

    t=st;

    reverse(st.begin(),st.end());

    cout << st<<t<<"\n";

}

int main(){

    ios_base::sync_with_stdio(0);

    cin.tie(0);

    int t;

    cin>>t;

    while (t--)

    {

        solve();

    }

    return 0;

}
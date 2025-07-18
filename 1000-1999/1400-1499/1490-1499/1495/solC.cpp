#include<bits/stdc++.h>

using namespace std;

char c[505][505];

int n,m;

void solve(){

    cin >> n >> m;

    for(int i=1;i<=n;i++){

        for(int j=1;j<=m;j++) cin >> c[i][j];

    }

    for(int i=(m%3==0)?2:m%3;i<=m;i+=3){

        for(int j=1;j<=n;j++) c[j][i]='X';

        if(i>=3){

            bool ok=false;

            for(int j=1;j<=n;j++){

                if(c[j][i-1]=='X'){c[j][i-2]='X';ok=true;break;}

                if(c[j][i-2]=='X'){c[j][i-1]='X';ok=true;break;}

            }

            if(!ok){c[1][i-1]=c[1][i-2]='X';}

        }

    }

    for(int i=1;i<=n;i++){

        for(int j=1;j<=m;j++) cout << c[i][j];

        cout << '\n';

    }

}

signed main(){

    ios_base::sync_with_stdio(false);

    cin.tie(NULL);cout.tie(NULL);

    int t;cin >> t;

    while(t--) solve();

}
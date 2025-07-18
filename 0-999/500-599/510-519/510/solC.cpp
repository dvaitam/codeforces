#include<bits/stdc++.h>

#define pb push_back

using namespace std;



typedef long long LL;



int color[26];

int color1[26];

int in[26];

vector<vector<int> >g(26);



vector<int> ans;



bool dfs1(int v)

{

    color1[v] = 1;

    for(int i=0;i<g[v].size();i++){

        int to = g[v][i];

        if(color1[to]==2)

            continue;

        if(color1[to]==1)

            return true;

        if( dfs1(to) )

            return true;

    }

    color1[v] = 2;

    return false;

}

void dfs(int v)

{

    color[v] = 1;

    for(int i=0;i<g[v].size();i++){

        int to = g[v][i];

        if( color[to]==2 )

            continue;

        if(color[to]==1){

            return ;

        }

        dfs( to );

    }

    color[v] = 2;

    ans.pb(v);

}

int main()

{

    int t;

    cin>>t;

    vector<string>v(t);

    for(int i=0;i<t;i++){

        cin>>v[i];

        v[i]+="#";

    }

    for(int i=1;i<t;i++){

        

        int k = 0;

        while( v[i][k] == v[i-1][k] )

            k++;

        if( (v[i-1][k] )=='#' )

            continue;

        if(( v[i][k] )=='#'){

            cout<<"Impossible"<<endl;

            return 0;

        }

        g[ v[i-1][k] -'a' ].pb( v[i][k] - 'a' );

    }

    for(int i=0;i<26;i++){

        if(color1[i]==0){

            bool f = dfs1(i);

            if(f){

                cout<<"Impossible"<<endl;

                return 0;

            }

        }

    }

    for(int i=0;i<26;i++){

        if(color[i]==0)

            dfs(i);

    }

    for(int i=ans.size()-1;i>=0;i--)

    {

        char c = ans[i]+'a';

        cout<<c;

    }

    cout<<endl;

    return 0;

}
#include<bits/stdc++.h>

using namespace std;

typedef long long ll;

const int N=3e5+7;

vector<int>mp[N],ans;

int f[N],nxt[N],id[N],vis[N],in[N];

signed main(){

    cin.tie(0);ios::sync_with_stdio(0);

    int n,m,rt=0;cin>>n>>m;

    for(int i=1;i<=n;i++)cin>>f[i];

    for(int i=1,x,y;i<=m;i++)cin>>x>>y,nxt[x]=y,vis[y]=1;

    for(int i=1,x=1;i<=n;i++,x=i)if(!vis[i]){

        while(x){

            id[x]=i;

            if(f[x]&&id[f[x]]!=i)mp[f[x]].push_back(i),in[i]++;

            x=nxt[x];

        }

        if(!in[i])rt=i;

    }

    queue<int>q;

    q.push(rt);

    while(!q.empty()){

        int x=q.front();q.pop();

        while(x){

            ans.push_back(x);

            for(auto i:mp[x]){

                if(!--in[i])q.push(i);

            }

            x=nxt[x];

        }

    }

    if(ans.size()==n)for(auto x:ans)cout<<x<<' ';

    else cout<<0;

}
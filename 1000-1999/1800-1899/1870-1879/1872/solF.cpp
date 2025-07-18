#pragma GCC optimize("Ofast,unroll-loops")
#include <bits/stdc++.h>
using namespace std;
int a[100005],c[100005],in[100005],f[100005];
vector <int> v[100005];
int main() {
    int T = 1, kase = 0;
    cin >> T;
    while (T--) {
        int n;
        scanf("%d",&n);
        for(int i=1;i<=n;i++) in[i]=0,v[i].clear(),f[i]=0;
        for(int i=1;i<=n;i++) scanf("%d",&a[i]),v[i].push_back(a[i]),in[a[i]]++;
        for(int i=1;i<=n;i++) scanf("%d",&c[i]);
        queue <int> q;
        vector <int> ans;
        for(int i=1;i<=n;i++) if(!in[i]) q.push(i),ans.push_back(i),f[i]=1;
        while(!q.empty()){
            int temp=q.front();
            q.pop();
            for(auto t:v[temp]){
                in[t]--;
                if(!in[t]) q.push(t),ans.push_back(t),f[t]=1;
            }
        }
        vector <int> others;
        for(int i=1;i<=n;i++) if(!f[i]) others.push_back(i);
        sort(others.begin(),others.end(),[&](int x,int y){
            return c[x]<c[y];
        });
        for(auto t:others){
            if(f[t]) continue;
            f[t]=1;
            int now=a[t];
            while(now!=t) ans.push_back(now),f[now]=1,now=a[now];
            ans.push_back(t);
        }
        for(int i=1;i<=n;i++) printf("%d%c",ans[i-1],i==n?'\n':' ');
    }
    return 0;
}
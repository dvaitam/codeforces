#include<bits/stdc++.h>
typedef long long ll;
using namespace std;
#define N 400400
vector<int> ans;
int S[N<<1],tot,n,b[N],c[N],dg[N],f[N];
struct Line {
    int v,next;
} e[N<<1];
int h[N<<1],cnt=2;
bool vis[N];
void Add(int u,int v) {
    e[cnt]=(Line) {
        v,h[u]
    };
    h[u]=cnt++;
}
int getf(int x) {
    return x==f[x]?x:f[x]=getf(f[x]);
}
void dfs(int u) {
    for(int &i=h[u]; i; i=e[i].next) {
        if(vis[i>>1])continue;
        vis[i>>1]=true;
        int j=i;
        dfs(e[i].v);
        ans.push_back(e[j].v);
    }
}
int main() {
    scanf("%d",&n);
    for(int i=1; i<n; ++i)scanf("%d",&b[i]);
    for(int i=1; i<n; ++i)scanf("%d",&c[i]);
    for(int i=1; i<n; ++i)if(b[i]>c[i]) {
            puts("-1");
            return 0;
        }
    for(int i=1; i<n; ++i)S[++tot]=b[i];
    for(int i=1; i<n; ++i)S[++tot]=c[i];
    sort(&S[1],&S[tot+1]);
    tot=unique(&S[1],&S[tot+1])-S-1;
    for(int i=1; i<n; ++i)b[i]=lower_bound(&S[1],&S[tot+1],b[i])-S;
    for(int i=1; i<n; ++i)c[i]=lower_bound(&S[1],&S[tot+1],c[i])-S;
    for(int i=1; i<=tot; ++i)f[i]=i;
    for(int i=1; i<n; ++i) {
        dg[b[i]]+=1;
        dg[c[i]]+=1;
        Add(b[i],c[i]);
        Add(c[i],b[i]);
        f[getf(b[i])]=getf(c[i]);
    }
    vector<int> A;
    for(int i=1; i<=tot; ++i)if(dg[i]&1)A.push_back(i);
    if((int)A.size()!=2&&(int)A.size()!=0) {
        puts("-1");
        return 0;
    }
    for(int i=1; i<=tot; ++i)if(getf(i)!=getf(1)) {
            puts("-1");
            return 0;
        }
    if(A.size()==2)dfs(A[0]),ans.push_back(A[0]);
    else dfs(1),ans.push_back(1);
    for(int i=0; i<n; ++i)printf("%d ",S[ans[i]]);
    puts("");
    return 0;
}
#include<cstdio>  
#include<cstring>  
#include<cstdlib>  
#include<iostream>  
#include<algorithm>  
#include<map>
#include<set>
#include<vector>
using namespace std;
#define fo(i,a,b) for(i=a;i<=b;++i)
using namespace std;
const int mn=100100;
int vis[mn],f[mn],a[mn],d[mn];
vector<int>ans1,ans2;
int n,m,i;
char ch;
void read(int &a)
{
    int p=1;
    for(ch=getchar();(ch<'0'||ch>'9') && (ch!='-');ch=getchar());
    if(ch=='-')a=0,p=-1;else a=ch-'0';
    for(ch=getchar();ch>='0'&& ch<='9';ch=getchar())a=a*10+ch-'0';
    a*=p;
}
int dfs(int x)
{
    if(f[x])return f[x];
    if(vis[x])return 0;
    vis[x]=1;
    f[x]=dfs(a[x]);
    return f[x]?f[x]:f[x]=x;
}
int main()
{
    read(n);
    fo(i,1,n)
        read(a[i]),++d[a[i]];
    fo(i,1,n)if(!d[i]||!vis[i])dfs(i);
    fo(i,1,n)
    if(!d[i])
    { 
        ans1.push_back(i),ans2.push_back(f[i]); 
        vis[f[i]]=2; 
    } 
    fo(i,1,n)if(f[i]==i&&vis[i]<2)ans1.push_back(i),ans2.push_back(i); 
    m=ans1.size(); 
    if(m==1&&ans1[0]==ans2[0])m=0; 
    printf("%d\n",m); 
    fo(i,0,m-1)printf("%d %d\n",ans2[i],ans1[(i+1)%m]);
    return 0;
}
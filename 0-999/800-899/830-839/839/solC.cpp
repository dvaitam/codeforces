#include<iostream>
#include<algorithm>
#include<cstring>
#include<cstdio>
#include<queue>
using namespace std;
#define rep(i,l,r) for(int i=l;i<=r;i++)
inline int read(){
    int ret=0,flag=1;char ch=getchar();
    while(ch<'0'||ch>'9'){if(ch=='-')flag=-1;ch=getchar();}
    while(ch>='0'&&ch<='9'){ret=ret*10+ch-'0';ch=getchar();}
    return ret*flag;
}
const double eps=1e-6;
int n;
const int maxn=1e5+10;
int head[maxn];
struct node{
    int next,to;
}e[maxn<<1];
int cnt;
inline void add(int u,int v){
    e[cnt].to=v;
    e[cnt].next=head[u];
    head[u]=cnt++;
}
double dp[maxn];
int vis[maxn];
inline void dfs(int x,int dep){
    int sss=0;
    double sum=0;
    for(int i=head[x];~i;i=e[i].next){
        if(!vis[e[i].to]){
            vis[e[i].to]=1;
            sss++;
            dfs(e[i].to,dep+1);
            sum+=dp[e[i].to];
        }
    }
    //sum+=sss;
    if(sss)
        sum=sum/(double)sss;
    dp[x]=sum;
    if(!sss)
        dp[x]=dep;
}
int main(int argc,const char * argv[]){
    memset(head,-1,sizeof(head));
    n=read();
    int x,y;
    rep(i,1,n-1){
        x=read();y=read();
        add(x,y);
        add(y,x);
    }
    vis[1]=1;
    dfs(1,0);
    printf("%.10lf\n",dp[1]);
    return 0;
}
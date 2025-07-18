#include<set>
#include<map>
#include<cmath>
#include<ctime>
#include<stack>
#include<queue>
#include<vector>
#include<cstdio>
#include<string>
#include<cstdlib>
#include<cstring>
#include<iostream>
#include<algorithm>
#define INF 2147483647
#define LL long long
#define ULL unsigned long long
#define MAXN 1000+100000
#define rep(i,s,t) for(int i=s;i<=t;++i)
#define dwn(i,s,t) for(int i=s;i>=t;--i)

using namespace std;
inline LL read(){
	LL x=0,f=1; char ch=getchar();
	for(;!isdigit(ch);ch=getchar())if(ch=='-')f=-1;
	for(;isdigit(ch);ch=getchar())x=ch-'0'+(x<<3)+(x<<1);
	return x*f;
}

inline void write(LL x){
	if(x<0)putchar('-'),x=-x;
    ULL t=10,len=1; while(t<=x)t*=10,len++;
    while(len--)t/=10,putchar(x/t+48),x%=t;
    return ;
}

int s[MAXN],fir[MAXN],nxt[MAXN<<1],to[MAXN<<1],num,n;
inline void AddEdge(int u,int v){
	nxt[++num]=fir[u],fir[u]=num,to[num]=v; return ;
}

namespace LCT{
	    int rnd[10],num[10],cnt[10];
    int rt,size,lch[10],rch[10],big[10];
    void iniit(){rt=0,size=0;return ;}
    
    inline void update(int u){
        big[u]=big[lch[u]]+big[rch[u]]+cnt[u];
        return ;
    }
    
    inline void lturn(int &u){
        int t=rch[u];rch[u]=lch[t],lch[t]=u;
        big[t]=big[u],update(u),u=t;return ;
    }
    
    inline void rturn(int &u){
        int t=lch[u];lch[u]=rch[t],rch[t]=u;
        big[t]=big[u],update(u),u=t;return ;
    }
    
    void insert(int &cur,int x){
        if(cur==0){
            cur=++size,num[cur]=x;
            cnt[cur]=1,rnd[cur]=rand();
            big[cur]=1; return ;
        }
        big[cur]++; 
        if(num[cur]==x)cnt[cur]++;
        else if(x>num[cur]){
            insert(rch[cur],x);
            if(rnd[rch[cur]]<rnd[cur])lturn(cur);
        }
        else {
            insert(lch[cur],x);
            if(rnd[lch[cur]]<rnd[cur])rturn(cur);
        }
        return ;
    }
    
    void del(int &k,int x){
        if(k==0)return ;
        else if(num[k]==x){
            if(cnt[k]>1){cnt[k]--,big[k]--;return ;}
            if(lch[k]*rch[k]==0){k=lch[k]+rch[k];return ;}
            if(rnd[lch[k]]>rnd[rch[k]])lturn(k),del(lch[k],x);
            else rturn(k),del(rch[k],x);
        }
        else if(num[k]>x)big[k]--,del(lch[k],x);
        else big[k]--,del(rch[k],x);
        update(k); return ;
    }
    
    int frank(int k,int x){
        if(k==0)return 0;
        else if(num[k]==x)return big[lch[k]]+1;
        else if(num[k]>x)return frank(lch[k],x);
        else return big[lch[k]]+cnt[k]+frank(rch[k],x);
    }
    
    int srank(int k,int x){
        if(k==0)return 1;
        else if(x<=big[lch[k]])return srank(lch[k],x);
        else if(x>big[lch[k]]+cnt[k])
            return srank(rch[k],x-big[lch[k]]-cnt[k]);
        else return num[k];
    }
    
    int bef(int k,int x){
        if(k==0)return -INF;
        else if(num[k]>=x)return bef(lch[k],x);
        else return max(num[k],bef(rch[k],x));
    }
    
    int aft(int k,int x){
        if(k==0)return INF;
        else if(num[k]<=x)return aft(rch[k],x);
        else return min(num[k],aft(lch[k],x));
    }
}
using namespace LCT;

void dfs(int x,int fa,int flg=1){
	for(int i=fir[x];i;i=nxt[i])if(to[i]!=fa)
		{flg=0;dfs(to[i],x);s[x]+=s[to[i]];}
	if(flg)s[x]=1; return ;
}

void init(){
	n=read();
	rep(i,2,n){
		int y=read();
		AddEdge(i,y),AddEdge(y,i);
    }
	return ;
}

void work(){
	dfs(1,0); sort(s+1,s+n+1);
	rep(i,1,n)write(s[i]),putchar(' ');
    return ;
}

int main(){
    init(); work(); return 0;
}
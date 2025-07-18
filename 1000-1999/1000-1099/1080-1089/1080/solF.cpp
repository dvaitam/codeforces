#include<iostream>
#include<cstring>
#include<cmath>
#include<algorithm>
#include<cstdio>
#include<iomanip>
#include<cstdlib>
#define MAXN 0x7fffffff
typedef long long LL;
const int N=1e5+5;
using namespace std;
inline int Getint(){register int x=0,f=1;register char ch=getchar();while(!isdigit(ch)){if(ch=='-')f=-1;ch=getchar();}while(isdigit(ch)){x=x*10+ch-'0';ch=getchar();}return x*f;}

struct Que{
    int l,r,p;
    bool operator < (const Que & a)const{return r<a.r;}
}sav_que[N*3];

struct Seg{int lson,rson,mn;}Tree[N*100];
int cnt,rt[N*3];
int data[N*3];

inline void Pushup(int v){Tree[v].mn=min(Tree[Tree[v].lson].mn,Tree[Tree[v].rson].mn);}
void Newcode(int &v,int x){
	v=++cnt;
	Tree[v].lson=Tree[x].lson,
	Tree[v].rson=Tree[x].rson,
	Tree[v].mn=Tree[x].mn;
}

int Update(int v,int l,int r,int pos,int val){
	int u;
	Newcode(u,v);
	if(l==r){Tree[u].mn=max(Tree[u].mn,val);return u;}
	int mid=l+r>>1;
	if(pos<=mid)Tree[u].lson=Update(Tree[u].lson,l,mid,pos,val);
	else Tree[u].rson=Update(Tree[u].rson,mid+1,r,pos,val);
	Pushup(u);
	return u;
}

int Query(int v,int l,int r,int L,int R){
	if(l>R||r<L)return MAXN;
	if(l>=L&&r<=R)return Tree[v].mn;
	int mid=l+r>>1;
	return min(Query(Tree[v].lson,l,mid,L,R),Query(Tree[v].rson,mid+1,r,L,R));
}

int main(){
    int n=Getint(),m=Getint(),K=Getint();
    
	for(register int i=1;i<=K;++i){
		sav_que[i].l=Getint();
		sav_que[i].r=data[i]=Getint();
		sav_que[i].p=Getint();
	}
    sort(data+1,data+K+1),sort(sav_que+1,sav_que+K+1);
    int cc=unique(data+1,data+K+1)-data-1;
    
	for(register int i=1;i<=K;++i){
	    int x=lower_bound(data+1,data+cc+1,sav_que[i].r)-data;
	    if(sav_que[i].r==sav_que[i-1].r)rt[x]=Update(rt[x],1,n,sav_que[i].p,sav_que[i].l);
	    else rt[x]=Update(rt[x-1],1,n,sav_que[i].p,sav_que[i].l);
	}
    
	for(register int i=1;i<=m;++i){
	    int a=Getint(),b=Getint(),x=Getint(),y=Getint();
	    int p=upper_bound(data+1,data+cc+1,y)-data-1;
	    cout<<(Query(rt[p],1,n,a,b)>=x?"yes":"no")<<'\n';
	    fflush(stdout);
	}
	return 0;
}
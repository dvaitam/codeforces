#include <bits/stdc++.h>
#include <bits/stdc++.h>
using namespace std;
#define iinf 2000000000
#define linf 1000000000000000000LL
#define ulinf 10000000000000000000ull
#define MOD1 1000000007LL
#define mpr make_pair
typedef long long LL;
typedef unsigned long long ULL;
typedef unsigned long UL;
typedef unsigned short US;
typedef pair < int , int > pii;
clock_t __stt;
inline void TStart(){__stt=clock();}
inline void TReport(){printf("\nTaken Time : %.3lf sec\n",(double)(clock()-__stt)/CLOCKS_PER_SEC);}
template < typename T > T MIN(T a,T b){return a<b?a:b;}
template < typename T > T MAX(T a,T b){return a>b?a:b;}
template < typename T > T ABS(T a){return a>0?a:(-a);}
template < typename T > void UMIN(T &a,T b){if(b<a) a=b;}
template < typename T > void UMAX(T &a,T b){if(b>a) a=b;}
int B,n,q,h[10000005],cnt;
LL c[10000005];
struct block{
	vector < int > hs;
	vector < LL > cs;
	void read(){
		int L,i,j,k;
		scanf("%d",&L);
		hs.resize(L);
		cs.resize(L);
		for(i=0;i<L;++i){
			scanf("%d",&hs[i]);
		}
		for(i=0;i<L;++i){
			scanf("%d",&k);
			cs[i]=(LL)k;
		}
	}
	void apply(LL mul){
		int i,j,k;
		for(i=0;i<(int)hs.size();++i){
			h[cnt+i]=hs[i];
		}
		for(i=0;i<(int)cs.size();++i){
			c[cnt+i]=cs[i]*mul;
		}
		cnt+=i;
	}
}bs[250005];
int lb[10000005],rb[10000005],stk[10000005],sz;
LL dp[10000005];
int main(){
    // inputting start
    // 数据结构记得初始化！ n，m别写反！
    int i,j,k;
	scanf("%d%d",&B,&n);
	for(i=0;i<B;++i){
		bs[i].read();
	}
	scanf("%d",&q);
	for(i=0;i<q;++i){
		int id,mul;
		scanf("%d%d",&id,&mul);
		bs[--id].apply(mul);
	}
    #ifdef LOCAL
        TStart();
    #endif
    // calculation start
    // 数据结构记得初始化！ n，m别写反！
    sz=0;
	for(i=0;i<n;++i){
		int L=i-h[i]+1;
		while(sz){
			if(stk[sz-1]-h[stk[sz-1]]+1<L) break;
			--sz;
		}
		if(sz && stk[sz-1]>=L){
			lb[i]=MIN(lb[stk[sz-1]],MAX(L,0));
		}
		else{
			lb[i]=MAX(L,0);
		}
		stk[sz++]=i;
	}
	sz=0;
	for(i=n-1;i>=0;--i){
		int R=i+h[i]-1;
		while(sz){
			if(stk[sz-1]+h[stk[sz-1]]-1>R) break;
			--sz;
		}
		if(sz && stk[sz-1]<=R){
			rb[i]=MAX(rb[stk[sz-1]],MIN(R,n-1));
		}
		else{
			rb[i]=MIN(R,n-1);
		}
		stk[sz++]=i;
	}
	sz=0;
	for(i=0;i<n;++i){
		dp[i]=linf;
		while(sz && rb[stk[sz-1]]<i) --sz;
		int d=(sz?stk[sz-1]:-1);
		if(sz && rb[d]>=i){// covered from left
			UMIN(dp[i],(LL)(d?(LL)dp[d-1]:0ll)+c[d]);
		}
		{// cover to left
			UMIN(dp[i],(LL)(lb[i]?(LL)dp[lb[i]-1]:0ll)+c[i]);
		}
		LL C=(LL)(i?dp[i-1]:0ll)+c[i];
		if(sz && rb[d]==rb[i]){
			if(C<(LL)(d?dp[d-1]:0ll)+c[d]){
				stk[sz-1]=i;
			}
		}
		else if(!sz || C<(LL)(d?dp[d-1]:0ll)+c[d]){
			stk[sz++]=i;
		}
	}
	printf("%lld\n",dp[n-1]);
    #ifdef LOCAL
        TReport();
    #endif
    return 0;
}
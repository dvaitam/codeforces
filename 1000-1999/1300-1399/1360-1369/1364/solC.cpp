#include<cstdio>
#include<cmath>
#include<cstring>
#include<algorithm>
using namespace std;
typedef long long ll;
const int mo=998244353;
#define M 100005

template <class T>
void rd(T &l){
	l=0;int f=1;char ch=getchar();
	while(ch<'0' || ch>'9'){ if(ch=='-') f=0; ch=getchar(); }
	while(ch>='0' && ch<='9'){ l=l*10+(ch^'0'); ch=getchar(); }
	l=(f?l:-l);
}

int sl[M];
int ans[M];

int main(){
    // freopen("21.in","r",stdin);
    // freopen("21.out","w",stdout);
    int n,m,x,y=0,f=1,ct=0;
    rd(n);
    for(int i=1;i<=n;++i){
        rd(x);
        if(x>i) f=0;
        if(f){
            sl[++ct]=i;
            while(y<x){
                ans[sl[ct--]]=y;
                ++y;
            }
        }
    }
    while(ct) ans[sl[ct--]]=x+1;
    if(f){
        for(int i=1;i<=n;++i) printf("%d ",ans[i]);
    }
    else puts("-1");
	return 0;
}
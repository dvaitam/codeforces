#include<cstdio>
#include<algorithm>
#include<cstring>
using namespace std;
#define M 200005
typedef long long ll;
const int mo=1e9+7;
int A[M],B[M];
inline void rd(int &l){
	l=0;int f=0;char ch=getchar();
	while(ch>'9'||ch<'0'){if(ch=='-')f=1;ch=getchar();}
	while(ch>='0'&&ch<='9'){l=l*10+ch-'0';ch=getchar();}
	if(f)l=-l;
}
int main(){
	int t,n,k,f1,f2;
	rd(t);
	while(t--){
		rd(n);k=f1=f2=0;
		for(int i=1;i<=n;++i)B[i]=0;
		for(int i=1;i<=n;++i)rd(A[i]),++B[A[i]],k=max(k,A[i]);
		for(int i=1;i<=k;++i)
			if(B[i]==1){
				f1=1;
			}
			else if(B[i]==2){
				if(f1){f2=1;break;}
			}
			else{f2=1;break;}
		if(f2){puts("0");continue;}
		if(n==k){printf("2\n0 %d\n%d 0\n",n,n);continue;}
		f1=f2=0;
		for(int i=1;i<=k;++i)B[i]=0;
		for(int i=1;i<=k;++i)++B[A[i]];
		for(int i=1;i<=k;++i)if(!B[i]){f1=1;break;}
		for(int i=1;i<=k;++i)B[i]=0;
		for(int i=n;i>n-k;--i)++B[A[i]];
		for(int i=1;i<=k;++i)if(!B[i]){f2=1;break;}
		printf("%d\n",(!f1)+((!f2)&&(f1||k!=n-k)));
		if(!f1)printf("%d %d\n",k,n-k);
		if(!f2){if(f1||k!=n-k)printf("%d %d\n",n-k,k);}
	}
	return 0;
}
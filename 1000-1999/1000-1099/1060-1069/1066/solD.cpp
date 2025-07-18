#include<cstdio>
#include<algorithm>
using namespace std;
char ch; bool fh;
inline void read(int &a){
	for(fh=0,ch=getchar();ch<'0'||ch>'9';ch=getchar())if(ch=='-')fh=1;
	for(a=0;ch>='0'&&ch<='9';ch=getchar())(a*=10)+=ch-'0'; if(fh)a=-a;
}
int n,m,k;
int A[200010];
int pos;
bool check(int dolen){
	int vi,tmp=0,sum=0;
	for(vi=dolen;vi<=n;vi++){
		tmp+=A[vi];
		if(tmp>k){
			sum++; tmp=A[vi];
		}
	}
	sum++;
	return (sum<=m);
}
int main(){
	int vi;
//	freopen("T.in","r",stdin);
	read(n); read(m); read(k);
	for(vi=1;vi<=n;vi++)read(A[vi]);
	int l=pos+1,r=n,mid;//��ʾװ��[mid..n]������ 
	for(;l!=r;){
		mid=(l+r)>>1;
		if(check(mid))r=mid;else l=mid+1;
	}
	printf("%d\n",n-l+1);
}
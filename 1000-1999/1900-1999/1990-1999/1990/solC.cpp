#include<bits/stdc++.h>
using namespace std;
int main() {
	long long n,t,a,b,c,p1[200005],p2[200005];
	long long sum;
	scanf("%lld",&t);
	while(t--) {
		scanf("%lld",&n);
		sum=a=b=c=0;
		memset(p1,0,sizeof(p1));
		memset(p2,0,sizeof(p2));
		while(n--) {
			scanf("%lld",&a);
			if(p1[a]) b=max(b,a);
			else p1[a]=1;
			if(p2[b]) c=max(c,b);
			else if(b) p2[b]=1;
			sum+=a+b+(n+1)*c;
		}
		printf("%lld\n",sum);
	}
}
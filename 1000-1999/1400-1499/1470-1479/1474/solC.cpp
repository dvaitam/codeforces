#include<cstdio>
#include<set>
#include<algorithm>
using namespace std;
int t,n,a[2005],b[1000005],ans[2005][2];
int main() {
	scanf("%d",&t);
	while(t--) {
		scanf("%d",&n);n*=2;
		for(int i=1;i<=n;++i) {scanf("%d",&a[i]);b[a[i]]++;}
		sort(a+1,a+1+n);
		int now,tf;
		for(int i=1;i<n;++i) {
			tf=a[n]+a[i];now=n;
			for(int j=1;j<=n/2;++j) {
				while(!b[a[now]]) now--;
				b[a[now]]--;
				if(b[tf-a[now]]) {
					b[tf-a[now]]--;
					ans[j][0]=tf-a[now];
					ans[j][1]=a[now];
					tf=a[now];
				}else goto lp;
			}
			printf("YES\n%d\n",a[n]+a[i]);
			for(int j=1;j<=n/2;++j) printf("%d %d\n",ans[j][0],ans[j][1]);
			goto lp2;
			lp:;
			for(int i=1;i<=n;++i) b[a[i]]=0;
			for(int i=1;i<=n;++i) b[a[i]]++;
		}
		printf("NO\n");
		lp2:;
		for(int i=1;i<=n;++i) b[a[i]]=0;
	}
	return 0;
}
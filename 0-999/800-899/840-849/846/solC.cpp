//Author: Lixiang
#include<stdio.h>
const int maxn=5001;
struct C{
	long long s[maxn];
	int a[maxn],N;
	void init(){
		scanf("%d",&N);
		for(int i=1;i<=N;i++){
			scanf("%d",&a[i]);
			s[i]=s[i-1]+a[i];
		}
	}
	void work(){
		long long sum=0,now=0,ans=-5000000000000ll;
		int st=1,ed,St=1,Ed=0,m0,m1,m2;
		for(int i=1;i<=N;i++){
			sum=sum+a[i];
			if(sum>=0){
				sum=0;
				st=i+1;
			}
			else{
				ed=i;
				if(sum<now){
					St=st;
					Ed=ed;
					now=sum;
				}
			}
			if(((s[i]-now)*2-s[N])>ans){
				ans=(s[i]-now)*2-s[N];
				m0=St-1;
				m1=Ed;
				m2=i;
			}
		}
		printf("%d %d %d\n",m0,m1,m2);
	}
}sol;
int main(){
	sol.init();
	sol.work();
	return 0;
}
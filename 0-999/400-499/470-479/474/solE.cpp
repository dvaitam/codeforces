#include<iostream>
#include<cstdio>

using namespace std;

long long maxn[100005],minn[100005],ans[100005];
int maxnn[100005],minnn[100005],tail,ans1[100005];

int main(){
	int n;
	long long x,d;
	scanf("%d %I64d",&n,&d);
	scanf("%I64d",&x);
	maxn[tail]=x;minn[tail]=x;
	maxnn[0]=1;minnn[0]=1;
	for(int i=2;i<=n;i++){
		scanf("%I64d",&x);
		if(minn[tail]+d<=x||maxn[tail]-d>=x){
			maxn[++tail]=x;maxnn[tail]=i;
			minn[tail]=x;minnn[tail]=i;
		}
		else	if(tail){
				if(minn[tail-1]+d<=x||maxn[tail-1]-d>=x){
					if(maxn[tail]<x){
						maxn[tail]=x;
						maxnn[tail]=i;
					}
					if(minn[tail]>x){
						minn[tail]=x;
						minnn[tail]=i;
					}
				}
			}
			else{
				if(maxn[tail]<x){
					maxn[tail]=x;
					maxnn[tail]=i;
				}
				if(minn[tail]>x){
					minn[tail]=x;
					minnn[tail]=i;
				}
			}
	}
	printf("%d\n",tail+1);
	ans[tail]=minn[tail];
	ans1[tail]=minnn[tail];
	for(int i=tail-1;i>=0;i--)
		if(minn[i]+d<=ans[i+1]){
			ans[i]=minn[i];
			ans1[i]=minnn[i];
		}
		else{
			ans[i]=maxn[i];
			ans1[i]=maxnn[i];
		}
	for(int i=0;i<tail;i++)
	printf("%d ",ans1[i]);
	printf("%d\n",ans1[tail]);
	return 0;
}
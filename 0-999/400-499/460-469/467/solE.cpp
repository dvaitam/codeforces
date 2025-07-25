#include <cstdio>
#include <cstdio>
#include <vector>
#include <algorithm>
using namespace std;
const int N=500000;
int a[N];
int c[N];
int q[N];
int fly[N];
int co[N];
int mark[N];
typedef pair<int,int> PII;
int st[N];
int stop;
void pop(){
	for(int i=0;i<stop;i++){
		co[st[i]]=0;
		mark[st[i]]=-1;
	}
	stop=0;
}
bool cmp(int a,int b){
	return c[a]<c[b];
}
int main(){
	int n;
	scanf("%d",&n);
	vector<PII> ans;
	for(int i=0;i<n;i++)scanf("%d",&c[i]);
	for(int i=0;i<n;i++)fly[i]=i-1,q[i]=i;
	for(int i=0;i<=n;i++)mark[i]=-1;
	sort(q,q+n,cmp);
	int now=1;
	for(int i=0;i<n;i++){
		if(i==0 || c[q[i]]==c[q[i-1]])a[q[i]]=now;
		else{
			now++;
			a[q[i]]=now;
		}
	}
	for(int i=0;i<n;i++){
		st[stop++]=a[i];
		co[a[i]]++;
		if(mark[a[i]]!=-1){
			ans.push_back(make_pair(mark[a[i]],i));
			pop();
		}
		else if(co[a[i]]==4){
			ans.push_back(make_pair(i,i));
			pop();
		}
		else if(co[a[i]]>=2){
			for(int j=i-1;j>=0;){
				if(a[i]==a[j]){fly[i]=j;break;}
				mark[a[j]]=i;j=fly[j];
			}
		}
	}
	printf("%d\n",(int)ans.size()*4);
	for(int i=0;i<ans.size();i++){
		int a=c[ans[i].first], b = c[ans[i].second];
		printf("%d %d %d %d",a,b,a,b);
		if(i==(int)ans.size()-1)puts("");
		else printf(" ");
	}
}
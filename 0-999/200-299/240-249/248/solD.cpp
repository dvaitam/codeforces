// LUOGU_RID: 150904603
#include<bits/stdc++.h>
using namespace std;
int n,t,a[500005];
char s[500005];
bool chk(int x){
	if(a[n]+x<0) return 0;
	int to=0,res=0;
	for(int i=1;i<=n;i++) if(s[i]=='H'||a[i-1]+x<0) to=i;
	for(int i=0;i<=n;i++){
		if(res+max(to-i,0)*2<=t) return 1;
		res+=(a[i]+x<0)*2+1;
	}
	return 0;
}
int main(){
	cin>>n>>t>>s+1;
	for(int i=1;i<=n;i++) a[i]=a[i-1]+(s[i]=='S')-(s[i]=='H');
	int l=0,r=n+1;
	while(l<r){
		int mid=l+r>>1;
		if(chk(mid)) r=mid;
		else l=mid+1;
	}
	if(l<=n) cout<<l<<endl;
	else cout<<"-1\n";
	return 0;
}

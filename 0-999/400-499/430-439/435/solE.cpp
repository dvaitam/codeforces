#include <bits/stdc++.h>
using namespace std;
typedef long long ll;
typedef unsigned int uint;
typedef unsigned long long ull;
const ll mod=998244353;
const ll inf=2e9;
const int N=2e5+5;
const int M=2e5+5;
const int V=2e6+5;
int n,m;
string s[1005],t[1005];
bool work(){
	int e[5]={-1,-1,-1,-1,-1};
	for(int i=1;i<=n;i++){
		int g[5]={-1,-1,-1,-1,-1};
		for(int j=1;j<=m;j++) if(s[i][j]!='0'){
			if(g[s[i][j]-'0']==!(j&1)||e[s[i][j]-'0']==!(i&1)) return 0;
			g[s[i][j]-'0']=(j&1),e[s[i][j]-'0']=(i&1);
		}
	}
	if((e[1]==0)+(e[2]==0)+(e[3]==0)+(e[4]==0)>2||(e[1]==1)+(e[2]==1)+(e[3]==1)+(e[4]==1)>2) return 0;
	int s1=0,s2=0,s3=0,s4=0;
	for(int i=1;i<=4;i++) if(e[i]!=-1){
		if(e[i]==1){
			if(s1) s2=i;
			else s1=i;
		}
		else{
			if(s3) s4=i;
			else s3=i;
		}
	}
	for(int i=1;i<=4;i++) if(e[i]==-1){
		if(!s1) s1=i;
		else if(!s2) s2=i;
		else if(!s3) s3=i;
		else s4=i;
	}
	for(int i=1;i<=n;i++){
		for(int j=1;j<=m;j++) if(s[i][j]!='0'){
			if((i&1)&&(j&1)&&s1!=s[i][j]-'0'){swap(s1,s2);break;}
			if((i&1)&&!(j&1)&&s2!=s[i][j]-'0'){swap(s1,s2);break;}
			if(!(i&1)&&(j&1)&&s3!=s[i][j]-'0'){swap(s3,s4);break;}
			if(!(i&1)&&!(j&1)&&s4!=s[i][j]-'0'){swap(s3,s4);break;}
		}
		for(int j=1;j<=m;j++){
			if((i&1)&&(j&1)) s[i][j]=char(s1+'0');
			if((i&1)&&!(j&1)) s[i][j]=char(s2+'0');
			if(!(i&1)&&(j&1)) s[i][j]=char(s3+'0');
			if(!(i&1)&&!(j&1)) s[i][j]=char(s4+'0');
		}
	}
	return 1;
}
void solve(int Ca){
	cin>>n>>m;
	for(int i=1;i<=n;i++) cin>>s[i],s[i]=" "+s[i];
	if(work()){
		for(int i=1;i<=n;i++,cout<<"\n") for(int j=1;j<=m;j++) cout<<s[i][j];
		return;
	}
	for(int i=1;i<=m;i++){
		t[i]=" ";
		for(int j=1;j<=n;j++) t[i]+=s[j][i];
	}
	swap(n,m);
	for(int i=1;i<=n;i++) s[i]=t[i];
	if(work()){
		for(int i=1;i<=m;i++,cout<<"\n") for(int j=1;j<=n;j++) cout<<s[j][i];
		return;
	}
	cout<<"0\n";
}
int main(){
	#ifdef ONLINE_JUDGE
	ios::sync_with_stdio(false); cin.tie(nullptr); cout.tie(nullptr);
	#endif
	#ifndef ONLINE_JUDGE
	freopen("test.in","r",stdin);
	freopen("test.out","w",stdout);
	#endif
	
	int Ca=1;
//	cin>>Ca;
	for(int i=1;i<=Ca;i++){
		solve(i);
	}
	return 0;
}
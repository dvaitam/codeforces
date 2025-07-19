#include<bits/stdc++.h>
using namespace std;
int main(){
	int T;
	int n,q,l,r;
	string a,b;
	cin>>T;
	while(T--){
		cin>>n>>q;
		cin>>a>>b;
		vector<vector<int> >s(n+1,vector<int>(26,0));
		for(int i=0;i<n;i++){
			s[i+1]=s[i];
			s[i+1][a[i]-'a']++;
			s[i+1][b[i]-'a']--;
		}
		while(q--){
			cin>>l>>r;
			l--;
			int ans=0;
			for(int i=0;i<26;i++){
				ans+=max((int)0,s[r][i]-s[l][i]);
			}
			cout<<ans<<"\n";
		}
	}
	return 0;
}
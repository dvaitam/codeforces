#include <bits/stdc++.h>

using namespace std;
const int mod = 1000000007;
const int N = 200005;

int main(){
	ios_base::sync_with_stdio(false);
    cin.tie(NULL);
    cout.tie(NULL);
#ifndef ONLINE_JUDGE
    freopen("in.txt", "r", stdin);
    freopen("out.txt", "w", stdout);
#endif
		string s;
		cin>>s;
		int n = s.size();
		int k;
		cin>>k;
		int c1=0,c2=0;
		for(auto i : s){
			if(i == '?')
				c1++;
			if(i == '*')
				c2++;
		}
		int ss = (int)s.size();
		ss -= 2*(c1 + c2);
		if(k < ss){
			cout<<"Impossible";
			return 0;
		}
		ss += c1 + c2;
		if(c2 <= 0 && k > ss){
			cout<<"Impossible";
			return 0;
		}
		if(k <= ss){
			//cout<<"123";
			int dif = ss - k;
			string ans ="";
			for(int i=0;i<n -1;i++){
				if(s[i] == '*' || s[i] == '?')
					continue;
				if(s[i+1] != '*' && s[i+1] != '?'){
					ans += s[i];
					//cout<<ans<<"\n";
				}
				if(dif > 0 && ( s[i+1] == '*' || s[i+1] == '?' ) ){
					dif--;
					continue;
				}
				if(dif == 0 && ( s[i+1] == '*' || s[i+1] == '?' ) )
				ans += s[i];
			}
			if(s[n-1] != '*' && s[n-1] != '?')
				ans += s[n-1];
			cout<<ans;
			return 0;
		}
		if(k > ss){
			//cout<<"456";
			int dif = k - ss;
			string ans ="";
			for(int i=0;i<n -1;i++){
				if( s[i] == '?')
					continue;
				if(s[i] == '*' && c2 > 1){
					c2--;
					continue;
				}
				if(s[i] == '*' && c2 == 1){
					c2--;
					for(int j=0;j<dif;j++){
						ans += s[i-1];
					}
					continue;
				}
					ans += s[i];
	
			}
			if(s[n-1] != '*' && s[n-1] != '?')
				ans += s[n-1];
			if(s[n-1] == '*' && c2 == 1){
					c2--;
					for(int j=0;j<dif;j++){
						ans += s[n-2];
					}
				}
			cout<<ans;
			return 0;
		}
	return 0;
}
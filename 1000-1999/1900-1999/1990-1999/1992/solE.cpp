#include<bits/stdc++.h>
using namespace std;
int t; 

int main(){
	ios::sync_with_stdio(0),cin.tie(0),cout.tie(0);
	
	cin>>t;
	while(t--){
		int n;
		cin>>n;
		
		string s=to_string(n);
		int len=s.size();
		for(int i=0;i<8;i++) s+=s;
		
		vector<pair<int,int>>q;
		for(int a=1;a<=10000;a++){
			for(int b=a*len-7;b<a*len;b++){
				if(b>=1&&b<=10000&&n*a-b==stoi(s.substr(0,a*len-b))){
					q.push_back({a,b});
				}
			}
		}
		cout<<q.size()<<'\n';
		for(auto &[a, b]:q) cout<<a<<' '<<b<<'\n';
	}
	
	return 0;
}
#include<bits/stdc++.h>
using namespace std;
int main(){
	int t;
	cin>>t;
	while(t--){
		string s;
		cin>>s;
		int i=1;
		while(s[i-1]!=s[i]&&i<s.length()){
			i++;
		}
		if(s[i-1]=='a')
		s.insert(i,1,'b');
		else
		s.insert(i,1,'a');
		cout<<s<<endl;
	}
}
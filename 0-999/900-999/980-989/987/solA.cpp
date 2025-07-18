#include <bits/stdc++.h>
using namespace std;
#define ll long long

int main(){
	string s;
	map<string,string> jem;

	jem.insert({"green","Time"});
	jem.insert({"blue","Space"});
	jem.insert({"purple","Power"});
	jem.insert({"orange","Soul"});
	jem.insert({"red","Reality"});
	jem.insert({"yellow","Mind"});
	int n;
	cin>>n;

	for(int i=0;i<n;i++){
		cin>>s;
		jem.erase(s);
	}
	cout<<jem.size()<<endl;
	for(auto k:jem){
		cout<<k.second<<endl;
	}
}
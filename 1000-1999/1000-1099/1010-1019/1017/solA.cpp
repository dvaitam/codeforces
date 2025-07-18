#include<bits/stdc++.h>
using namespace std;
int main(){
	int a;
	cin>>a;
	vector<long long>vec;
	for(int i=0;i<a;i++){
		int g,h,k,j;
		cin>>g>>h>>k>>j;
		vec.push_back(g+h+k+j);
	}
	int u=vec[0];
	sort(vec.begin(),vec.end());
	int y=0;
	for(int i=vec.size()-1;i>=0;i--){
		y++;
		if(vec[i]==u){
			cout<<y<<endl;
			return 0;
		}
	}
}
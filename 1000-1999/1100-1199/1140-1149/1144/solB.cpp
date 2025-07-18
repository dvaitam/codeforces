#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
#include <cmath>
#include <climits>
#include <map>

//author:Swastik Banerjee


using namespace std;

typedef long long ll;

int main() {

	ios_base::sync_with_stdio(false);
	cin.tie(0);

#ifndef ONLINE_JUDGE
	freopen("input.txt","r",stdin);
	freopen("output.txt","w",stdout);
#endif

int n;
cin>>n;
int arr[n];
vector <int> v1;
vector <int>v2;
for(int i=0;i<n;i++){
	cin>>arr[i];
}
for(int j=0;j<n;j++){
	if(arr[j]%2==0)
		v1.push_back(arr[j]);
	else
		v2.push_back(arr[j]);
}
sort(v1.begin(),v1.end());
sort(v2.begin(),v2.end());
int x=v1.size()-v2.size();
int sum=0;
if(abs(x)==1 || x==0)
	cout<<"0"<<endl;
else if(x>=2){
	for(int i=0;i<x-1;i++){
		sum+=v1[i];
	}
	cout<<sum<<endl;
}
else if(x<=-2){
	for(int i=0;i<abs(x)-1;i++){
		sum+=v2[i];
	}
	cout<<sum<<endl;
}




}
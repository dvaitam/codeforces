#include <bits/stdc++.h>

using namespace std;
vector<int >v;
int n,k;
void find(int i){

	int be=1;
	i-=(k+1)*2;
	be+=min(i,k);

	v.push_back(be);
	be+=(2*k+1);
	//cout<<be<<endl;
	while(be<=n){

		v.push_back(be);
		be+=(2*k+1);
	}

	cout<<v.size()<<endl;
	for(int i=0;i<(int)v.size();++i){
		cout<<v[i]<<" ";

	}
	cout<<endl;
}

int main() {

	cin>>n>>k;
	if(n<(k+1)*2){
	cout<<1<<endl;
	cout<<min(n,k+1)<<endl;
	return 0;
	}

	for(int i=(k+1)*2;i<=(2*k+1)*2;++i){
		if((n-i)%(2*k+1)==0){
		find(i);
		return 0;
		}
	}
	cout<<n<<endl;
	for(int i=1;i<=n;++i){
		cout<<i<<" ";
	}
	cout<<endl;

}
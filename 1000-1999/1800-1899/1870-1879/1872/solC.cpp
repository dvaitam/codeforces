#include <iostream>
#include <algorithm>
#include <cmath>
#include <map>
#include <string>
using namespace std;
int gcd(int a,int b) {
	return b?gcd(b,a%b):a;
}
int main() {
	int t=0;
	cin>>t;
	while(t--){
		long long l=0,r=0;
		cin>>l>>r;
		for(int i=2;i*i<=r;++i){
			int k=r/i;
			k*=i;
			if(k<l)
				continue;
			if(min(i,k-i)!=1){
				cout<<i<<" "<<k-i<<endl;
				goto aaa;
			}
		}
		cout<<"-1"<<endl;;
		aaa:
			continue;
	}
    return 0;
}
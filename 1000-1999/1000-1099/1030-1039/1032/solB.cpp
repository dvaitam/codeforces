#include <iostream>
using namespace std;

int main() {
	int n,k,a[5]={0},sm=0,mx=0;
	
	string s,r[5];
	cin>>s;
	n=s.size();
	k=n/20;
	if(n%20)k++;
	for(int i=0;i<k;i++){
	a[i]=n/k;
	sm+=a[i];
	}
	mx=n/k;
	int i=0;
	while(sm<n){
		sm++;
		a[i]++;
		if (mx<a[i])mx=a[i];
		i++;
	}
	int p=0;
	for (i=0;i<k;i++){
		for(int t=0;t<a[i];t++)
		r[i]+=s[t+p];
	
	p+=a[i];
	if (mx>a[i])r[i]+="*";
	}
	cout<<k<<" "<<mx<<endl;
	for (i=0;i<k;i++)cout<<r[i]<<endl;
	return 0;
}
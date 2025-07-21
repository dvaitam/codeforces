#include<iostream>
using namespace std;
void solve(){
	int n;
	cin>>n;
	if(n<=5){
		cout<<n/2+1<<endl;
		for(int i=1;i<=n;i++) cout<<i/2+1<<' ';
		cout<<endl;
		return;
	}
	cout<<4<<endl;
	for(int i=1;i<=n;i++)
		cout<<i%4+1<<' ';
	cout<<endl;
}
int main(){
	int T;
	cin>>T;
	while(T--) solve();
}

#include<bits/stdc++.h>
using namespace std;
int main(){
int n;
cin>>n;
int i,div,q;
for( i=9;i>0;i--){
	if(n%i==0){
		q=n/i;
		div=i;
		break;
	}
}
if(i==1){
	cout<<n<<endl;
	for(int j=0;j<n;j++){
		cout<<1<<" ";
	}
}
else{
	cout<<q<<endl;
	for(int j=0;j<q;j++){
		cout<<div<<" ";
	}
}
}
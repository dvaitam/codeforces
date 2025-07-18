#include <bits/stdc++.h>
using namespace std;

string a, b;

int main(){

	int n, m;
	cin>>n>>m;
	
	if(n<m){

		for(int i = 0;i<n;i++)
			a+='1', b+='1';
		
		if(2*n-m>=0)
			for(int d=2*n-m;d>0;d--)
				b+='0';
		if(2*n-m<0)
			for(int d=m-2*n;d>0;d--)
				b = '1'+b;

		cout<<b<<endl<<a;

	}else{

		for(int i = 0;i<n;i++)
			a+='1';
		for(int i =0;i<m-1;i++)
			b+='9';
		for(int i=0;i<n-m;i++)
			b+='8';
		b+='9';
		cout<<b<<endl<<a;	

	}
		


			
		




}
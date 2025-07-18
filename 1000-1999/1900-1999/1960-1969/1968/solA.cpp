#include <stdio.h>
#include <bits/stdc++.h>
using namespace std;
#define rep(i,n) for (int i = 0; i < (n); ++i)
#define Inf32 1000000001
#define Inf64 1000000000000000001

int main(){
	
	int _t;
	cin>>_t;
	
	rep(_,_t){
		
		int x;
		cin>>x;
		int ma = 1;
		for(int i=1;i<x;i++){
			if(gcd(i,x)+i > gcd(ma,x)+ma)ma = i;
		}
		
		cout<<ma<<endl;
	}
	
	return 0;
}
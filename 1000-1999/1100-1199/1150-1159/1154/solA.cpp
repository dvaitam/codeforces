#include <bits/stdc++.h>
#include <cmath>
#include <iostream>
using namespace std ;
int main (){
	int n[4];
	for(int i =0;i<4;i++)
	cin>>n[i];
	sort(n, n+4);
	for (int i =0;i <3;i++)
	cout<<n[3]-n[i]<<" ";
	
}
#include<iostream>

#include<math.h>

using namespace std;

int a[200];

int main(){



	int n,k;

	cin >> n >> k;

	n= n*2 +1;



	for(int i= 0 ; i < n ; i++)		cin>> a[i];

	

	for(int i=0 ; i < n ; i++){

		if(i % 2 == 1 && k>0 ){

			if(a[i]-1 > a[i-1] && a[i]-1 > a[i+1] ){

				k--;

				cout<<a[i]-1<<" ";

			}

			else

				cout<<a[i]<<" ";

		}

		else

			cout<<a[i]<<" ";

	}





}
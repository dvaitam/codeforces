#include <bits/stdc++.h>

using namespace std;

int n , x[49][49] , oddc = 1  , evenc = 2 , flag = 0 ;

int main() {

	cin >> n ;

	for(int i = 0 ; i < n ; i++ ){x[i][n/2]=oddc;oddc+=2;}

	for(int i = 0 ; i < n ; i++ )if(!x[n/2][i]){x[n/2][i]=oddc;oddc+=2;}

	for(int i = 0 ; i < n/2 ; i++ )

	{	for(int j = 0 ; j < n/2 ; j++ )

		{

			if(oddc-2==n*n)

			{

				flag = 1 ;

				break;

			}

			x[i][j]=oddc;

			oddc+=2;

			x[n-i-1][j]=oddc;

			oddc+=2;

			x[i][n-j-1]=oddc;

			oddc+=2;

			x[n-i-1][n-j-1]=oddc;

			oddc+=2;

		}

		if(flag)

			break;

	}

	for(int a = 0 ; a < n ; a++ )

		for(int b = 0 ; b < n ; b++ )

		{

			if(!x[a][b])

			{

				cout << evenc ;

				evenc+=2;

			}

			else

				cout << x[a][b];

			if(b==n-1)

				cout << endl;

			else cout << " "; 

		}

}
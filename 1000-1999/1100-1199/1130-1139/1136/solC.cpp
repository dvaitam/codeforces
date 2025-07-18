#include <bits/stdc++.h>
#pragma comment(linker, "/stack:200000000")
#pragma GCC optimize("Ofast")
#pragma GCC target("sse,sse2,sse3,ssse3,sse4,popcnt,abm,mmx,avx,tune=native")

#include <iostream>
#include <algorithm>

#define MAX_N 505

using namespace std;

int N, M;
int arr1[MAX_N][MAX_N], arr2[MAX_N][MAX_N];
int temp1[MAX_N], temp2[MAX_N];

void fastscan(int &number) 
{ 
    //variable to indicate sign of input number 
    bool negative = false; 
    register int c; 
  
    number = 0; 
  
    // extract current character from buffer 
    c = getchar(); 
    if (c=='-') 
    { 
        // number is negative 
        negative = true; 
  
        // extract the next character from the buffer 
        c = getchar(); 
    } 
  
    // Keep on extracting characters if they are integers 
    // i.e ASCII Value lies from '0'(48) to '9' (57) 
    for (; (c>47 && c<58); c=getchar()) 
        number = number *10 + c - 48; 
  
    // if scanned input has a negative sign, negate the 
    // value of the input number 
    if (negative) 
        number *= -1; 
} 

int main() {
	fastscan(N);
	fastscan(M);
	for(int i=0; i<N; i++)
		for(int j=0; j<M; j++)
			fastscan(arr1[i][j]);
	for(int i=0; i<N; i++)
		for(int j=0; j<M; j++)
			fastscan(arr2[i][j]);
	
	// check
	for(int i=0; i<=N+M-2; i++) {
		int num=0;
		for(int j= ((i<N) ? i : N-1); j>=0 && i-j>=0 && i-j<M; j--) {
			temp1[num] = arr1[j][i-j];
			temp2[num] = arr2[j][i-j];
			num++;
		}

		sort(temp1, temp1+num);
		sort(temp2, temp2+num);

		bool flag = true;
		for(int j=0; j<num; j++)
			if(temp1[j] != temp2[j]) {
				flag = false;
				break;
			}
		if(!flag) {
			cout << "NO" << endl;
			return 0;
		}
	}

	cout << "YES" << endl;
}
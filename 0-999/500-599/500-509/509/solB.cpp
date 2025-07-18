#include <bits/stdc++.h>

using namespace std;

int num[101],mini=INT_MAX;

int main() {

	int n,k,maxi=0;

	cin >> n >> k;

	int arr[n];

	for(int i=0;i<n;i++) { 

		cin >> arr[i];

	    num[arr[i]]=1;	

	    maxi = max(maxi,arr[i]);

	    mini = min(mini,arr[i]);

    }

    if(maxi-mini>k) cout << "NO";

    else {

		cout << "YES" << endl;

		for(int i=0;i<n;i++) {

			int pop = 1, flag = 0;

			for(int j=1;j<=arr[i];j++) {

				cout << pop << " ";

				if(flag) pop++;

				if(j==mini && !flag) flag = 1;

			}

			cout << endl;

		}	

	}

	return 0;

    }
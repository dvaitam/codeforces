#include<bits/stdc++.h>
using namespace std;
#define FOR(i,s,k) for(int i=s; i<k; i++)

#define ll long long 

int main() {

int arr[5010]={0};
bool crr[5010]={0};
int n,k;
cin>>n>>k;

//int arr[n];
ll sum=0;

int temp,count=0;
bool flag=1;

int brr[n];

for(int i=0; i<n; i++) {
	cin>>temp;
	arr[temp]++;
	brr[i] = temp;

}

int max=0;

for(int i=0; i<5010; i++){

	if(arr[i]>k) {
		flag=0;
		// cout<<"arr[i] zyaada h"<<endl;
		break;
	}
	
	if(arr[i]>max) max = arr[i];

	if(arr[i]!=0) {
		count+=arr[i];
		
	}	


	if(arr[i]!=0 && arr[i]!=1){
		arr[i] = sum + arr[i];
		sum = arr[i];
	}

	if(arr[i]==1) crr[i]=1;	

	
}

if(count<k) {flag=0 ;
	// cout<<"count kam h"<<endl;
}
int temp4=1;

if(flag==1) {
	cout<<"YES"<<endl;
	
	int color = 0;
	
	for(int i=0; i<n; i++) {

		color = arr[brr[i]];
		
		color = color%k==0? k:color%k; 

		arr[brr[i]]--;

		if(crr[brr[i]]==1)	{color = (sum+1)%k==0?k:(sum+1)%k; sum++;}

		cout<<color<<" ";
	}

	cout<<endl;
}

else cout<<"NO"<<endl;


return 0;
}
#include <bits/stdc++.h>

using namespace std;

#include<stdio.h>
#include<vector>
#include<algorithm>
#include<math.h>
long long a[100001];
int main(){
	int n;
	scanf("%d",&n);
	int i,j,k;
	int m=n/2;	
	for(i=1;i<n;i+=2){
		scanf("%lld",&a[i]);
	}
	long long cur=100000000000001LL;
	for(i=n-1;i>=1;i-=2){
		bool ok=false;
		for(k=2-a[i]%2;(long long)(k)*k<=a[i];k+=2){
			if(a[i]%k) continue;
			long long miden=a[i]/k;
			if((miden-k)%2) continue;
			long long upper=(a[i]/k)+(k-1);
			if(upper + 2 < cur){
				ok=true;
				long long lower=(a[i]/k)-(k-1);
				a[i+1] = ((upper+2) + (cur-2)) * ((cur - (upper+2))/2) / 2;
				cur = lower;
				break;
			}
		}
		if(!ok){
			break;
		}
	}
	a[0] = (1 + (cur-2)) * ((cur - 1)/2) / 2;
	if(a[0]>0&&i<1){
		puts("Yes");
		for(i=0;i<n;i++){
			printf("%lld%c",a[i]," \n"[i==n-1]);
		}
	} else {
		puts("No");
	}
}
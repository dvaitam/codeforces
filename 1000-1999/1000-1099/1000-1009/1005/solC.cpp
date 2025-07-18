#include<stdio.h>
#include<algorithm>
using namespace std;
int a[120010],pow2[60];
int binarySearch(int key,int len)
{
	int left = 0;
	int right = len - 1;
	while(left <= right)
	{
		int mid = (right - left) / 2 + left;
		if (a[mid] == key)
			return mid;
		else if (a[mid] > key)
			right = mid - 1;
		else
			left = mid + 1;
	}
	return -1;
}
int main(){
	int tot=0,ans=0,n,i,j;
	pow2[1]=2;
	for(i=2;i<=30;i++){
		pow2[i]=pow2[i-1]*2;
	}
	scanf("%d",&n);
	for(i=0;i<n;i++)
		scanf("%d",&a[i]);
	sort(a,a+n);
	for(i=0;i<n;i++)
		for(j=1;j<=30;j++)
	    if(pow2[j]>a[i]){
	    	int k=binarySearch(pow2[j]-a[i],n);
	    	if(k!=-1&&(k!=i||k!=0&&a[k-1]==a[k]||k!=n-1&&a[k+1]==a[k])){
	      	ans++;break;
		  }
		}
	      
	printf("%d",n-ans);
	return 0;
	  
}
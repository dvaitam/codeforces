#include <cstdio>
#include <iostream> 

#define SIZE 2001

using namespace std;

//int b[SIZE];
int N;

long long compute(int d, int cur, int last){
	long long hi=1000000001, lo=0;
	while (lo<hi){
		long long mid=(hi+lo)>>1;
		long long v=cur+mid*d;
		if (v<=last) lo=mid+1;
		else hi=mid;
	}
	return lo;
}
int main(){
	int d;
	scanf("%d %d", &N, &d);
	int last=-1;
	long long moves=0;
	for(int i=0;i<N;i++){
		int cur;
		scanf("%d", &cur);
		//fprintf(stderr,"%d\n", cur);
		if (last>=cur){
			long long to_add=compute(d,cur,last);
			moves+=to_add;
			last=cur+d*to_add;
		}
		else
			last=cur;
	}
	cout<<moves<<endl;
	
	return 0;
}
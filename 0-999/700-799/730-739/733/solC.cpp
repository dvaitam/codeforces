#include <cstdio>

#include <vector>

#define ll long long

const int maxn=507;

int n,k;

int a[maxn],b[maxn],flag[maxn];

int nxt[maxn];

struct node{

	int pos;

	char type;

	//0 left 1 right

};//结构体后面忘了加分号

std::vector<node> ans;

void init(int left,int right){

	register int i;

	for(i=left;i<right;++i){

		nxt[i]=i+1;

	}

	nxt[right]=-1;

}

bool check(int left,int right){

	if(nxt[left]!=-1) return false;

	return true;

}

int findIndex(int left,int index){

	int cnt=0;

	register int i;

	for(i=left;i!=-1;i=nxt[i]){

		cnt++;

		if(index==i){

			return cnt;

		}

	}

	return -1;

	//error

}

int solve(int left,int right,int pre){

		ll mx=-1;

		int index,first=1;

		register int i;

		for(i=left;i!=-1;i=nxt[i]){

			if(nxt[i]!=-1&&a[i]!=a[nxt[i]]){

				if(first){

					first=0;

					mx=a[i]+a[nxt[i]];

					index=i;

				}

				else{

					ll temp=a[i]+a[nxt[i]];

					if(temp>mx){

						mx=temp;

						index=i;

					}

				}

			}

		}

		if(mx!=-1){

			int pos1=findIndex(left,index);

			int pos2=findIndex(left,nxt[index]);

			if(a[index]>a[nxt[index]]){

				ans.push_back((node){pre+pos1,'R'});

			}

			else{

				ans.push_back((node){pre+pos2,'L'});

			}

			a[index]+=a[nxt[index]];//这两句话位置不对，之前做的早了

			nxt[index]=nxt[nxt[index]];

		}

		return mx;

}

int main(){

	scanf("%d",&n);

	register int i;

	ll suma=0,sumb=0;

	for(i=0;i<n;++i){

		scanf("%d",a+i);

		suma+=a[i];

	}

	scanf("%d",&k);

	for(i=0;i<k;++i){

		scanf("%d",b+i);

		sumb+=b[i];

	}

	if(suma!=sumb) {

		printf("NO\n");

		return 0;

	}

	ll temp=0;

	int cur=0,cnt=0;

	for(i=0;i<n;++i){

		temp+=a[i];

		if(temp==b[cur]){

			flag[cnt++]=i;

			++cur;

			temp=0;

		}

		else if(temp>b[cur]){

			printf("NO\n");//NO 打成 No

			return 0;

		}

	}

	if(cur!=k) {

		printf("NO\n");

		return 0;

	}

	register int j;

	int No=0;

	for(i=0;i<cnt;++i){

		int left,right=flag[i];

		if(i==0){

			left=0;

		}

		else{

			left=flag[i-1]+1;

		}

		init(left,right);

		while(!check(left,right)){

			int p=solve(left,right,i);

			if(p==-1){

				No=1;

				break;

			}

		}

		if(No){

			break;

		}

	}

	if(No) printf("NO\n");

	else{

		printf("YES\n");

		for(i=0;i<ans.size();++i){

			node t=ans[i];

			printf("%d %c\n",t.pos,t.type);

		}

	}

	return 0;

}
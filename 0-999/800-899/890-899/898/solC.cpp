#include<cstdio>
#include<cstring>
#include<algorithm>
using namespace std;
int n;
struct ss{
	char name[20],num[20];
}node[400];
char NA[20],NU[20];
int cnt[400];
int CmpNA(ss a,ss b){
	int i=0;
	while(a.name[i]!='\0'||b.name[i]!='\0'){
		if(a.name[i]!=b.name[i])return a.name[i]>b.name[i]?1:-1;
		i++;
	}
	if(a.name[i]==b.name[i])return 0;
	return a.name[i]=='\0'?-1:1;
}
int CmpNU(ss a,ss b){
	int i=0;
	while(a.num[i]!='\0'||b.num[i]!='\0'){
		if(a.num[i]!=b.num[i])return a.num[i]>b.num[i]?1:-1;
		i++;
	}
	if(a.num[i]==b.num[i])return 0;
	return a.num[i]=='\0'?-1:1;
}
bool cmp(ss a,ss b){
	return CmpNA(a,b)==1||(CmpNA(a,b)==0&&CmpNU(a,b)==1);
}
bool check(int );
int main()
{
	int i,j,k,l,many;
	int len;
	scanf("%d",&n);k=0;
	for(i=1;i<=n;i++){
		scanf("%s",NA);
		scanf("%d",&many);
		for(j=1;j<=many;j++){
			scanf("%s",NU),len=strlen(NU),k++;
			for(l=0;l<20;l++)
				node[k].name[l]=NA[l];
			for(l=0;l<len;l++)
				node[k].num[l]=NU[len-l-1];
		}
	}
	sort(node+1,node+k+1,cmp);
	l=0;
	for(i=1;i<=k;i++){
		if(i==1||CmpNA(node[i-1],node[i])!=0)
			l++,cnt[l]++;
		else{
			if(check(i))
				cnt[l]++;
		}
	}
	printf("%d\n",l);
	l=0;
	for(i=1;i<=k;i++){
		if(i==1||CmpNA(node[i-1],node[i])!=0){
			if(i!=1)printf("\n");
			printf("%s ",node[i].name);
			l++;printf("%d ",cnt[l]);
			len=strlen(node[i].num);
			for(j=0;j<len;j++)
				printf("%c",node[i].num[len-j-1]);
			printf(" ");
		}
		else{
			if(check(i)){
				len=strlen(node[i].num);
				for(j=0;j<len;j++)
					printf("%c",node[i].num[len-j-1]);
				printf(" ");
			}
		}
	}
}
bool check(int x){
	int i=0;
	while(node[x].num[i]!='\0'){
		if(node[x].num[i]!=node[x-1].num[i])return true;
		i++;
	}
	return false;
}
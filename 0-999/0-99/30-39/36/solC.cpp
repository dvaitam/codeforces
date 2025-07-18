#include <stdio.h>

int n;
int h[3100], r[3100], R[3100];

double height_intersect(int h0, int r0, int R0, int h1, int r1, int R1) {
	double res1 = (r1-r0)/double(R0-r0)*h0;
	double res2 = (R0-r1)/double(R1-r1)*h1;
	double res3 = (R1-r0)/double(R0-r0)*h0;
	if (res2>h1)
		res2=-1000;
	else
		res2 = h0-res2;

	if (res3>h0) res3=h0;
	if (res3<0) res3=0;
	res3-=h1;
	
	if (r1>=R0)
		return h0;
	if (res1<0)
		res1=0;
	if (res2>h0)
		res2=h0;
	if (res1>res2 && res1>res3)
		return res1;
	if (res3>res2)
		return res3;
	return res2;
}

double sh[3100];
int ch[3100], cr[3100], cR[3100];
int curr;

int main() {
	FILE *fp=fopen("input.txt", "r");
	fscanf(fp, "%d", &n);
	for (int i=0; i<n; i++)
		fscanf(fp, "%d%d%d", h+i, r+i, R+i);
	fclose(fp);

	//double height=0;

	curr=1;
	sh[0]=0.0;
	ch[0]=h[0];
	cr[0]=r[0];
	cR[0]=R[0];

	for (int i=1; i<n; i++) {
		double mh=0.0;
		for (int j=0; j<curr; j++) {
			double hh = height_intersect(ch[j], cr[j], cR[j], h[i], r[i], R[i]) + sh[j];
			if (hh>mh)
				mh=hh;
		}

		int prog=0;
		for (int j=0; j<curr; j++)
			if (sh[j]+ch[j]>mh) {
				sh[prog]=sh[j];
				ch[prog]=ch[j];
				cr[prog]=cr[j];
				cR[prog]=cR[j];
				prog++;
			}

		sh[prog]=mh;
		ch[prog]=h[i];
		cr[prog]=r[i];
		cR[prog]=R[i];

		curr=prog+1;
	}

	double res=0.0;
	for (int i=0; i<curr; i++)
		if (sh[i]+ch[i]>res)
			res=sh[i]+ch[i];

	fp=fopen("output.txt", "w");
	fprintf(fp, "%.12lf\n", res);
	fclose(fp);

	return 0;
}
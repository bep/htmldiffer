# HTMLDiffer (Work In Progress)

hugo -b https://example.org/dir1 -d /mysitecompare/site1
hugo -b https://example.org/dir2 -d /mysitecompare/site2

```bash
htmldiffer --dir1 /mysitecompare/site1 --dir2 /mysitecompare/site2 --outdir /mysitecompare/result
```

Then from `/mysitecompare/result` you can run `hugo server` and open http://localhost:1313/.
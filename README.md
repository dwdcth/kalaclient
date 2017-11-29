
# kala client ( 
a pure kala client for [kala](https://github.com/ajvb/kala)
## Why making a new wheel
### A clean kala client library
The default kala library will make "jobdb.db" and "jobdb.db.lock" files.
It makes the files  in "jobdb.go", and there is no way to disable.

### Fix Delete job
In kala,the delete job url "job/id" will redirect to "job/id/".
But the default go http client will not redirect with http methods and headers.
See [this issue](https://github.com/golang/go/issues/4800).

### Result
So I replace the sling with resty and take out all the struts depended.

### Usage
The same with kala default lib.

## (Chinese) 为什么造个新轮子
### 需要一个干净的库
使用默认的client会产生"jobdb.db"、"jobdb.db.lock"，有点强迫症，不需要。
### 无法删除job
删除时会把 "job/id" 重定向到 "job/id/"。而go语言默认的http client在重定向时不会把http header和http method带上。
### 用法
和kala官方一样


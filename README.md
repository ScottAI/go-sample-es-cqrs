# 简介

这是一个使用Go语言实现的ES-CQRS的案例，通过完成一个代办事项（todo）项目进行模拟

## 使用docker运行

```
docker build -t go-sample-es-cqrs .
docker run -it --rm --name todoapp go-sample-es-cqrs
```

## 不使用docker运行

不使用docker运行,  按照下面的步骤：

```
go get https://github.com/ScottAI/go-sample-es-cqrs

cd $GOPATH/github.com/ScottAI/go-sample-es-cqrs

go build

./go-sample-es-cqrs
```
: 

T然后运行并访问 @ http://localhost:8787




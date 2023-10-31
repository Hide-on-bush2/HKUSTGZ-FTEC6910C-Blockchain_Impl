# HW2-Chain

## Setup

1. Clone the [repository](https://classroom.github.com/a/J71MnJVR) to the local.
2. Because of the firework, we must change the proxy of Golang from `direct` to `aliyun`. We can achieve this by running the following commands
```
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
```
You can check whether these commands work by running `go env` and check two variables in the list.

![](./img/1.png)

3. Go to the cloned repository and run `go mod tidy` or `go get ./...` to download all the required packages. The expected output is as follows:

![](./img/2.png)

## 
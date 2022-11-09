# ddltomodel

把 [go-zero](https://github.com/zeromicro/go-zero) 的 [goctl](https://github.com/zeromicro/go-zero/tree/master/tools/goctl) 摘出来的东西

## 安装

```shell
go install github.com/dabao-zhao/ddltomodel
```

## 使用

```
ddltomodel ddl -s=".\test.sql" -d="test"

ddltomodel datasource -u="root:123456@tcp(127.0.0.1:3306)/test" -t="*" -d="test"
```
# xlsx2proto 配置表导出工具

1. 导出xlsx为json文件和go文件
2. 数组用;和#分割，注意避开特殊字符

``` shell
# 导出为配置文件
./xlsx2proto -data=json -in=./examples/ -out=./examples/config.json

# 导出为配置go文件
./xlsx2proto -data=examples -in=./examples/ -out=./examples/config.go

```
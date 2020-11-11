# GFUZZ

*wfuzz copied to Golang*



### V1.1
- **修改请求依赖为 github.com/levigross/grequests**
- 添加参数-X/--method,可以指定请求方法
- 修复占位符在2及以上没有将其替换的bug

### V1.2

- 修复filter不计算占位符的bug
- 修复请求方法无法使用占位符的bug


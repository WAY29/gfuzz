# GFUZZ

*wfuzz copied to Golang*



### V1.1
- **修改请求依赖为 github.com/levigross/grequests**
- 添加参数-X/--method,可以指定请求方法
- 修复占位符在2及以上没有将其替换的bug

### V1.2

- 修复filter不计算占位符的bug
- 修复请求方法无法使用占位符的bug

### V1.3
- 修复POST参数错误改为GET参数的问题
- 修复请求错误依然显示请求状况的问题
- 添加参数-S/--session, 控制在请求时是否使用session
- 添加参数-f/--file, 是否将结果输出到文件

### V1.4
- 优化输出逻辑
- 添加payload:stdin, 从标准输入中获取payload
- 添加payload:dirwalk, 从某个目录递归获取文件相对路径作为payload
- 添加参数-req_delay, 控制请求超时时间
- 添加参数-conn_delay, 控制连接超时时间

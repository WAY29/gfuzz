# GFUZZ

*wfuzz copied to Golang*



### V1.6
- 修改参数
    - --json 不再自动将单引号替换成单引号,添加@filepath来从文件中读取json数据
- 修复bug
    - json参数不计算占位符




### V1.5
- 删除某些测试语句
- 添加参数--sx,--hx, 用于显示/隐藏响应中存在特定内容的过滤器
- 修复session每次调用都会使用新的session请求的问题
- 修改文件输出的参数为-o,--output
- 修改--json的参数为string, 直接使用--json {jsondata}来提交json数据,在windows下可以使用单引号代替双引号
- 添加--filter的别名-f
- 修改默认线程数为32

### V1.4
- 优化输出逻辑
- 添加payload:stdin, 从标准输入中获取payload
- 添加payload:dirwalk, 从某个目录递归获取文件相对路径作为payload
- 添加参数-req_delay, 控制请求超时时间
- 添加参数-conn_delay, 控制连接超时时间

### V1.3
- 修复POST参数错误改为GET参数的问题
- 修复请求错误依然显示请求状况的问题
- 添加参数-S/--session, 控制在请求时是否使用session
- 添加参数-f/--file, 是否将结果输出到文件

### V1.2

- 修复filter不计算占位符的bug
- 修复请求方法无法使用占位符的bug


### V1.1
- **修改请求依赖为 github.com/levigross/grequests**
- 添加参数-X/--method,可以指定请求方法
- 修复占位符在2及以上没有将其替换的bug




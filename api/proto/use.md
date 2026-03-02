1. 定义 Proto 文件
   ↓
2. 运行 make proto 编译
   ↓
3. 生成 Go 代码（Server 接口 + Client 接口 + Message 类型）
   ↓
4. 服务端实现 Server 接口
   ↓
5. 客户端使用 Client 接口调用远程服务
   ↓
6. gRPC 自动序列化/反序列化数据
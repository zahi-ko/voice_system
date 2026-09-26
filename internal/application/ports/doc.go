// Package ports 定义应用层依赖的端口接口。
//
// 存储契约：AudioStore 持久化的是 PCM 样本（canonical float32 mono，
// WAV 封装为存储实现细节）。编码格式只存在于摄取/导出边沿
// （上传 decode、下载 encode），消费样本的服务不经任何编解码。
package ports

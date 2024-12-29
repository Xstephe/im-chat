package websocket

import "net/http"

// DailOptions 定义了用于配置拨号选项的函数类型。
//
// 该类型是一个函数，接受一个 *dialOption 指针，并对其进行修改。
type DailOptions func(option *dailOption)

// dialOption 结构体保存了 WebSocket 连接的配置选项。
//
// 该结构体包括连接的 HTTP 头部和连接路径模式等设置。
type dailOption struct {
	pattern string
	header  http.Header
}

// newDialOptions 创建一个具有默认值的新的 dialOption 结构体，并根据传入的选项进行配置。
//
// 该函数初始化 dialOption 结构体的默认值，并应用提供的 DialOptions 函数。
//
// 参数:
//   - opts: 可选的 DialOptions 函数，用于配置 WebSocket 连接选项。
//
// 返回:
//   - dialOption: 配置好的 WebSocket 连接选项。
func newDailOptions(opts ...DailOptions) dailOption {
	o := dailOption{
		pattern: "/ws",
		header:  nil,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return o
}

// WithClientPatten 返回一个设置连接路径模式的 DialOptions 函数。
//
// 该函数返回一个 DialOptions 函数，用于设置 WebSocket 连接的路径模式。
//
// 参数:
//   - pattern: 连接路径模式，例如 "/ws"。
//
// 返回:
//   - DialOptions: 配置连接路径模式的函数。
func WithClientPatten(pattern string) DailOptions {
	return func(opt *dailOption) {
		opt.pattern = pattern
	}
}

// WithClientHeader 返回一个设置 HTTP 头部的 DialOptions 函数。
//
// 该函数返回一个 DialOptions 函数，用于设置 WebSocket 连接的 HTTP 头部。
//
// 参数:
//   - header: HTTP 头部设置。
//
// 返回:
//   - DialOptions: 配置 HTTP 头部的函数。
func WithClientHeader(header http.Header) DailOptions {
	return func(opt *dailOption) {
		opt.header = header
	}
}

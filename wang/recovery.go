package wang

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// 错误处理中间件
// 在server发生panic等致命错误时进行recover 捕获错误，友好的显示错误信息给用户客户端
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				mess := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(mess))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}

// print stack trace for debug
func trace(message string) string {
	var pcs [32]uintptr

	// runtime.Callers 用来返回调用栈的程序计数器, 第 0 个 Caller 是 Callers 本身，第 1 个是上一层 trace，第 2 个是再上一层的 defer func。
	// 因此，为了日志简洁一点，我们跳过了前 3 个 Caller。
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		// 获取对应的函数
		fn := runtime.FuncForPC(pc)
		// 获取到调用该函数的文件名和行号
		filename, line := fn.FileLine(pc)

		str.WriteString(fmt.Sprintf("\n\t%s: %d", filename, line))
	}
	return str.String()
}

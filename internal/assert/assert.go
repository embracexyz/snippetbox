package assert

import (
	"strings"
	"testing"
)

func Assert[T comparable](t *testing.T, got, want T) {
	// 该调用使得在测试失败，打印t.Errorf时打印的文件名和代码行号是调用Assert所在文件
	t.Helper()

	if got != want {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func StringContains(t *testing.T, got, want string) {
	t.Helper()

	if !strings.Contains(got, want) {
		t.Errorf("got = %q, want %q", got, want)
	}
}

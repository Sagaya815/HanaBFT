package messages

import (
	"crypto/md5"
	"strconv"
	"time"
)

type key int

type value []byte

func NextKeyValue(num int) (key, value) {
	h := md5.New()
	content := make([]byte, 0)
	nowTime := time.Now().UnixNano()
	content = append(content, []byte(strconv.Itoa(int(nowTime)))...)
	content = append(content, strconv.Itoa(num)...)
	h.Write(content)
	v := h.Sum(nil)
	return key(num), v
}

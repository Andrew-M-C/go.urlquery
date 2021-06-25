package urlquery

// Marshal 将 kv 解析为字节流
func Marshal(v interface{}) ([]byte, error) {
	kv, err := marshalToValues(v)
	if err != nil {
		return nil, err
	}
	s := kv.Encode()
	return []byte(s), nil
}

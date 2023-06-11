package sequence

type Redis struct {
}

func NewRedis(addr string) Sequence {
	return &Redis{}
}

func (r *Redis) Next() (seq uint64, err error) {
	// 使用 redis 实现发号器
	// incr
	return
}

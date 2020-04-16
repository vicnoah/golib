package vbool

// BoolGroup 真值组,用于判断很多次执行是否全部执行成功或全部失败
type BoolGroup struct {
	bls []bool
}

// Add 添加执行结果
func (b *BoolGroup) Add(bl bool) {
	b.bls = append(b.bls, bl)
}

// True 判断是否全true
func (b *BoolGroup) True() bool {
	if len(b.bls) == 0 {
		return false
	}
	for _, v := range b.bls {
		if v == false {
			return false
		}
	}
	return true
}

// False 判断是否全false
func (b *BoolGroup) False() bool {
	if len(b.bls) == 0 {
		return true
	}
	for _, v := range b.bls {
		if v == true {
			return false
		}
	}
	return true
}

// Len 判断数据长度
func (b *BoolGroup) Len() int {
	return len(b.bls)
}

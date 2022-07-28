package bitmap

//用于实际存储数据的字节数组
type BitMap []byte

//创建一个新的BitMap
func New() *BitMap {
	b := BitMap(make([]byte, 0))
	return &b
}

/*
由${param}传入的byte数组创建的BitMap
*/
func FromBytes(bytes []byte) *BitMap {
	bm := BitMap(bytes)
	return &bm
}

/*
获取当前BitMap的BitSize
*/
func (b *BitMap) Bitsize() int {
	return len(*b) * 8
}

/*
获取BitMap的byte数组，即容器
*/
func (b *BitMap) GetBytes() []byte {
	return *b
}

func (b *BitMap) SetValue(offset int64, val byte) {
	index := offset / 8
	bitOffset := offset % 8

	temp := byte(1 << bitOffset)

	b.grow(offset + 1) //确保容量是够的

	//val 判断是添加还是删除： 0 OR 1
	if val > 0 {
		(*b)[index] |= temp
	} else {
		(*b)[index] &^= temp //bit clear
	}
}

func (b *BitMap) GetValue(offset int64) byte {
	index := offset / 8
	bitOffset := offset % 8

	if index >= int64(len(*b)) {
		return 0
	}

	return ((*b)[index] >> byte(bitOffset)) & 0x01

}

/*
将bit的size转换为byte的size
不对外暴露
*/
func getByteSize(bitSize int64) int64 {
	if bitSize%8 == 0 {
		return bitSize / 8
	} else {
		return bitSize/8 + 1
	}
}

/*
根据给定的bitSize给BitMap扩容
*/
func (b *BitMap) grow(bitSize int64) {

	byteSize := getByteSize(bitSize)

	//获取需要的size和现有BitMap的size的差
	gap := byteSize - int64(len(*b))
	if gap <= 0 {
		return
	}
	*b = append(*b, make([]byte, gap)...)
}

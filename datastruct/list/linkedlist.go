package list

//A LinkedList struct keeps the pointer to the head and tail, which are neither dummy node.
type LinkedList struct {
	head *node
	tail *node
	size int
}

type node struct {
	prev *node
	next *node
	val  interface{}
}

//add new val into a LinkedList at the tail position defaultly.
func (list *LinkedList) Add(val interface{}) {
	if list == nil {
		panic("list does not exist")
	}
	n := &node{
		val: val,
	}
	if list.tail == nil {
		list.head = n
		list.tail = n
	} else {
		n.prev = list.tail
		list.tail.next = n
		list.tail = n
	}
	list.size++
}

//api to get the node with given index, index: [0,size)
func (list *LinkedList) findByIndex(index int) *node {
	var n *node

	if index < list.size/2 {
		n = list.head
		for i := 0; i < index; i++ {
			n = n.next
		}
	} else {
		n = list.tail
		for i := list.size - 1; i > index; i-- {
			n = n.prev
		}
	}
	return n
}

//get the val with the given index, index:[0, size)
func (list *LinkedList) Get(index int) interface{} {
	if list == nil {
		panic("list does not exist")
	}
	if index < 0 || index >= list.size {
		panic("index out of bound")
	}

	return list.findByIndex(index).val
}

//set the val of the node with the given index as the new val, index: [0, size]
func (list *LinkedList) Set(index int, val interface{}) {
	if list == nil {
		panic("list does not exist")
	}
	if index < 0 || index > list.size {
		panic("index out of bound")
	}
	//if index == list.size, a new node should be insert into the list
	if index == list.size {
		list.Add(val)
		return
	}

	n := list.findByIndex(index)
	n.val = val
}

//insert a new node into the list at the position before the node of the given index, index: [0, size]
func (list *LinkedList) Insert(index int, val interface{}) {
	if list == nil {
		panic("list does not exist")
	}
	if index < 0 || index > list.size {
		panic("index out of bound")
	}
	if index == list.size {
		list.Add(val)
		return
	}

	pivot := list.findByIndex(index)

	n := &node{
		val:  val,
		prev: pivot.prev,
		next: pivot,
	}
	//pivot is the head
	if pivot.prev == nil {
		list.head = n
	} else {
		n.prev.next = n
	}
	pivot.prev = n
	list.size++

}

//Remove node
func (list *LinkedList) removeNode(n *node) {
	if n.prev == nil {
		list.head = n.next
	} else {
		n.prev.next = n.next
	}

	if n.next == nil {
		list.tail = n.prev
	} else {
		n.next.prev = n.prev
	}

	n.prev = nil
	n.next = nil

	list.size--
}

//Remove node by the given node, and return the val of this node
func (list *LinkedList) Remove(index int) (val interface{}) {
	if list == nil {
		panic("list does not exist")
	}

	if index < 0 || index >= list.size {
		panic("index out of bound")
	}

	n := list.findByIndex(index)

	list.removeNode(n)

	return n.val
}

func (list *LinkedList) RemoveLast() (val interface{}) {
	if list == nil {
		panic("list does not exist")
	}

	n := list.tail

	list.removeNode(n)

	return n.val
}

func (list *LinkedList) Len() int {
	if list == nil {
		panic("list does not exist")
	}
	return list.size
}

func (list *LinkedList) ForEach(consumer func(int, interface{}) bool) {
	if list == nil {
		panic("list does not exist")
	}

	n := list.head
	i := 0

	for n != nil {
		goNext := consumer(i, n.val)

		if !goNext {
			break
		}

		i++
		n = n.next
	}
}

func Make(vals ...interface{}) *LinkedList {
	list := LinkedList{}

	for _, v := range vals {
		list.Add(v)
	}
	return &list
}

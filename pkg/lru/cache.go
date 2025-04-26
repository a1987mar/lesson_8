package lru

type LruCache interface {
	// Якщо наш кеш вже повний (ми досягли нашого capacity)
	// то має видалитись той елемент, який ми до якого ми доступались (читали) найдавніше
	Put(key, value string)
	Get(key string) (string, bool)
}

type node struct {
	key, value string
	prev, next *node
}

type casheLru struct {
	capac int
	items map[string]*node
	head  *node
	tail  *node
}

func NewLruCache(capacity int) LruCache {
	return &casheLru{
		capac: capacity,
		items: make(map[string]*node),
	}
}

func (l *casheLru) Put(key, value string) {
	if n, ok := l.items[key]; ok {
		n.value = value
		l.moveFront(n)
		return
	}
	newNode := &node{key: key, value: value}
	l.items[key] = newNode
	l.add(newNode)
	if len(l.items) > l.capac {
		l.removeTail()
	}
}

func (l *casheLru) Get(key string) (string, bool) {
	if n, ok := l.items[key]; ok {
		l.moveFront(n)
		return n.value, true
	}
	return "", false
}

func (l *casheLru) add(n *node) {
	n.next = l.head
	if l.head != nil {
		l.head.prev = n
	}
	l.head = n
	if l.tail == nil {
		l.tail = n
	}
}

func (l *casheLru) removeNode(n *node) {
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		l.head = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	} else {
		l.tail = n.prev
	}
	n.prev = nil
	n.next = nil
}

func (l *casheLru) removeTail() {
	if l.tail == nil {
		return
	}
	delete(l.items, l.tail.key)
	l.removeNode(l.tail)
}

func (l *casheLru) moveFront(n *node) {
	if l.head == n {
		return
	}
	l.removeNode(n)
	l.add(n)
}

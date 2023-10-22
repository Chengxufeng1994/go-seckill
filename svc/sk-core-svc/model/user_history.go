package model

import "sync"

// ProductCountMgr 產品數量管理
type ProductCountMgr struct {
	productCount map[int]int
	lock         sync.RWMutex
}

func NewProductCountMgr() *ProductCountMgr {
	productMgr := &ProductCountMgr{
		productCount: make(map[int]int, 128),
	}
	return productMgr
}

// 商品数量
func (p *ProductCountMgr) Count(productId int) (count int) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	count = p.productCount[productId]
	return
}

// 新增商品
func (p *ProductCountMgr) Add(productId, count int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	cur, ok := p.productCount[productId]
	if !ok {
		cur = count
	} else {
		cur += count
	}
	p.productCount[productId] = cur
}

type UserBuyHistory struct {
	History map[int]int
	Lock    sync.RWMutex
}

func (p *UserBuyHistory) GetProductBuyCount(productId int) int {
	p.Lock.RLock()
	defer p.Lock.RUnlock()

	count, _ := p.History[productId]
	return count
}

func (p *UserBuyHistory) Add(productId, count int) {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	cur, ok := p.History[productId]
	if !ok {
		cur = count
	} else {
		cur += count
	}

	p.History[productId] = cur
}

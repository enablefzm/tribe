package gameroom

import (
	"errors"
)

// 缓存管理器
//	所有被放在内存的管理对象都要在这里被处理
//	1、item 管理

// 缓存管理器里的对象接口
type IFCache interface {
	Name() string // 缓存管理器名称
	Max() int     // 获取当前缓存里的总数
	Count() int   // 获取当前缓存数量总数
	Release()     // 释放当前对象空间
	Save() error  // 保存缓存里的对象
}

// 缓存对管理器
type manageCache struct {
	arrCache []IFCache
}

// 获取指定的Cache对象
func (this *manageCache) GetCache(name string) (IFCache, error) {
	for _, obCache := range this.arrCache {
		if obCache.Name() == name {
			return obCache, nil
		}
	}
	return nil, errors.New("没有发现指定的Cache")
}

// 放入Cache
func (this *manageCache) PutCache(obCache IFCache) error {
	_, err := this.GetCache(obCache.Name())
	if err == nil {
		return errors.New("已存在名为" + obCache.Name() + "对象，操作不成功")
	}
	this.arrCache = append(this.arrCache, obCache)
	return nil
}

// 获取Cache列表
func (this *manageCache) GetList() []string {
	result := make([]string, 0, len(this.arrCache))
	for _, tCache := range this.arrCache {
		result = append(result, tCache.Name())
	}
	return result
}

// 保存所有管理缓存对象
func (this *manageCache) Save() error {
	for _, obCache := range this.arrCache {
		err := obCache.Save()
		if err != nil {
			return errors.New("保存所有缓对象出错：" + err.Error())
		}
	}
	return nil
}

// 清空所有对象的不必要的数据释放内存暂用
func (this *manageCache) Release() {
	for _, obCache := range this.arrCache {
		// 启用线程
		go obCache.Release()
	}
}

var OBManageCache = &manageCache{
	arrCache: make([]IFCache, 0, 5),
}

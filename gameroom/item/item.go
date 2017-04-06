package item

import (
	"errors"
	"fmt"
	"strings"
	"tribe/baseob"
	"vava6/vatools"
)

/**
 物品类别
	101	普通物品没有具体作用
	201	材料-草药
	202	材料-矿石
	203	材料-皮革
	204 材料-木才
	205 材料-布
	301 弓
	302 矛
	303 盾
	304 刀
	305 法杖
	401 头部
	402 胸部
	403	腿部
	404 脚部（鞋子）
	405	饰品
	501 食物-肉
	502 食物-水果
	503 食物-水
	504 食物-植物
 物品基本属性可以被共享的
	eNeed			int // 食用需要多少饥饿值
	ePow			int // 食物增加力量值
	eStamina		int // 食物增加体能值
	eAqile			int // 食物增加敏捷值
	eIq				int // 食物增加智力值
	reMinPow		int // ？随机增长 力量 数值下限
	reMaxPow		int // ？随机增长 力量 数值上限
	reMinStamina    int // ？随机增长 体能 数值下限
	reMaxStamina	int // ？随机增长 体能 数值上限
	reMinAqile		int // ？随机增长 敏捷 数值下限
	reMaxAqile		int // ？随机增长 敏捷 数值上限
	reMinIq			int // ？随机增长 智力 数值下限
	reMaxIq			int // ？随机增长 智力 数值上限
	nedLv 		    int // 使用需要的级别
	nedHType		int // 使用英雄的要求
 物品属性
	minAtt 			int	// 最小攻击
	maxAtt 			int	// 最大攻击
	speed           int // 攻击速度
	pow  			int // 力量
	stamina	 		int	// 体能
	agile   		int	// 敏捷
	iq 				int // 智力
	hit	   			int	// 命中
	crit			int	// 暴击
	def				int // 防御力
	wareType		int // 装备类型：1布甲、2皮甲、3重甲
	addAtt			int // 附加增加攻击力
	addPow			int // 附加增加力量
	addHit			int // 附加增加命中值
	addCrit			int // 附加增加暴击值
	addStamina 		int // 附加增加体能值
	addAgile		int // 附加增加敏捷值
	addIq 			int // 附加增加智力值
	addDef		    int // 附加增加防御值
	func			int // 可以使用的功能
 随机生成参数
	rAtt			int	// 生成随机攻击力 在原有att基础上增加的值
	rPow            int // 生成随机力量
	rHit			int // 生成随机命中力 在原有hit基础上增加的值
	rCrit			int // 生成随机暴击值 ---
	rStamina		int // 生成随机体能值
	rAgile			int // 生成随机敏捷值
	rIq 			int // 生成随机智力值
	rDef            int // 生成随机防御值
	roAtt 	     	int // 随机生成攻击力选项之一
	roPow			int // 随机生成力量选项之一
	roHit			int // 随机生成命中值选项之一
	roCrit			int // 随机生成暴击值选项之一
	roStamina  		int // 随机生成体能值
	roAgile    		int // 随机生成敏捷值
	roIq 			int // 随机生成智力选项之一
	roDef           int // 随机生成防御力
	roNum   		int // 随机生成数
*/
type Item struct {
	baseob.BaseOB
	itemID           int            // 物品基类ID
	itemName         string         // 名称
	itemImg          string         // 图片
	itemType         uint16         // 商品类型
	itemType2        uint16         // 物品二级类
	itemQuality      uint8          // 商品质量
	itemMsg          string         // 商品祥细信息
	itemValue        uint16         // 价值
	itemLevel        uint16         // 物品等级
	superposition    bool           // 是否可以叠加
	itemProperty     map[string]int // 属性
	itemRndProperty  map[string]int // 生成物品的随机参数
	itemBaseProperty map[string]int // 物品基本的其它属性(可被共享的属性)
}

// 从数据库里加载Item
func NewItem(itemID int) (*Item, error) {
	ob := &Item{}
	rs, err := ob.LoadDB("d_item", "*", "itemID", map[string]interface{}{"itemID": itemID})
	if err != nil {
		return nil, errors.New("没有数据存在")
	}
	ob.itemID = itemID
	ob.itemImg = rs["itemImg"]
	ob.itemName = rs["itemName"]
	ob.itemType = vatools.SUint16(rs["itemType"])
	ob.itemQuality = vatools.SUint8(rs["itemQuality"])
	ob.itemMsg = rs["itemMsg"]
	ob.itemValue = vatools.SUint16(rs["itemValue"])
	ob.itemLevel = vatools.SUint16(rs["itemLevel"])
	if rs["superposition"] == "1" {
		ob.superposition = true
	} else {
		ob.superposition = false
	}
	// 解析属性
	ob.initProperty(rs["itemProperty"])
	// 解析随机生成属性
	ob.initCreateProperty(rs["itemRndProperty"])
	// 获取物品基本属性
	ob.itemBaseProperty = ob.initStrProperty(rs["itemBaseProperty"])
	return ob, nil
}

func (this *Item) ItemID() int {
	return this.itemID
}

func (this *Item) ItemName() string {
	return this.itemName
}

func (this *Item) initProperty(strVal string) {
	this.itemProperty = this.initStrProperty(strVal)
}

func (this *Item) initCreateProperty(strVal string) {
	this.itemRndProperty = this.initStrProperty(strVal)
}

func (this *Item) initStrProperty(strVal string) map[string]int {
	arr := strings.Split(strVal, ",")
	il := len(arr)
	result := make(map[string]int, il)
	for i := 0; i < il; i++ {
		arrField := strings.Split(arr[i], "=")
		if len(arrField) == 2 {
			// 如果是攻击属性配置
			if arrField[0] == "att" {
				// 获取最高攻击和最低攻击值
				arrAtt := strings.Split(arrField[1], "-")
				if len(arrAtt) == 2 {
					result["minAtt"] = vatools.SInt(arrAtt[0])
					result["maxAtt"] = vatools.SInt(arrAtt[1])
					if result["maxAtt"] < result["minAtt"] {
						result["maxAtt"] = result["minAtt"] + 1
					}
				}
			} else {
				result[arrField[0]] = vatools.SInt(arrField[1])
			}
		}
	}
	return result
}

// 获取物品的基本属性(可被共享的属性)
func (this *Item) GetBaseProperty() map[string]int {
	result := make(map[string]int, len(this.itemBaseProperty))
	for k, v := range this.itemBaseProperty {
		result[k] = v
	}
	return result
}

func (this *Item) GetRndProperty() map[string]int {
	result := make(map[string]int)
	rndOnley := make(map[string]int)
	// 判断是否有随机挑选的物品信息
	roNum := 0
	for k, v := range this.itemRndProperty {
		switch k {
		case "roAtt", "roHit", "roCrit", "roStamina", "roAgile", "roIq", "roDef":
			rndOnley[k] = v
		case "roNum":
			roNum = v
		default:
			result[k] = this.getRndPropertyValue(v)
		}
	}
	// 生成随机生成其中几项
	if roNum > 0 {
		for i := 0; i < roNum; i++ {
			il := len(rndOnley)
			if il < 1 {
				continue
			}
			j := 1
			l := vatools.CRnd(1, il)
			for k, v := range rndOnley {
				if j != l {
					j++
					continue
				}
				result[k] = this.getRndPropertyValue(v)
				delete(rndOnley, k)
				break
			}
		}
	}
	// 返回结果值
	return result
}

// 随机获得的物品数值
func (this *Item) getRndPropertyValue(v int) int {
	if v > 0 {
		// 多少机率出现好的物品
		rndVal := vatools.CRnd(0, 1000)
		tmpMaxVal := v
		switch {
		case rndVal < 400:
			tmpMaxVal = v / 3
		case rndVal < 800:
			tmpMaxVal = v / 2
		}
		if tmpMaxVal > 0 {
			tmpMaxVal = vatools.CRnd(0, tmpMaxVal)
		}
		return tmpMaxVal
	} else {
		return 0
	}
}

func (this *Item) GetInfo() string {
	return fmt.Sprint("名称：", this.itemName, " 描述：", this.itemMsg)
}

func (this *Item) GetFieldInfo() map[string]interface{} {
	res := map[string]interface{}{
		"itemID":        this.itemID,
		"itemName":      this.itemName,
		"itemImg":       this.itemImg,
		"itemType":      this.itemType,
		"itemMsg":       this.itemMsg,
		"itemValue":     this.itemValue,
		"itemLevel":     this.itemValue,
		"superposition": this.superposition,
		"itemProperty":  this.itemProperty,
	}
	return res
}

func (this *Item) Superposition() bool {
	return this.superposition
}

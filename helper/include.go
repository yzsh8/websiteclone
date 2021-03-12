package helper

import (
	"reflect"
	"strconv"
)

/**
* 这里定义一些常用的断函数
 */

/**
* 利用反射机制，对接口数据进行取值
 */
func UnknowToString(unknow interface{}) string {
	result := ""

	if unknow == nil {
		return result
	}

	switch unknow.(type) {
	case string:
		result, _ = unknow.(string)
	case uint:
		v1, _ := unknow.(uint)
		result = strconv.Itoa(int(v1))
	case uint8:
		v1, _ := unknow.(uint8)
		result = strconv.Itoa(int(v1))
	case uint16:
		v1, _ := unknow.(uint16)
		result = strconv.Itoa(int(v1))
	case uint32:
		v1, _ := unknow.(uint32)
		result = strconv.Itoa(int(v1))
	case uint64:
		v1, _ := unknow.(uint64)
		result = strconv.Itoa(int(v1))
	case int:
		v1, _ := unknow.(int)
		result = strconv.Itoa(v1)
	case int8:
		v1, _ := unknow.(int8)
		result = strconv.Itoa(int(v1))
	case int16:
		v1, _ := unknow.(int16)
		result = strconv.Itoa(int(v1))
	case int32:
		v1, _ := unknow.(int32)
		result = strconv.Itoa(int(v1))
	case int64:
		v1, _ := unknow.(int64)
		result = strconv.Itoa(int(v1))
	case float32:
		v1, _ := unknow.(float32)
		result = strconv.FormatFloat(float64(v1), 'f', -1, 64)
	case float64:
		v1, _ := unknow.(float64)
		result = strconv.FormatFloat(v1, 'f', -1, 64)
	case bool:
		v1, _ := unknow.(bool)
		result = strconv.FormatBool(v1)
	case []byte:
		v1, _ := unknow.([]byte)
		result = string(v1)
	default:
		result = ""
	}
	return result
}

/**
* 计算数组切片中的最大值
 */
func MaxFunc(l []float64) (max float64) {
	max = l[0]
	for _, v := range l {
		if v > max {
			max = v
		}
	}
	return max
}

/**
* 计算数组切片中的最小值
 */
func MinFunc(l []float64) (min float64) {
	min = l[0]
	for _, v := range l {
		if v < min {
			min = v
		}
	}
	return min
}

/**
* 两个字段相差不对的结构体，直接赋值
* binding  		要修改的结构体
* value 		有数据的结构体
 */
func StructAssign(binding interface{}, value interface{}) {
	bVal := reflect.ValueOf(binding).Elem() //获取reflect.Type类型
	vVal := reflect.ValueOf(value).Elem()   //获取reflect.Type类型
	vTypeOfT := vVal.Type()
	for i := 0; i < vVal.NumField(); i++ {
		// 在要修改的结构体中查询有数据结构体中相同属性的字段，有则修改其值
		name := vTypeOfT.Field(i).Name
		if ok := bVal.FieldByName(name).IsValid(); ok {
			v1 := vVal.Field(i).IsValid()
			if v1 == false {
				continue
			}
			if vVal.Field(i).CanInterface() == true {
				vv := vVal.Field(i).Interface()
				v2 := reflect.ValueOf(vv)
				bVal.FieldByName(name).Set(v2)
			}
		}
	}
}

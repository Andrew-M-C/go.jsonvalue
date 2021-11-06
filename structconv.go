package jsonvalue

import (
	"encoding/json"
)

// Export convert jsonvalue to another type of parameter. The target parameter type should match the type of *V.
//
// Export 将 *jsonvalue.V 转到符合原生 encoding/json 的一个 struct 中。该函数只是个便利的函数封装，途中需要一次序列化和反序列化，
// 性能不是最优。
func (v *V) Export(dst interface{}) error {
	b, err := v.Marshal()
	if err != nil {
		return err
	}

	return json.Unmarshal(b, dst)
}

// Import convert json value from a marsalable parameter to *V.
//
// Import 将符合 encoding/json 的 struct 转为 *jsonvalue.V 类型。该函数只是个便利的函数封装，途中需要一次序列化和反序列化，
// 性能不是最优。后续有计划改为无需序列化直接转换，并且支持浮点数的 Inf, -Inf, NaN 的特殊处理，敬请期待。
func Import(src interface{}) (*V, error) {
	b, err := json.Marshal(src)
	if err != nil {
		return &V{}, err
	}
	return Unmarshal(b)
}

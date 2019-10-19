// Package jsonvalue is for JSON parsing and setting. It is used in situations those Go structures cannot achieve, or "map[string]interface{}" could not do properbally.
//
// As a quick start:
// 	v := jsonvalue.NewObject()
// 	v.SetString("Hello, JSON").At("someObject", "someObject", "someObject", "message")  // automatically create sub objects
// 	fmt.Println(v.MustMarshalString())                                                  // marshal to string type
// 	// Output:
// 	// {"someObject":{"someObject":{"someObject":{"message":"Hello, JSON!"}}}
// If you want to parse raw JSON data, use Unmarshal()
// 	raw := []byte(`{"message":"hello, world"}`)
// 	v, err := jsonvalue.Unmarshal(raw)
// 	s, _ := v.GetString("message")
// 	fmt.Println(s)
// 	// Output:
// 	// hello, world
package jsonvalue

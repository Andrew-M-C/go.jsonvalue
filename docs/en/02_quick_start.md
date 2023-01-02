<font size=6>Quick Start</font>

[Prev Page](./01_introduction.md) | [Contents](./README.md) | [Next Page](./03_set.md)

---

To create a complex JSON object as follow:

```json
{
	"obj": {
		"obj": {
			"obj": {
				"str": "Hello, JSON!"
			}
		}
	}
}
```

Just three lines with jsonvalue:

```go
	v := jsonvalue.NewObject()
	v.Set("Hello, JSON").At("obj", "obj", "obj", "str")
	fmt.Println(v.MustMarshalString())
```

Output: `{"obj":{"obj":{"obj":{"str":"Hello, JSON!"}}}`

On the other hand, if we want to fetch one single data from the JSON bytes above, it is also easy:

```go
const raw = `{"obj": {"obj": {"obj": {"str": "Hello, JSON!"}}}}`
s := jsonvalue.MustUnmarshalString(s).MustGet("obj", "obj", "obj", "str").String()
fmt.Println(s)
```

Output: `Hello, JSON!`



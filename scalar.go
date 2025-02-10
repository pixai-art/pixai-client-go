package pixai_client

type JSONObject map[string]any

func ToJSONObject(o any) JSONObject {
	switch s := o.(type) {
	case JSONObject:
		return s
	case map[string]any:
		return JSONObject(s)
	}
	return nil
}

func (o *JSONObject) GetGraphQLType() string { return "JSONObject" }

func (o *JSONObject) GetString(key string) string {
	if o == nil {
		return ""
	}
	if v, ok := (*o)[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func (o *JSONObject) GetObject(key string) JSONObject {
	if o == nil {
		return nil
	}
	if v, ok := (*o)[key]; ok {
		return ToJSONObject(v)
	}
	return nil
}

func (o *JSONObject) GetArray(key string) []any {
	if o == nil {
		return nil
	}
	if v, ok := (*o)[key]; ok {
		if s, ok := v.([]any); ok {
			return s
		}
	}
	return nil
}

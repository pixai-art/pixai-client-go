package pixai_client

type JSONObject map[string]any

func (o *JSONObject) GetGraphQLType() string { return "JSONObject" }

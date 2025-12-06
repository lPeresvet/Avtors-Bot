package model

type GisKznZoneIDResp struct {
	Type          string `json:"type"`
	TotalFeatures string `json:"totalFeatures"`
	Features      []struct {
		Type       string      `json:"type"`
		Id         string      `json:"id"`
		Geometry   interface{} `json:"geometry"`
		Properties struct {
			Key int64 `json:"key"`
		} `json:"properties"`
	} `json:"features"`
	Crs interface{} `json:"crs"`
}

type FuncZoneDetailsInfo struct {
	Key          int64  `json:"key"`
	Alias        string `json:"alias"`
	Project      string `json:"project"`
	Title        string `json:"title"`
	GroupsFields []struct {
		Name   *string `json:"name"`
		Fields []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"fields"`
	} `json:"groupsFields"`
	Images    []interface{} `json:"images"`
	FilesInfo []interface{} `json:"filesInfo"`
	Comments  []interface{} `json:"comments"`
}

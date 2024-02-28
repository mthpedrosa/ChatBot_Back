package models

type Trigger struct {
	Audio   string         `json:"audio" bson:"audio"` // Node for format audio
	Message TriggerOptions `json:"message" bson:"message"`
}

type TriggerOptions struct {
	Options []Option `json:"options" bson:"options"`
	Default string   `json:"default" bson:"default"`
}

type Option struct {
	Content  string `json:"content" bson:"content"`
	NextNode string `json:"next_node" bson:"next_node"`
}

type Button struct {
	ID       string `json:"id" bson:"id"`
	Title    string `json:"title" bson:"title"`
	NextNode string `json:"next_node" bson:"next_node"`
}

type Node struct {
	ID         string `json:"id" bson:"id"`
	Name       string `json:"name" bson:"name"`
	Type       string `json:"type" bson:"type"`
	Parameters struct {
		Buttons       []Button `json:"buttons" bson:"buttons"`
		Triggers      Trigger  `json:"triggers" bson:"triggers"`
		ErrorOutput   string   `json:"error_output" bson:"error_output"`
		SuccessOutput string   `json:"success_output" bson:"success_output"`
		Command       string   `json:"command" bson:"command"`
		Content       string   `josn:"content" bson:"content"`
	} `json:"parameters"`
	Position []int `json:"position" bson:"position"` //used in the fron
}

/*type Campo struct {
	Entidade  string `json:"entidade"`
	Tipo      string `json:"tipo"`
	CampoNome string `json:"campo_name" bson:"campo_name"`
	Conteudo  string `json:"conteudo"`
}*/

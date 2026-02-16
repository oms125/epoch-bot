package player

import (
	"fmt"
	"strings"
	"slices"
)

type (
	Armory []*Tool

	Tool struct {
		*ToolData
		*ToolMetadata
	}
	ToolData struct {
		ID int
		Name string
		MaxDurability int
		Slot int
	}
	ToolMetadata struct {
		Durability int `db:"durability"`
		Equipped int `db:"equipped"`
	}
)

var ArmoryFields = []Field{
	{
		Name: "item_id",
		Type: "INTEGER",
	},
	{
		Name: "durability",
		Type: "INTEGER",
	},
	{
		Name: "equipped",
		Type: "INTEGER",
	},
}

//Game Logic
func (t *ToolData) Type() int { return TYPE_TOOL }
func (t *ToolMetadata) Type() int { return TYPE_TOOL }
func (t *Tool) Type() int { return TYPE_TOOL }

func NewDefaultTool(id int) *Tool {
	data := ITEM_DATA[id].(*ToolData)
	return &Tool{
		ToolData: data,
		ToolMetadata: &ToolMetadata{
			Durability: data.MaxDurability,
		},
	}
}

func (p *Player) RemoveTool(idx int) {
	p.Arm = slices.Delete(p.Arm, idx, idx+1)
}

func (p *Player) AddTool(tool *Tool) error {
	if len(p.Arm) >= p.ArmSize {
		return &ArmFullError{}
	}
	p.Arm = append(p.Arm, tool)
	return nil
}

//Helpers
func (arm Armory) String() string {
	var armString strings.Builder; 
	for _, tool := range arm {
		fmt.Fprintf(&armString, "%s: (%d/%d) %d\n", tool.Name, tool.Durability, tool.MaxDurability, tool.Equipped)
	}
	return armString.String()
}
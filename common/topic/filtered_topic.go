package topic

import (
	"fmt"

	"github.com/blocklords/sds/common/data_type/key_value"
)

type TopicFilter struct {
	Organizations  []string `json:"o,omitempty"`
	Projects       []string `json:"p,omitempty"`
	NetworkIds     []string `json:"n,omitempty"`
	Groups         []string `json:"g,omitempty"`
	Smartcontracts []string `json:"s,omitempty"`
	Events         []string `json:"e,omitempty"`
}

func NewFilterTopic(o []string, p []string, n []string, g []string, s []string, e []string) TopicFilter {
	return TopicFilter{
		Organizations:  o,
		Projects:       p,
		NetworkIds:     n,
		Groups:         g,
		Smartcontracts: s,
		Events:         e,
	}
}

func (t *TopicFilter) Len(level uint8) int {
	switch level {
	case ORGANIZATION_LEVEL:
		return len(t.Organizations)
	case PROJECT_LEVEL:
		return len(t.Projects)
	case NETWORK_ID_LEVEL:
		return len(t.NetworkIds)
	case GROUP_LEVEL:
		return len(t.Groups)
	case SMARTCONTRACT_LEVEL:
		return len(t.Smartcontracts)
	case FULL_LEVEL:
		return len(t.Events)
	default:
		return len(t.Organizations) + len(t.Projects) + len(t.NetworkIds) + len(t.Groups) + len(t.Smartcontracts) + len(t.Events)
	}
}

// topic key
func (t *TopicFilter) Key() TopicKey {
	return TopicKey(t.ToString())
}

// list of path
func list(properties []string) string {
	str := ""
	for _, v := range properties {
		str += "," + v
	}

	return str
}

// Convert the topic filter object to the topic filter string.
func (t *TopicFilter) ToString() string {
	str := ""
	if len(t.Organizations) > 0 {
		str += "o:" + list(t.Organizations) + ";"
	}
	if len(t.Projects) > 0 {
		str += "p:" + list(t.Projects) + ";"
	}
	if len(t.NetworkIds) > 0 {
		str += "n:" + list(t.NetworkIds) + ";"
	}
	if len(t.Groups) > 0 {
		str += "g:" + list(t.Groups) + ";"
	}
	if len(t.Smartcontracts) > 0 {
		str += "s:" + list(t.Smartcontracts) + ";"
	}
	if len(t.Events) > 0 {
		str += "e:" + list(t.Events) + ";"
	}

	return str
}

func NewFromKeyValueParameter(parameters key_value.KeyValue) (*TopicFilter, error) {
	topic_filter_map, err := parameters.GetKeyValue("topic_filter")
	if err != nil {
		return nil, fmt.Errorf("missing `topic_filter` parameter")
	}

	return NewFromKeyValue(topic_filter_map), nil
}

// Converts the JSON object to the topic.TopicFilter
func NewFromKeyValue(parameters key_value.KeyValue) *TopicFilter {
	topic_filter := TopicFilter{
		Organizations:  []string{},
		Projects:       []string{},
		NetworkIds:     []string{},
		Groups:         []string{},
		Smartcontracts: []string{},
		Events:         []string{},
	}

	organizations, err := parameters.GetStringList("o")
	if err == nil {
		topic_filter.Organizations = organizations
	}
	projects, err := parameters.GetStringList("p")
	if err == nil {
		topic_filter.Projects = projects
	}
	network_ids, err := parameters.GetStringList("n")
	if err == nil {
		topic_filter.NetworkIds = network_ids
	}
	groups, err := parameters.GetStringList("g")
	if err == nil {
		topic_filter.Groups = groups
	}
	smartcontracts, err := parameters.GetStringList("s")
	if err == nil {
		topic_filter.Smartcontracts = smartcontracts
	}
	logs, err := parameters.GetStringList("e")
	if err == nil {
		topic_filter.Events = logs
	}

	return &topic_filter
}

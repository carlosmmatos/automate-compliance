package opencontrol

type Family string

type OpenControlEntry struct {
	ControlKey           string           `json:"control_key"`
	CoveredBy            []string         `json:"covered_by"`
	ImplementationStatus string           `json:"implementation_status"`
	Narrative            []NarrativeEntry `json:"narrative"`
}

type NarrativeEntry struct {
	Key  string `json:"key,omitempty"`
	Text string `json:"text"`
}

func ControlEntryWithDefaults() OpenControlEntry {
	return OpenControlEntry{
		CoveredBy: []string{},
		Narrative: []NarrativeEntry{},
	}
}

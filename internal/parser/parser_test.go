package parser

import (
	"reflect"
	"testing"

	v3c "github.com/opencontrol/compliance-masonry/pkg/lib/components/versions/3_1_0"
)

func buildControlEntryWithDefaults(key string, n v3c.NarrativeSection) v3c.Satisfies {
	oce := v3c.Satisfies{}
	oce.ControlKey = key
	oce.Narrative = append(oce.Narrative, n)
	return oce
}

func TestParser_normalizeFamily(t *testing.T) {
	type args struct {
		family string
	}

	p := NewParser()

	tests := []struct {
		name string
		args args
		want controlFamily
	}{
		{
			"simple family", args{"ACCESS CONTROL"}, controlFamily("ACCESS_CONTROL"),
		},
		{
			"family with extra spaces", args{"ACCESS   CONTROL"}, controlFamily("ACCESS_CONTROL"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := p.normalizeFamily(tt.args.family); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.normalizeFamily() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseControl(t *testing.T) {
	type args struct {
		control string
	}

	p := NewParser()

	tests := []struct {
		name    string
		args    args
		want    v3c.Satisfies
		wantErr bool
	}{
		{
			"Simple control is successfully parsed",
			args{"AC-1"},
			buildControlEntryWithDefaults("AC-1", v3c.NarrativeSection{
				Text: "Text only",
			}),
			false,
		},
		{
			"Control with enhancement is successfully parsed",
			args{"AC-2a."},
			buildControlEntryWithDefaults("AC-2", v3c.NarrativeSection{
				Key:  "a",
				Text: "Text for enhancement",
			}),
			false,
		},
		{
			"Control with enhancement and sub-enhancement is successfully parsed",
			args{"AC-2a.1."},
			buildControlEntryWithDefaults("AC-2", v3c.NarrativeSection{
				Key:  "a.1",
				Text: "Text for enhancement",
			}),
			false,
		},
		{
			"Simple Sub-control is successfully parsed",
			args{"AC-2 (1)"},
			buildControlEntryWithDefaults("AC-2 (1)", v3c.NarrativeSection{
				Text: "Text only",
			}),
			false,
		},
		{
			"Simple Sub-control with high number is successfully parsed",
			args{"SC-42 (3)"},
			buildControlEntryWithDefaults("SC-42 (3)", v3c.NarrativeSection{
				Text: "Text only",
			}),
			false,
		},
		{
			// At some point there was a bug in which AC-2 (10) wasn't
			// parsed properly. This ensures that such values are parsed.
			"Simple Sub-control (above 9) is successfully parsed",
			args{"AC-2 (21)"},
			buildControlEntryWithDefaults("AC-2 (21)", v3c.NarrativeSection{
				Text: "Text only",
			}),
			false,
		},
		{
			"Sub-control with enhancement is successfully parsed",
			args{"AC-3 (3)(a)"},
			buildControlEntryWithDefaults("AC-3 (3)", v3c.NarrativeSection{
				Key:  "a",
				Text: "Text for enhancement",
			}),
			false,
		},
		{
			"Sub-control with enhancement and sub-enhancement is successfully parsed",
			args{"AC-3 (3)(b)(1)"},
			buildControlEntryWithDefaults("AC-3 (3)", v3c.NarrativeSection{
				// TODO(jaosorior): What's the appropriate value here?
				//                  I couldn't find examples... and so I
				//                  left it at  just the enhancement number.
				Key:  "b",
				Text: "Text for enhancement",
			}),
			false,
		},
		{
			"Sub-control with enhancement and extra spaces is successfully parsed",
			args{"AC-3   (3)(a)"},
			buildControlEntryWithDefaults("AC-3 (3)", v3c.NarrativeSection{
				Key:  "a",
				Text: "Text for enhancement",
			}),
			false,
		},
		{
			"Malformed control returns an error",
			args{"SC-43a)"},
			v3c.Satisfies{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.parseControl(tt.args.control)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseControl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseControl() = %v, want %v", got, tt.want)
			}
		})
	}
}

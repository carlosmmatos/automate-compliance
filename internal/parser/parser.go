package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/carlosmmatos/automate-compliance/internal/opencontrol"
)

type Parser struct {
	// whitespace regex
	wre *regexp.Regexp
	// simple control
	simpleCtrl *regexp.Regexp
	// control with enhancements
	ctrlEnh *regexp.Regexp
	// subcontrol with enhancements
	simpleSubCtrl *regexp.Regexp
	// subcontrol with enhancements
	subCtrlEnh *regexp.Regexp
	data       map[opencontrol.Family]map[string]opencontrol.OpenControlEntry
}

func NewParser() *Parser {
	return &Parser{
		// whitespace regex
		wre: regexp.MustCompile(`\s+`),
		// simple control ([1] control family [2] control number)
		simpleCtrl: regexp.MustCompile(`^([A-Z]+)-([0-9]+)$`),
		// control with enhancements ([1] control family [2] control number [3] enhancement)
		ctrlEnh: regexp.MustCompile(`^([A-Z]+)-([0-9]+)([a-z]\.([1-9]\.)?)$`),
		// subcontrol with enhancements ([1] control family [2] control number [3] sub-control)
		simpleSubCtrl: regexp.MustCompile(`^([A-Z]+)-([0-9]+) (\([0-9]+\))$`),
		// subcontrol with enhancements ([1] control family [2] control number [3] sub-control [4] enhancement)
		// NOTE(jaosorior): This ignores any sub-entries in the enhancement... so it'll
		// match AC-3 (3)(b)(1) and AC-3 (3)(b)(2) as the same entry -> AC-3 (3)(b)
		subCtrlEnh: regexp.MustCompile(`^([A-Z]+)-([0-9]+) (\([0-9]+\))\(([a-z])\).*$`),
		data:       make(map[opencontrol.Family]map[string]opencontrol.OpenControlEntry),
	}
}

func (p *Parser) ParseEntry(family, control string) error {
	nfamily := p.normalizeFamily(family)

	ctrls, foundFam := p.data[nfamily]

	if !foundFam {
		// initialize control entries
		ctrls = make(map[string]opencontrol.OpenControlEntry)
		p.data[nfamily] = ctrls
	}

	parsedCtrl, err := p.parseControl(control)
	if err != nil {
		return err
	}

	storedCtrl, foundCtrl := ctrls[parsedCtrl.ControlKey]

	if !foundCtrl {
		ctrls[parsedCtrl.ControlKey] = parsedCtrl
		return nil
	}

	ctrls[parsedCtrl.ControlKey] = mergeControls(storedCtrl, parsedCtrl)
	return nil
}

// normalizeFamily normalizes the family name into something more
// fitting for OpenControl.
//
// NOTE(jaosorior): This currently only replaces spaces for underscores...
// we should probably replace this function with something that gets a
// standardized name somehow
func (p *Parser) normalizeFamily(family string) opencontrol.Family {
	return opencontrol.Family(p.wre.ReplaceAllString(family, "_"))
}

func (p *Parser) parseControl(control string) (opencontrol.OpenControlEntry, error) {
	// control without extra whitespaces
	ctrlNw := p.wre.ReplaceAllString(control, " ")

	ctrl := opencontrol.ControlEntryWithDefaults()
	if p.simpleCtrl.MatchString(ctrlNw) {
		matches := p.simpleCtrl.FindStringSubmatch(ctrlNw)
		ctrl.ControlKey = getControlKey(matches)
		ctrl.Narrative = append(ctrl.Narrative, getTextOnlyNarrative())
		return ctrl, nil
	} else if p.ctrlEnh.MatchString(ctrlNw) {
		matches := p.ctrlEnh.FindStringSubmatch(ctrlNw)
		ctrl.ControlKey = getControlKey(matches)
		ctrl.Narrative = append(ctrl.Narrative, getNarrativeForEnhancement(matches[3]))
		return ctrl, nil
	} else if p.simpleSubCtrl.MatchString(ctrlNw) {
		matches := p.simpleSubCtrl.FindStringSubmatch(ctrlNw)
		ctrl.ControlKey = getSubControlKey(matches)
		ctrl.Narrative = append(ctrl.Narrative, getTextOnlyNarrative())
		return ctrl, nil
	} else if p.subCtrlEnh.MatchString(ctrlNw) {
		matches := p.subCtrlEnh.FindStringSubmatch(ctrlNw)
		ctrl.ControlKey = getSubControlKey(matches)
		ctrl.Narrative = append(ctrl.Narrative, getNarrativeForEnhancement(matches[4]))
		return ctrl, nil
	} else {
		// no match
		return ctrl, fmt.Errorf("Couldn't parse control")
	}
}

func (p *Parser) GetData() map[opencontrol.Family]map[string]opencontrol.OpenControlEntry {
	return p.data
}

func getControlKey(matches []string) string {
	return fmt.Sprintf("%s-%s", matches[1], matches[2])
}

func getSubControlKey(matches []string) string {
	return fmt.Sprintf("%s-%s %s", matches[1], matches[2], matches[3])
}

func getTextOnlyNarrative() opencontrol.NarrativeEntry {
	// TODO(jaosorior): get text from spreadsheet
	return opencontrol.NarrativeEntry{
		Text: "Text only",
	}
}

func getNarrativeForEnhancement(enhancement string) opencontrol.NarrativeEntry {
	// TODO(jaosorior): get text from spreadsheet
	return opencontrol.NarrativeEntry{
		Key:  normalizeEnhacementKey(enhancement),
		Text: "Text for enhancement",
	}
}

func normalizeEnhacementKey(e string) string {
	return strings.TrimRight(e, ".")
}

func mergeControls(old, new opencontrol.OpenControlEntry) opencontrol.OpenControlEntry {
	// The controlKey is the same so we don't need to merge these.

	// TODO(jaosorior): Gotta handle implementation status

	old.Narrative = append(old.Narrative, new.Narrative...)
	return old
}

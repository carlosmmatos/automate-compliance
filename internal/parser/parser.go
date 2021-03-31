package parser

import (
	"fmt"
	"regexp"
	"strings"

	v3c "github.com/opencontrol/compliance-masonry/pkg/lib/components/versions/3_1_0"
)

type controlFamily string

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
	subCtrlEnh     *regexp.Regexp
	// subcontrol with extra enhancements
	subCtrlEnhPlus *regexp.Regexp
	data           map[controlFamily]map[string]v3c.Satisfies
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
		subCtrlEnh: regexp.MustCompile(`^([A-Z]+)-([0-9]+) (\([0-9]+\))\(([a-z])\)$`),
		// subcontrol with additional enhancements ([1] control family [2] control number [3] sub-control 
		// [4] enhancement + [5] additional_enhancement)
		subCtrlEnhPlus: regexp.MustCompile(`^([A-Z]+)-([0-9]+) (\([0-9]+\))\(([a-z])\)\(([0-9]+)\)$`),
		data:       make(map[controlFamily]map[string]v3c.Satisfies),
	}
}

func (p *Parser) ParseEntry(family, control string) error {
	nfamily := p.normalizeFamily(family)

	ctrls, foundFam := p.data[nfamily]

	if !foundFam {
		// initialize control entries
		ctrls = make(map[string]v3c.Satisfies)
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
func (p *Parser) normalizeFamily(family string) controlFamily {
	// Ensure we take care of whitespace issues
	nFamily := p.wre.ReplaceAllString(family, "_")
	return controlFamily(ParseFamily(nFamily))
}

// parseControl parses a NIST 800-53 control and ensures it conforms to the 
// OpenControl Satisfies struct.
func (p *Parser) parseControl(control string) (v3c.Satisfies, error) {
	// control without extra whitespaces
	ctrlNw := p.wre.ReplaceAllString(control, " ")

	ctrl := v3c.Satisfies{}
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
	} else if p.subCtrlEnhPlus.MatchString(ctrlNw) {
		matches := p.subCtrlEnhPlus.FindStringSubmatch(ctrlNw)
		ctrl.ControlKey = getSubControlKey(matches)
		ctrl.Narrative = append(ctrl.Narrative, getNarrativeForEnhancementPlus(matches))
		return ctrl, nil
	} else {
		// no match
		return ctrl, fmt.Errorf("couldn't parse control")
	}
}

func (p *Parser) GetData() map[controlFamily]map[string]v3c.Satisfies {
	removeNarrative(p.data)
	return p.data
}

// removeNarrative removes the first text Narrative from controls with additional controls
// and enhancements. Fixes https://github.com/carlosmmatos/automate-compliance/issues/12
func removeNarrative(p map[controlFamily]map[string]v3c.Satisfies) {
	for _, ctrl := range p {
		for _, v := range ctrl {
			if len(v.Narrative) > 1 {
				// when length of Narrative > 1, we know we can safely remove the first
				// index, which would be for the primary control key.
				ctrl[v.ControlKey] = remove(v)
			}
		}
	}
}

// remove deletes the first element of a slice, while maintaining order
func remove(old v3c.Satisfies) v3c.Satisfies {
	old.Narrative = append(old.Narrative[:0], old.Narrative[1:]...)
	return old
}

func getControlKey(matches []string) string {
	return fmt.Sprintf("%s-%s", matches[1], matches[2])
}

func getSubControlKey(matches []string) string {
	return fmt.Sprintf("%s-%s %s", matches[1], matches[2], matches[3])
}

func getTextOnlyNarrative() v3c.NarrativeSection {
	// TODO(jaosorior): get text from spreadsheet
	return v3c.NarrativeSection{
		Text: "Text only",
	}
}

func getNarrativeForEnhancement(enhancement string) v3c.NarrativeSection {
	// TODO(jaosorior): get text from spreadsheet
	return v3c.NarrativeSection{
		Key:  normalizeEnhancementKey(enhancement),
		Text: "Text for enhancement",
	}
}

func getNarrativeForEnhancementPlus(matches []string) v3c.NarrativeSection {
	// handles use case: AC-3 (3)(b)(2)
	return v3c.NarrativeSection{
		Key: normalizeEnhancementPlusKey(matches),
		Text: "Text for enhancement plus",
	}
}

func normalizeEnhancementKey(e string) string {
	return strings.TrimRight(e, ".")
}

func normalizeEnhancementPlusKey(matches []string) string {
	// Return b.2 from AC-3 (3)(b)(2)
	return fmt.Sprintf("%s.%s", matches[4], matches[5])
}

func mergeControls(old, new v3c.Satisfies) v3c.Satisfies {
	// The controlKey is the same so we don't need to merge these.

	// TODO(jaosorior): Gotta handle implementation status

	old.Narrative = append(old.Narrative, new.Narrative...)
	return old
}

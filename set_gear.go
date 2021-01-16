package main

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func setGear(context ExecutionContext) error {

	var err error

	recursive := *context.Options["recursive"].(*bool)

	matches, _ := resolveToMatches(context.Args[0], recursive, true)

	nameValuePair := strings.Split(context.Args[1], "=")

	path := strings.Split(nameValuePair[0], ".")
	newValue := nameValuePair[1]

	for _, match := range matches {

		source, _ := filepath.Abs(match)

		sourceFile, _ := ioutil.ReadFile(source)

		var preset PresetXMLV5Raw
		_ = xml.Unmarshal(sourceFile, &preset)

		attrName := path[len(path)-1]

		switch path[len(path)-2] {
		case "Preset":
			updateAttr(&preset.Attrs, attrName, newValue)
		case "Chain":
			updateAttr(&preset.Chain.Attrs, attrName, newValue)
		case "Input":
			updateAttr(&preset.Input.Attrs, attrName, newValue)
		case "Tuner":
			if path[len(path)-3] == "Preset" {
				updateAttr(&preset.Tuner.Attrs, attrName, newValue)
			} else {
				updateAttr(&preset.Tuner.Tuner.Attrs, attrName, newValue)
			}
		case "Slot0":
			switch path[len(path)-3] {
			case "StompA1":
				err = updateAttr(&preset.StompA1.Slot0.Attrs, attrName, newValue)
			case "StompA2":
				err = updateAttr(&preset.StompA2.Slot0.Attrs, attrName, newValue)
			case "StompStereo":
				err = updateAttr(&preset.StompStereo.Slot0.Attrs, attrName, newValue)
			case "StompB1":
				err = updateAttr(&preset.StompB1.Slot0.Attrs, attrName, newValue)
			case "StompB2":
				err = updateAttr(&preset.StompB2.Slot0.Attrs, attrName, newValue)
			case "StompB3":
				err = updateAttr(&preset.StompB3.Slot0.Attrs, attrName, newValue)
			case "LoopFxA":
				err = updateAttr(&preset.LoopFxA.Slot0.Attrs, attrName, newValue)
			case "LoopFxB":
				err = updateAttr(&preset.LoopFxB.Slot0.Attrs, attrName, newValue)
			case "LoopFxC":
				err = updateAttr(&preset.LoopFxC.Slot0.Attrs, attrName, newValue)
			case "RackA":
				err = updateAttr(&preset.RackA.Slot0.Attrs, attrName, newValue)
			case "RackB":
				err = updateAttr(&preset.RackB.Slot0.Attrs, attrName, newValue)
			case "RackC":
				err = updateAttr(&preset.RackC.Slot0.Attrs, attrName, newValue)
			case "RackDI":
				err = updateAttr(&preset.RackDI.Slot0.Attrs, attrName, newValue)
			case "RackMaster":
				err = updateAttr(&preset.RackMaster.Slot0.Attrs, attrName, newValue)
			default:
				err = errors.New("invalid or unsupported path " + strings.Join(path, "."))
			}
		case "Slot1":
			switch path[len(path)-3] {
			case "StompA1":
				err = updateAttr(&preset.StompA1.Slot1.Attrs, attrName, newValue)
			case "StompA2":
				err = updateAttr(&preset.StompA2.Slot1.Attrs, attrName, newValue)
			case "StompStereo":
				err = updateAttr(&preset.StompStereo.Slot1.Attrs, attrName, newValue)
			case "StompB1":
				err = updateAttr(&preset.StompB1.Slot1.Attrs, attrName, newValue)
			case "StompB2":
				err = updateAttr(&preset.StompB2.Slot1.Attrs, attrName, newValue)
			case "StompB3":
				err = updateAttr(&preset.StompB3.Slot1.Attrs, attrName, newValue)
			case "LoopFxA":
				err = updateAttr(&preset.LoopFxA.Slot1.Attrs, attrName, newValue)
			case "LoopFxB":
				err = updateAttr(&preset.LoopFxB.Slot1.Attrs, attrName, newValue)
			case "LoopFxC":
				err = updateAttr(&preset.LoopFxC.Slot1.Attrs, attrName, newValue)
			case "RackA":
				err = updateAttr(&preset.RackA.Slot1.Attrs, attrName, newValue)
			case "RackB":
				err = updateAttr(&preset.RackB.Slot1.Attrs, attrName, newValue)
			case "RackC":
				err = updateAttr(&preset.RackC.Slot1.Attrs, attrName, newValue)
			case "RackDI":
				err = updateAttr(&preset.RackDI.Slot1.Attrs, attrName, newValue)
			case "RackMaster":
				err = updateAttr(&preset.RackMaster.Slot1.Attrs, attrName, newValue)
			default:
				err = errors.New("invalid or unsupported path " + strings.Join(path, "."))
			}
		case "Slot2":
			switch path[len(path)-3] {
			case "StompA1":
				err = updateAttr(&preset.StompA1.Slot2.Attrs, attrName, newValue)
			case "StompA2":
				err = updateAttr(&preset.StompA2.Slot2.Attrs, attrName, newValue)
			case "StompStereo":
				err = updateAttr(&preset.StompStereo.Slot2.Attrs, attrName, newValue)
			case "StompB1":
				err = updateAttr(&preset.StompB1.Slot2.Attrs, attrName, newValue)
			case "StompB2":
				err = updateAttr(&preset.StompB2.Slot2.Attrs, attrName, newValue)
			case "StompB3":
				err = updateAttr(&preset.StompB3.Slot2.Attrs, attrName, newValue)
			case "LoopFxA":
				err = updateAttr(&preset.LoopFxA.Slot2.Attrs, attrName, newValue)
			case "LoopFxB":
				err = updateAttr(&preset.LoopFxB.Slot2.Attrs, attrName, newValue)
			case "LoopFxC":
				err = updateAttr(&preset.LoopFxC.Slot2.Attrs, attrName, newValue)
			case "RackMaster":
				err = updateAttr(&preset.RackMaster.Slot2.Attrs, attrName, newValue)
			default:
				err = errors.New("invalid or unsupported path " + strings.Join(path, "."))
			}
		case "Slot3":
			switch path[len(path)-3] {
			case "StompA1":
				err = updateAttr(&preset.StompA1.Slot3.Attrs, attrName, newValue)
			case "StompA2":
				err = updateAttr(&preset.StompA2.Slot3.Attrs, attrName, newValue)
			case "StompB1":
				err = updateAttr(&preset.StompB1.Slot3.Attrs, attrName, newValue)
			case "StompB2":
				err = updateAttr(&preset.StompB2.Slot3.Attrs, attrName, newValue)
			case "StompB3":
				err = updateAttr(&preset.StompB3.Slot3.Attrs, attrName, newValue)
			case "LoopFxA":
				err = updateAttr(&preset.LoopFxA.Slot3.Attrs, attrName, newValue)
			case "LoopFxB":
				err = updateAttr(&preset.LoopFxB.Slot3.Attrs, attrName, newValue)
			case "LoopFxC":
				err = updateAttr(&preset.LoopFxC.Slot3.Attrs, attrName, newValue)
			case "RackMaster":
				err = updateAttr(&preset.RackMaster.Slot3.Attrs, attrName, newValue)
			default:
				err = errors.New("invalid or unsupported path " + strings.Join(path, "."))
			}
		case "Slot4":
			switch path[len(path)-3] {
			case "StompA1":
				err = updateAttr(&preset.StompA1.Slot4.Attrs, attrName, newValue)
			case "StompA2":
				err = updateAttr(&preset.StompA2.Slot4.Attrs, attrName, newValue)
			case "StompB1":
				err = updateAttr(&preset.StompB1.Slot4.Attrs, attrName, newValue)
			case "StompB2":
				err = updateAttr(&preset.StompB2.Slot4.Attrs, attrName, newValue)
			case "StompB3":
				err = updateAttr(&preset.StompB3.Slot4.Attrs, attrName, newValue)
			case "RackMaster":
				err = updateAttr(&preset.RackMaster.Slot4.Attrs, attrName, newValue)
			default:
				err = errors.New("invalid or unsupported path " + strings.Join(path, "."))
			}
		case "Slot5":
			switch path[len(path)-3] {
			case "StompA1":
				err = updateAttr(&preset.StompA1.Slot5.Attrs, attrName, newValue)
			case "StompA2":
				err = updateAttr(&preset.StompA2.Slot5.Attrs, attrName, newValue)
			case "StompB1":
				err = updateAttr(&preset.StompB1.Slot5.Attrs, attrName, newValue)
			case "StompB2":
				err = updateAttr(&preset.StompB2.Slot5.Attrs, attrName, newValue)
			case "StompB3":
				err = updateAttr(&preset.StompB3.Slot5.Attrs, attrName, newValue)
			case "RackMaster":
				err = updateAttr(&preset.RackMaster.Slot5.Attrs, attrName, newValue)
			default:
				err = errors.New("invalid or unsupported path " + strings.Join(path, "."))
			}
		case "Amp":
			switch path[len(path)-3] {
			case "AmpA":
				err = updateAttr(&preset.AmpA.Amp.Attrs, attrName, newValue)
			case "AmpB":
				err = updateAttr(&preset.AmpB.Amp.Attrs, attrName, newValue)
			case "AmpC":
				err = updateAttr(&preset.AmpC.Amp.Attrs, attrName, newValue)
			default:
				err = errors.New("invalid or unsupported path " + strings.Join(path, "."))
			}
		case "AmpA":
			err = updateAttr(&preset.AmpA.Attrs, attrName, newValue)
		case "AmpB":
			err = updateAttr(&preset.AmpB.Attrs, attrName, newValue)
		case "AmpC":
			err = updateAttr(&preset.AmpC.Attrs, attrName, newValue)
		case "Cab":
			switch path[len(path)-3] {
			case "CabA":
				err = updateAttr(&preset.CabA.Cab.Attrs, attrName, newValue)
			case "CabB":
				err = updateAttr(&preset.CabB.Cab.Attrs, attrName, newValue)
			case "CabC":
				err = updateAttr(&preset.CabC.Cab.Attrs, attrName, newValue)
			default:
				err = errors.New("invalid or unsupported path " + strings.Join(path, "."))
			}
		case "CabA":
			err = updateAttr(&preset.CabA.Attrs, attrName, newValue)
		case "CabB":
			err = updateAttr(&preset.CabB.Attrs, attrName, newValue)
		case "CabC":
			err = updateAttr(&preset.CabC.Attrs, attrName, newValue)
		case "Studio":
			err = updateAttr(&preset.Studio.Attrs, attrName, newValue)
		case "Output":
			err = updateAttr(&preset.Output.Attrs, attrName, newValue)
		case "MidiAssignments":
			err = updateAttr(&preset.MidiAssignments.Attrs, attrName, newValue)
		case "MetaInfo":
			err = updateAttr(&preset.MetaInfo.Attrs, attrName, newValue)
		default:
			err = errors.New("invalid or unsupported path " + strings.Join(path, "."))
		}

		if err == nil {
			newData, _ := xml.MarshalIndent(preset, "", "    ")
			ioutil.WriteFile(source, selfClose(append([]byte(XMLPrefix), append(newData[:], '\n')...)), 0664)
		}

	}

	return err
}

func updateAttr(attrs *[]xml.Attr, name string, value string) error {
	updated := false
	for i, thisAttr := range *attrs {
		if thisAttr.Name.Local == name {
			(*attrs)[i].Value = value
			updated = true
		}
	}
	if !updated {
		return errors.New("attribute not found: " + name)
	} else {
		return nil
	}
}

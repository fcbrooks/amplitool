package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

func listGear(context ExecutionContext) error {

	if len(context.Args) < 1 {
		return errors.New("list gear requires a preset")
	}

	raw := *context.Options["raw"].(*bool)
	details := *context.Options["details"].(*bool)

	source, err := filepath.Abs(context.Args[0])

	if !isFile(source) || !isInPresetsFolder(source) {
		return errors.New("first argument must be a preset file")
	}

	format, err := presetFormatVersion(source)

	sourceFile, err := ioutil.ReadFile(source)

	if err != nil {
		return err
	}

	if raw || format != "at5p" {
		fmt.Fprintln(out, string(sourceFile))
	} else {

		var sourcePreset PresetXMLV5
		err = xml.Unmarshal(sourceFile, &sourcePreset)

		var cabAttrs CabAttrs
		err = xml.Unmarshal(sourceFile, &cabAttrs)

		if err != nil {
			return err
		}

		var indent []string
		var output string

		output = fmt.Sprintln(strings.Join(indent, "") + filepath.Base(source))

		indent = append(indent, "    ")

		if sourcePreset.Chain.Preset == "Chain11" {
			printAmp(&output, fromAmpA(sourcePreset.AmpA), "AmpA", indent, details)
			printCab(&output, fromCabA(sourcePreset.CabA), "CabA", indent, cabAttrs.CabA.Attrs, details)
			printStomp(&output, fromStompA1(sourcePreset.StompA1), "StompA1", indent, details)
			printStomp(&output, fromStompB1(sourcePreset.StompB1), "StompB1", indent, details)
			printStomp(&output, fromLoopFxA(sourcePreset.LoopFxA), "LoopFxA", indent, details)
			printStomp(&output, fromRackA(sourcePreset.RackA), "RackA", indent, details)
			printStomp(&output, fromRackDI(sourcePreset.RackDI), "RackDI", indent, details)
			printStomp(&output, fromRackMaster(sourcePreset.RackMaster), "RackMaster", indent, details)
		}

		if sourcePreset.Chain.Preset == "Chain12" {
			printAmp(&output, fromAmpA(sourcePreset.AmpA), "AmpA", indent, details)
			printCab(&output, fromCabA(sourcePreset.CabA), "CabA", indent, cabAttrs.CabA.Attrs, details)
			printAmp(&output, fromAmpB(sourcePreset.AmpB), "AmpB", indent, details)
			printCab(&output, fromCabB(sourcePreset.CabB), "CabB", indent, cabAttrs.CabB.Attrs, details)
			printStomp(&output, fromStompA1(sourcePreset.StompA1), "StompA1", indent, details)
			printStomp(&output, fromStompB1(sourcePreset.StompB1), "StompB1", indent, details)
			printStomp(&output, fromStompB2(sourcePreset.StompB2), "StompB2", indent, details)
			printStomp(&output, fromLoopFxA(sourcePreset.LoopFxA), "LoopFxA", indent, details)
			printStomp(&output, fromLoopFxB(sourcePreset.LoopFxB), "LoopFxB", indent, details)
			printStomp(&output, fromRackA(sourcePreset.RackA), "RackA", indent, details)
			printStomp(&output, fromRackB(sourcePreset.RackB), "RackB", indent, details)
			printStomp(&output, fromRackDI(sourcePreset.RackDI), "RackDI", indent, details)
			printStomp(&output, fromRackMaster(sourcePreset.RackMaster), "RackMaster", indent, details)
		}

		if sourcePreset.Chain.Preset == "Chain13" {
			printAmp(&output, fromAmpA(sourcePreset.AmpA), "AmpA", indent, details)
			printCab(&output, fromCabA(sourcePreset.CabA), "CabA", indent, cabAttrs.CabA.Attrs, details)
			printAmp(&output, fromAmpB(sourcePreset.AmpB), "AmpB", indent, details)
			printCab(&output, fromCabB(sourcePreset.CabB), "CabB", indent, cabAttrs.CabB.Attrs, details)
			printAmp(&output, fromAmpC(sourcePreset.AmpC), "AmpC", indent, details)
			printCab(&output, fromCabC(sourcePreset.CabC), "CabC", indent, cabAttrs.CabC.Attrs, details)
			printStomp(&output, fromStompA1(sourcePreset.StompA1), "StompA1", indent, details)
			printStomp(&output, fromStompB1(sourcePreset.StompB1), "StompB1", indent, details)
			printStomp(&output, fromStompB2(sourcePreset.StompB2), "StompB2", indent, details)
			printStomp(&output, fromStompB3(sourcePreset.StompB3), "StompB3", indent, details)
			printStomp(&output, fromLoopFxA(sourcePreset.LoopFxA), "LoopFxA", indent, details)
			printStomp(&output, fromLoopFxB(sourcePreset.LoopFxB), "LoopFxB", indent, details)
			printStomp(&output, fromLoopFxC(sourcePreset.LoopFxC), "LoopFxC", indent, details)
			printStomp(&output, fromRackA(sourcePreset.RackA), "RackA", indent, details)
			printStomp(&output, fromRackB(sourcePreset.RackB), "RackB", indent, details)
			printStomp(&output, fromRackC(sourcePreset.RackC), "RackC", indent, details)
			printStomp(&output, fromRackDI(sourcePreset.RackDI), "RackDI", indent, details)
			printStomp(&output, fromRackMaster(sourcePreset.RackMaster), "RackMaster", indent, details)
		}

		if sourcePreset.Chain.Preset == "Chain22" {
			printAmp(&output, fromAmpA(sourcePreset.AmpA), "AmpA", indent, details)
			printCab(&output, fromCabA(sourcePreset.CabA), "CabA", indent, cabAttrs.CabA.Attrs, details)
			printAmp(&output, fromAmpB(sourcePreset.AmpB), "AmpB", indent, details)
			printCab(&output, fromCabB(sourcePreset.CabB), "CabB", indent, cabAttrs.CabB.Attrs, details)
			printStomp(&output, fromStompA1(sourcePreset.StompA1), "StompA1", indent, details)
			printStomp(&output, fromStompA2(sourcePreset.StompA2), "StompA2", indent, details)
			printStomp(&output, fromStompB1(sourcePreset.StompB1), "StompB1", indent, details)
			printStomp(&output, fromStompB2(sourcePreset.StompB2), "StompB2", indent, details)
			printStomp(&output, fromLoopFxA(sourcePreset.LoopFxA), "LoopFxA", indent, details)
			printStomp(&output, fromLoopFxB(sourcePreset.LoopFxB), "LoopFxB", indent, details)
			printStomp(&output, fromRackA(sourcePreset.RackA), "RackA", indent, details)
			printStomp(&output, fromRackB(sourcePreset.RackB), "RackB", indent, details)
			printStomp(&output, fromRackDI(sourcePreset.RackDI), "RackDI", indent, details)
			printStomp(&output, fromRackMaster(sourcePreset.RackMaster), "RackMaster", indent, details)
		}

		fmt.Fprintln(out, output)

	}

	return nil
}

func printAmp(output *string, amp GenericAmp, ampType string, indent []string, details bool) {
	*output = fmt.Sprintln(*output + strings.Join(indent, "") + ampType + ": " + getValueOrKey(Amps, amp.Model))
	if details {
		indent = append(indent, "    ")
		for _, a := range amp.Amp.Attrs {
			*output = fmt.Sprintln(*output + strings.Join(indent, "") + " " + a.Name.Local + ": " + a.Value)
		}
		indent = indent[:len(indent)-1]
	}
}

func printCab(output *string, cab GenericCab, cabType string, indent []string, attrs []xml.Attr, details bool) {
	*output = fmt.Sprintln(*output + strings.Join(indent, "") + cabType + ": " + getValueOrKey(Cabs, cab.CabModel))
	indent = append(indent, "    ")
	for _, a := range attrs {
		if details {
			*output = fmt.Sprintln(*output + strings.Join(indent, "") + " " + a.Name.Local + ": " + a.Value)
		} else {
			if strings.Index(a.Name.Local, "SpeakerModel") == 0 {
				speakerCount := 4
				if SpeakerCount[cab.CabModel] != 0 {
					speakerCount = SpeakerCount[cab.CabModel]
				}
				speakerNumber, _ := strconv.Atoi(a.Name.Local[12:])
				if speakerNumber < speakerCount {
					*output = fmt.Sprintln(*output + strings.Join(indent, "") + " " + a.Name.Local + ": " + getValueOrKey(Speakers, a.Value))
				}
			}
		}
	}
	for _, a := range cab.Cab.Attrs {
		if details {
			if a.Name.Local == "Mic0Model" || a.Name.Local == "Mic1Model" {
				*output = fmt.Sprintln(*output + strings.Join(indent, "") + " " + a.Name.Local + ": " + getValueOrKey(Mics, a.Value))
			} else if strings.Index(a.Name.Local, "SpeakerModel") == 0 {
				speakerCount := 4
				if SpeakerCount[cab.CabModel] != 0 {
					speakerCount = SpeakerCount[cab.CabModel]
				}
				speakerNumber, _ := strconv.Atoi(a.Name.Local[12:])
				if speakerNumber < speakerCount {
					*output = fmt.Sprintln(*output + strings.Join(indent, "") + " " + a.Name.Local + ": " + getValueOrKey(Speakers, a.Value))
				}
			} else {
				*output = fmt.Sprintln(*output + strings.Join(indent, "") + " " + a.Name.Local + ": " + a.Value)
			}
		} else {
			if a.Name.Local == "Mic0Model" || a.Name.Local == "Mic1Model" || a.Name.Local == "RoomType" || a.Name.Local == "RoomMicType" || strings.Index(a.Name.Local, "SpeakerModel") == 0 {
				if a.Name.Local == "Mic0Model" || a.Name.Local == "Mic1Model" {
					*output = fmt.Sprintln(*output + strings.Join(indent, "") + " " + a.Name.Local + ": " + getValueOrKey(Mics, a.Value))
				} else if strings.Index(a.Name.Local, "SpeakerModel") == 0 {
					speakerCount := 4
					if SpeakerCount[cab.CabModel] != 0 {
						speakerCount = SpeakerCount[cab.CabModel]
					}
					speakerNumber, _ := strconv.Atoi(a.Name.Local[12:])
					if speakerNumber < speakerCount {
						*output = fmt.Sprintln(*output + strings.Join(indent, "") + " " + a.Name.Local + ": " + getValueOrKey(Speakers, a.Value))
					}
				} else {
					*output = fmt.Sprintln(*output + strings.Join(indent, "") + " " + a.Name.Local + ": " + a.Value)
				}
			}
		}

	}
	indent = indent[:len(indent)-1]
}

func printStomp(output *string, stomp GenericStomp, stompType string, indent []string, details bool) {
	if !allStompsEmpty(stomp) {
		*output = fmt.Sprintln(*output + strings.Join(indent, "") + stompType)
		indent = append(indent, "    ")
		for i := 0; i < stomp.StompCount; i++ {
			switch i {
			case 0:
				if stomp.Stomp0 != EmptySlotGUID {
					*output = fmt.Sprintln(*output + strings.Join(indent, "") + "Slot" + strconv.Itoa(i) + ": " + getValueOrKey(FX, stomp.Stomp0))
					if details {
						printSlot(output, GenericSlot{stomp.Slot0.Attrs}, "Slot0", indent)
					}
				}
			case 1:
				if stomp.Stomp1 != EmptySlotGUID {
					*output = fmt.Sprintln(*output + strings.Join(indent, "") + "Slot" + strconv.Itoa(i) + ": " + getValueOrKey(FX, stomp.Stomp1))
					if details {
						printSlot(output, GenericSlot{stomp.Slot1.Attrs}, "Slot"+strconv.Itoa(i), indent)
					}
				}
			case 2:
				if stomp.Stomp2 != EmptySlotGUID {
					*output = fmt.Sprintln(*output + strings.Join(indent, "") + "Slot" + strconv.Itoa(i) + ": " + getValueOrKey(FX, stomp.Stomp2))
					if details {
						printSlot(output, GenericSlot{stomp.Slot1.Attrs}, "Slot"+strconv.Itoa(i), indent)
					}
				}
			case 3:
				if stomp.Stomp3 != EmptySlotGUID {
					*output = fmt.Sprintln(*output + strings.Join(indent, "") + "Slot" + strconv.Itoa(i) + ": " + getValueOrKey(FX, stomp.Stomp3))
					if details {
						printSlot(output, GenericSlot{stomp.Slot1.Attrs}, "Slot"+strconv.Itoa(i), indent)
					}
				}
			case 4:
				if stomp.Stomp4 != EmptySlotGUID {
					*output = fmt.Sprintln(*output + strings.Join(indent, "") + "Slot" + strconv.Itoa(i) + ": " + getValueOrKey(FX, stomp.Stomp4))
					if details {
						printSlot(output, GenericSlot{stomp.Slot1.Attrs}, "Slot"+strconv.Itoa(i), indent)
					}
				}
			case 5:
				if stomp.Stomp5 != EmptySlotGUID {
					*output = fmt.Sprintln(*output + strings.Join(indent, "") + "Slot" + strconv.Itoa(i) + ": " + getValueOrKey(FX, stomp.Stomp5))
					if details {
						printSlot(output, GenericSlot{stomp.Slot1.Attrs}, "Slot"+strconv.Itoa(i), indent)
					}
				}
			}
		}
		indent = indent[:len(indent)-1]
	}
}

func printSlot(output *string, slot GenericSlot, slotType string, indent []string) {
	indent = append(indent, "    ")
	for _, a := range slot.Attrs {
		*output = fmt.Sprintln(*output + strings.Join(indent, "") + " " + a.Name.Local + ": " + a.Value)
	}
	indent = indent[:len(indent)-1]
}

func allStompsEmpty(stomp GenericStomp) bool {
	return (stomp.Stomp0 == EmptySlotGUID || stomp.Stomp0 == "") && (stomp.Stomp1 == EmptySlotGUID || stomp.Stomp1 == "") && (stomp.Stomp2 == EmptySlotGUID || stomp.Stomp2 == "") && (stomp.Stomp3 == EmptySlotGUID || stomp.Stomp3 == "") && (stomp.Stomp4 == EmptySlotGUID || stomp.Stomp4 == "") && (stomp.Stomp5 == EmptySlotGUID || stomp.Stomp5 == "")
}

func getValueOrKey(valueMap map[string]string, key string) string {
	value := key
	if valueMap[key] != "" {
		value = valueMap[key]
	}
	return value
}

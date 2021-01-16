/*
Copyright (C) 2021 fcbrooks

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func removeGear(context ExecutionContext) error {
	return copyOrRemove(context, true)
}

func copyGear(context ExecutionContext) error {
	return copyOrRemove(context, false)
}

func copyOrRemove(context ExecutionContext, removeFx bool) error {

	copyAllAmps := false
	copyAllCabs := false
	copyAllFx := false
	overwriteFx := false
	appendFx := false
	insertFx := false
	copyCabWithAmp := true

	if !removeFx {
		copyAllAmps = *context.Options["allamps"].(*bool)
		copyAllCabs = *context.Options["allcabs"].(*bool)
		copyAllFx = *context.Options["allfx"].(*bool)
		overwriteFx = *context.Options["overwritefx"].(*bool)
		insertFx = *context.Options["insertfx"].(*bool)
		appendFx = !overwriteFx && !insertFx && !copyAllFx
		copyCabWithAmp = !*context.Options["nocabwithamp"].(*bool)
	}

	if removeFx {
		if len(context.Args) < 2 {
			return errors.New("remove gear requires a source preset and gear type")
		}
	} else {
		if !copyAllAmps && !copyAllCabs && !copyAllFx {
			if len(context.Args) < 3 {
				return errors.New("copy gear requires a source preset, destination, and source gear name")
			}
		} else {
			if (appendFx || insertFx) && copyAllFx {
				return errors.New("appending or inserting effects is not compatible with copy all effects")
			}
			if len(context.Args) < 2 {
				return errors.New("copying all requires a source preset and destination")
			}
		}
	}

	source, err := filepath.Abs(context.Args[0])

	if !isFile(source) || !isInPresetsFolder(source) {
		return errors.New("first argument must be a preset file")
	}

	recursive := *context.Options["recursive"].(*bool)

	if !removeFx {

		target, _ := filepath.Abs(context.Args[1])

		if isDir(target) && !recursive {
			fmt.Fprintln(out, "-r not specified; omitting directory")
			return nil
		}

	}

	format, err := presetFormatVersion(source)

	if format != "at5p" {
		return errors.New("copy gear only supported for version 5 presets")
	}

	sourceFile, err := ioutil.ReadFile(source)

	if err != nil {
		return err
	}

	var sourcePreset PresetXMLV5
	err = xml.Unmarshal(sourceFile, &sourcePreset)

	if err != nil {
		return err
	}

	matches := []string{}

	if removeFx {
		matches, err = resolveToMatches(context.Args[0], recursive, true)
	} else {
		matches, err = resolveToMatches(context.Args[1], recursive, true)
	}

	gearMap := map[string]string{}
	targetSlot := ""

	if !copyAllAmps && !copyAllCabs && !copyAllFx {
		if len(context.Args) > 3 {
			gearMap[context.Args[2]] = context.Args[3]
		} else {
			if removeFx {
				gearMap[context.Args[1]] = ""
				if len(context.Args) > 2 {
					targetSlot = context.Args[2]
				}
			} else {
				gearMap[context.Args[2]] = ""
			}
		}
	} else {
		if copyAllAmps {
			gearMap["AmpA"] = "AmpA"
			gearMap["AmpB"] = "AmpB"
			gearMap["AmpC"] = "AmpC"
		}
		if copyAllCabs {
			gearMap["CabA"] = "CabA"
			gearMap["CabB"] = "CabB"
			gearMap["CabC"] = "CabC"
		}
		if copyAllFx {
			gearMap["StompA1"] = "StompA1"
			gearMap["StompA2"] = "StompA2"
			gearMap["StompStereo"] = "StompStereo"
			gearMap["StompB1"] = "StompB1"
			gearMap["StompB2"] = "StompB2"
			gearMap["StompB3"] = "StompB3"
			gearMap["LoopFxA"] = "LoopFxA"
			gearMap["LoopFxB"] = "LoopFxB"
			gearMap["LoopFxC"] = "LoopFxC"
			gearMap["RackA"] = "RackA"
			gearMap["RackB"] = "RackB"
			gearMap["RackC"] = "RackC"
			gearMap["RackDI"] = "RackDI"
			gearMap["RackMaster"] = "RackMaster"
		}
	}

	if err != nil {
		return nil
	}

	for _, target := range matches {

		targetFile, err := ioutil.ReadFile(target)

		if err != nil {
			return err
		}

		var targetPreset PresetXMLV5
		err = xml.Unmarshal(targetFile, &targetPreset)

		if err != nil {
			return err
		}

		for sourceType, targetType := range gearMap {

			switch sourceType {
			case "AmpA":
				if targetType == "" {
					targetPreset.AmpA = sourcePreset.AmpA
					if copyCabWithAmp {
						targetPreset.CabA = sourcePreset.CabA
					}
				} else {
					if err = copyAmp(fromAmpA(sourcePreset.AmpA), &targetPreset, targetType); err != nil {
						return err
					} else {
						if copyCabWithAmp {
							sourceCab := sourcePreset.CabA
							targetCab := "Cab" + targetType[len(targetType)-1:]
							if err = copyCab(fromCabA(sourceCab), &targetPreset, targetCab); err != nil {
								return err
							}
						}
					}
				}
			case "AmpB":
				if targetType == "" {
					targetPreset.AmpB = sourcePreset.AmpB
					if copyCabWithAmp {
						targetPreset.CabB = sourcePreset.CabB
					}
				} else {
					if err = copyAmp(fromAmpB(sourcePreset.AmpB), &targetPreset, targetType); err != nil {
						return err
					} else {
						if copyCabWithAmp {
							sourceCab := sourcePreset.CabB
							targetCab := "Cab" + targetType[len(targetType)-1:]
							if err = copyCab(fromCabB(sourceCab), &targetPreset, targetCab); err != nil {
								return err
							}
						}
					}
				}
			case "AmpC":
				if targetType == "" {
					targetPreset.AmpC = sourcePreset.AmpC
					if copyCabWithAmp {
						targetPreset.CabC = sourcePreset.CabC
					}
				} else {
					if err = copyAmp(fromAmpC(sourcePreset.AmpC), &targetPreset, targetType); err != nil {
						return err
					} else {
						if copyCabWithAmp {
							sourceCab := sourcePreset.CabC
							targetCab := "Cab" + targetType[len(targetType)-1:]
							if err = copyCab(fromCabC(sourceCab), &targetPreset, targetCab); err != nil {
								return err
							}
						}
					}
				}
			case "StompA1":
				sourceStomp := sourcePreset.StompA1
				if targetType == "" {
					targetType = "StompA1"
				}
				if removeFx {
					removeStomps(fromStompA1(sourceStomp), &targetPreset, targetType, targetSlot)
				} else {
					if appendFx || insertFx {
						if err = appendOrInsertStomps(fromStompA1(sourceStomp), &targetPreset, targetType, insertFx); err != nil {
							return err
						}
					} else {
						if err = replaceStomps(fromStompA1(sourceStomp), &targetPreset, targetType); err != nil {
							return err
						}
					}
				}
			case "StompA2":
				sourceStomp := sourcePreset.StompA2
				if targetType == "" {
					targetType = "StompA2"
				}
				if removeFx {
					removeStomps(fromStompA2(sourceStomp), &targetPreset, targetType, targetSlot)
				} else {
					if appendFx || insertFx {
						if err = appendOrInsertStomps(fromStompA2(sourceStomp), &targetPreset, targetType, insertFx); err != nil {
							return err
						}
					} else {
						if err = replaceStomps(fromStompA2(sourceStomp), &targetPreset, targetType); err != nil {
							return err
						}
					}
				}
			case "StompB1":
				sourceStomp := sourcePreset.StompB1
				if targetType == "" {
					targetType = "StompB1"
				}
				if removeFx {
					removeStomps(fromStompB1(sourceStomp), &targetPreset, targetType, targetSlot)
				} else {
					if appendFx || insertFx {
						if err = appendOrInsertStomps(fromStompB1(sourceStomp), &targetPreset, targetType, insertFx); err != nil {
							return err
						}
					} else {
						if err = replaceStomps(fromStompB1(sourceStomp), &targetPreset, targetType); err != nil {
							return err
						}
					}
				}
			case "StompB2":
				sourceStomp := sourcePreset.StompB2
				if targetType == "" {
					targetType = "StompB2"
				}
				if removeFx {
					removeStomps(fromStompB2(sourceStomp), &targetPreset, targetType, targetSlot)
				} else {
					if appendFx || insertFx {
						if err = appendOrInsertStomps(fromStompB2(sourceStomp), &targetPreset, targetType, insertFx); err != nil {
							return err
						}
					} else {
						if err = replaceStomps(fromStompB2(sourceStomp), &targetPreset, targetType); err != nil {
							return err
						}
					}
				}
			case "StompB3":
				sourceStomp := sourcePreset.StompB3
				if targetType == "" {
					targetType = "StompB3"
				}
				if removeFx {
					removeStomps(fromStompB3(sourceStomp), &targetPreset, targetType, targetSlot)
				} else {
					if appendFx || insertFx {
						if err = appendOrInsertStomps(fromStompB3(sourceStomp), &targetPreset, targetType, insertFx); err != nil {
							return err
						}
					} else {
						if err = replaceStomps(fromStompB3(sourceStomp), &targetPreset, targetType); err != nil {
							return err
						}
					}
				}
			case "StompStereo":
				sourceStomp := sourcePreset.StompStereo
				if targetType == "" {
					targetType = "StompStereo"
				}
				if removeFx {
					removeStomps(fromStompStereo(sourceStomp), &targetPreset, targetType, targetSlot)
				} else {
					if appendFx || insertFx {
						if err = appendOrInsertStomps(fromStompStereo(sourceStomp), &targetPreset, targetType, insertFx); err != nil {
							return err
						}
					} else {
						if err = replaceStomps(fromStompStereo(sourceStomp), &targetPreset, targetType); err != nil {
							return err
						}
					}
				}
			case "LoopFxA":
				sourceStomp := sourcePreset.LoopFxA
				if targetType == "" {
					targetType = "LoopFxA"
				}
				if removeFx {
					removeStomps(fromLoopFxA(sourceStomp), &targetPreset, targetType, targetSlot)
				} else {
					if appendFx || insertFx {
						if err = appendOrInsertStomps(fromLoopFxA(sourceStomp), &targetPreset, targetType, insertFx); err != nil {
							return err
						}
					} else {
						if err = replaceStomps(fromLoopFxA(sourceStomp), &targetPreset, targetType); err != nil {
							return err
						}
					}
				}
			case "LoopFxB":
				sourceStomp := sourcePreset.LoopFxB
				if targetType == "" {
					targetType = "LoopFxB"
				}
				if removeFx {
					removeStomps(fromLoopFxB(sourceStomp), &targetPreset, targetType, targetSlot)
				} else {
					if appendFx || insertFx {
						if err = appendOrInsertStomps(fromLoopFxB(sourceStomp), &targetPreset, targetType, insertFx); err != nil {
							return err
						}
					} else {
						if err = replaceStomps(fromLoopFxB(sourceStomp), &targetPreset, targetType); err != nil {
							return err
						}
					}
				}
			case "LoopFxC":
				sourceStomp := sourcePreset.LoopFxC
				if targetType == "" {
					targetType = "LoopFxC"
				}
				if removeFx {
					removeStomps(fromLoopFxC(sourceStomp), &targetPreset, targetType, targetSlot)
				} else {
					if appendFx || insertFx {
						if err = appendOrInsertStomps(fromLoopFxC(sourceStomp), &targetPreset, targetType, insertFx); err != nil {
							return err
						}
					} else {
						if err = replaceStomps(fromLoopFxC(sourceStomp), &targetPreset, targetType); err != nil {
							return err
						}
					}
				}
			case "CabA":
				sourceCab := sourcePreset.CabA
				if targetType == "" {
					targetPreset.CabA = sourceCab
				} else {
					if err = copyCab(fromCabA(sourceCab), &targetPreset, targetType); err != nil {
						return err
					}
				}
			case "CabB":
				sourceCab := sourcePreset.CabB
				if targetType == "" {
					targetPreset.CabB = sourceCab
				} else {
					if err = copyCab(fromCabB(sourceCab), &targetPreset, targetType); err != nil {
						return err
					}
				}
			case "CabC":
				sourceCab := sourcePreset.CabC
				if targetType == "" {
					targetPreset.CabC = sourceCab
				} else {
					if err = copyCab(fromCabC(sourceCab), &targetPreset, targetType); err != nil {
						return err
					}
				}
			case "Studio":
				targetPreset.Studio = sourcePreset.Studio
			case "RackA":
				sourceStomp := sourcePreset.RackA
				if targetType == "" {
					targetType = "RackA"
				}
				if removeFx {
					removeStomps(fromRackA(sourceStomp), &targetPreset, targetType, targetSlot)
				} else {
					if appendFx || insertFx {
						if err = appendOrInsertStomps(fromRackA(sourceStomp), &targetPreset, targetType, insertFx); err != nil {
							return err
						}
					} else {
						if err = replaceStomps(fromRackA(sourceStomp), &targetPreset, targetType); err != nil {
							return err
						}
					}
				}
			case "RackB":
				sourceStomp := sourcePreset.RackB
				if targetType == "" {
					targetType = "RackB"
				}
				if removeFx {
					removeStomps(fromRackB(sourceStomp), &targetPreset, targetType, targetSlot)
				} else {
					if appendFx || insertFx {
						if err = appendOrInsertStomps(fromRackB(sourceStomp), &targetPreset, targetType, insertFx); err != nil {
							return err
						}
					} else {
						if err = replaceStomps(fromRackB(sourceStomp), &targetPreset, targetType); err != nil {
							return err
						}
					}
				}
			case "RackC":
				sourceStomp := sourcePreset.RackC
				if targetType == "" {
					targetType = "RackC"
				}
				if removeFx {
					removeStomps(fromRackC(sourceStomp), &targetPreset, targetType, targetSlot)
				} else {
					if appendFx || insertFx {
						if err = appendOrInsertStomps(fromRackC(sourceStomp), &targetPreset, targetType, insertFx); err != nil {
							return err
						}
					} else {
						if err = replaceStomps(fromRackC(sourceStomp), &targetPreset, targetType); err != nil {
							return err
						}
					}
				}
			case "RackDI":
				sourceStomp := sourcePreset.RackDI
				if targetType == "" {
					targetType = "RackDI"
				}
				if removeFx {
					removeStomps(fromRackDI(sourceStomp), &targetPreset, targetType, targetSlot)
				} else {
					if appendFx || insertFx {
						if err = appendOrInsertStomps(fromRackDI(sourceStomp), &targetPreset, targetType, insertFx); err != nil {
							return err
						}
					} else {
						if err = replaceStomps(fromRackDI(sourceStomp), &targetPreset, targetType); err != nil {
							return err
						}
					}
				}
			case "RackMaster":
				sourceStomp := sourcePreset.RackMaster
				if targetType == "" {
					targetType = "RackMaster"
				}
				if removeFx {
					removeStomps(fromRackMaster(sourceStomp), &targetPreset, targetType, targetSlot)
				} else {
					if appendFx || insertFx {
						if err = appendOrInsertStomps(fromRackMaster(sourceStomp), &targetPreset, targetType, insertFx); err != nil {
							return err
						}
					} else {
						if err = replaceStomps(fromRackMaster(sourceStomp), &targetPreset, targetType); err != nil {
							return err
						}
					}
				}
			case "Tuner":
				targetPreset.Tuner = sourcePreset.Tuner
			default:
				return errors.New("gear copy not supported for " + context.Args[2])
			}

		}

		targetData, err := xml.MarshalIndent(targetPreset, "", "    ")

		err = ioutil.WriteFile(target, selfClose(append([]byte(XMLPrefix), append(targetData[:], '\n')...)), 0664)

	}

	if err != nil {
		return err
	}

	return nil
}

func copyAmp(sourceAmp GenericAmp, targetPreset *PresetXMLV5, targetSlot string) error {
	switch targetSlot {
	case "AmpA":
		targetPreset.AmpA = AmpA{
			Bypass:       sourceAmp.Bypass,
			Mute:         sourceAmp.Mute,
			OutputVolume: sourceAmp.OutputVolume,
			Model:        sourceAmp.Model,
			Amp:          sourceAmp.Amp,
		}
	case "AmpB":
		targetPreset.AmpB = AmpB{
			Bypass:       sourceAmp.Bypass,
			Mute:         sourceAmp.Mute,
			OutputVolume: sourceAmp.OutputVolume,
			Model:        sourceAmp.Model,
			Amp:          sourceAmp.Amp,
		}
	case "AmpC":
		targetPreset.AmpC = AmpC{
			Bypass:       sourceAmp.Bypass,
			Mute:         sourceAmp.Mute,
			OutputVolume: sourceAmp.OutputVolume,
			Model:        sourceAmp.Model,
			Amp:          sourceAmp.Amp,
		}
	default:
		return errors.New("incompatible slots")
	}
	return nil
}

func copyCab(sourceCab GenericCab, targetPreset *PresetXMLV5, targetSlot string) error {
	switch targetSlot {
	case "CabA":
		targetPreset.CabA = CabA{
			Bypass:                     sourceCab.Bypass,
			Mute:                       sourceCab.Mute,
			CabModel:                   sourceCab.CabModel,
			SpeakerModel0:              sourceCab.SpeakerModel0,
			SpeakerModel1:              sourceCab.SpeakerModel1,
			SpeakerModel2:              sourceCab.SpeakerModel2,
			SpeakerModel3:              sourceCab.SpeakerModel3,
			Studio_Mic1_Level:          sourceCab.Studio_Mic1_Level,
			Studio_Mic1_Pan:            sourceCab.Studio_Mic1_Pan,
			Studio_Mic1_Mute:           sourceCab.Studio_Mic1_Mute,
			Studio_Mic1_Angle:          sourceCab.Studio_Mic1_Angle,
			Studio_Mic1_Solo:           sourceCab.Studio_Mic1_Solo,
			Studio_Mic1_Phase:          sourceCab.Studio_Mic1_Phase,
			Studio_Mic2_Level:          sourceCab.Studio_Mic2_Level,
			Studio_Mic2_Pan:            sourceCab.Studio_Mic2_Pan,
			Studio_Mic2_Mute:           sourceCab.Studio_Mic2_Mute,
			Studio_Mic2_Angle:          sourceCab.Studio_Mic2_Angle,
			Studio_Mic2_Solo:           sourceCab.Studio_Mic2_Solo,
			Studio_Mic2_Phase:          sourceCab.Studio_Mic2_Phase,
			Studio_Room_Level:          sourceCab.Studio_Room_Level,
			Studio_Room_Width:          sourceCab.Studio_Room_Width,
			Studio_Room_Mute:           sourceCab.Studio_Room_Mute,
			Studio_Room_Solo:           sourceCab.Studio_Room_Solo,
			Studio_Room_Phase:          sourceCab.Studio_Room_Phase,
			Studio_Bus_Level:           sourceCab.Studio_Bus_Level,
			Studio_Bus_Pan:             sourceCab.Studio_Bus_Pan,
			Studio_Bus_Mute:            sourceCab.Studio_Bus_Mute,
			Studio_Bus_Solo:            sourceCab.Studio_Bus_Solo,
			Studio_Bus_Phase:           sourceCab.Studio_Bus_Phase,
			Studio_DI_Level:            sourceCab.Studio_DI_Level,
			Studio_DI_Pan:              sourceCab.Studio_DI_Pan,
			Studio_DI_Mute:             sourceCab.Studio_DI_Mute,
			Studio_DI_Solo:             sourceCab.Studio_DI_Solo,
			Studio_DI_Phase:            sourceCab.Studio_DI_Phase,
			DI_PhaseDelay:              sourceCab.DI_PhaseDelay,
			Studio_LeslieCab_Horn_VolL: sourceCab.Studio_LeslieCab_Horn_VolL,
			Studio_LeslieCab_Horn_VolR: sourceCab.Studio_LeslieCab_Horn_VolR,
			Studio_LeslieCab_Horn_PanL: sourceCab.Studio_LeslieCab_Horn_PanL,
			Studio_LeslieCab_Horn_PanR: sourceCab.Studio_LeslieCab_Horn_PanR,
			Studio_LeslieCab_Horn_Mute: sourceCab.Studio_LeslieCab_Horn_Mute,
			Studio_LeslieCab_Horn_Solo: sourceCab.Studio_LeslieCab_Horn_Solo,
			Studio_LeslieCab_Drum_VolL: sourceCab.Studio_LeslieCab_Drum_VolL,
			Studio_LeslieCab_Drum_VolR: sourceCab.Studio_LeslieCab_Drum_VolR,
			Studio_LeslieCab_Drum_PanL: sourceCab.Studio_LeslieCab_Drum_PanL,
			Studio_LeslieCab_Drum_PanR: sourceCab.Studio_LeslieCab_Drum_PanR,
			Studio_LeslieCab_Drum_Mute: sourceCab.Studio_LeslieCab_Drum_Mute,
			Studio_LeslieCab_Drum_Solo: sourceCab.Studio_LeslieCab_Drum_Solo,
			Cab:                        sourceCab.Cab,
		}
	case "CabB":
		targetPreset.CabB = CabB{
			Bypass:                     sourceCab.Bypass,
			Mute:                       sourceCab.Mute,
			CabModel:                   sourceCab.CabModel,
			SpeakerModel0:              sourceCab.SpeakerModel0,
			SpeakerModel1:              sourceCab.SpeakerModel1,
			SpeakerModel2:              sourceCab.SpeakerModel2,
			SpeakerModel3:              sourceCab.SpeakerModel3,
			Studio_Mic1_Level:          sourceCab.Studio_Mic1_Level,
			Studio_Mic1_Pan:            sourceCab.Studio_Mic1_Pan,
			Studio_Mic1_Mute:           sourceCab.Studio_Mic1_Mute,
			Studio_Mic1_Angle:          sourceCab.Studio_Mic1_Angle,
			Studio_Mic1_Solo:           sourceCab.Studio_Mic1_Solo,
			Studio_Mic1_Phase:          sourceCab.Studio_Mic1_Phase,
			Studio_Mic2_Level:          sourceCab.Studio_Mic2_Level,
			Studio_Mic2_Pan:            sourceCab.Studio_Mic2_Pan,
			Studio_Mic2_Mute:           sourceCab.Studio_Mic2_Mute,
			Studio_Mic2_Angle:          sourceCab.Studio_Mic2_Angle,
			Studio_Mic2_Solo:           sourceCab.Studio_Mic2_Solo,
			Studio_Mic2_Phase:          sourceCab.Studio_Mic2_Phase,
			Studio_Room_Level:          sourceCab.Studio_Room_Level,
			Studio_Room_Width:          sourceCab.Studio_Room_Width,
			Studio_Room_Mute:           sourceCab.Studio_Room_Mute,
			Studio_Room_Solo:           sourceCab.Studio_Room_Solo,
			Studio_Room_Phase:          sourceCab.Studio_Room_Phase,
			Studio_Bus_Level:           sourceCab.Studio_Bus_Level,
			Studio_Bus_Pan:             sourceCab.Studio_Bus_Pan,
			Studio_Bus_Mute:            sourceCab.Studio_Bus_Mute,
			Studio_Bus_Solo:            sourceCab.Studio_Bus_Solo,
			Studio_Bus_Phase:           sourceCab.Studio_Bus_Phase,
			Studio_DI_Level:            sourceCab.Studio_DI_Level,
			Studio_DI_Pan:              sourceCab.Studio_DI_Pan,
			Studio_DI_Mute:             sourceCab.Studio_DI_Mute,
			Studio_DI_Solo:             sourceCab.Studio_DI_Solo,
			Studio_DI_Phase:            sourceCab.Studio_DI_Phase,
			DI_PhaseDelay:              sourceCab.DI_PhaseDelay,
			Studio_LeslieCab_Horn_VolL: sourceCab.Studio_LeslieCab_Horn_VolL,
			Studio_LeslieCab_Horn_VolR: sourceCab.Studio_LeslieCab_Horn_VolR,
			Studio_LeslieCab_Horn_PanL: sourceCab.Studio_LeslieCab_Horn_PanL,
			Studio_LeslieCab_Horn_PanR: sourceCab.Studio_LeslieCab_Horn_PanR,
			Studio_LeslieCab_Horn_Mute: sourceCab.Studio_LeslieCab_Horn_Mute,
			Studio_LeslieCab_Horn_Solo: sourceCab.Studio_LeslieCab_Horn_Solo,
			Studio_LeslieCab_Drum_VolL: sourceCab.Studio_LeslieCab_Drum_VolL,
			Studio_LeslieCab_Drum_VolR: sourceCab.Studio_LeslieCab_Drum_VolR,
			Studio_LeslieCab_Drum_PanL: sourceCab.Studio_LeslieCab_Drum_PanL,
			Studio_LeslieCab_Drum_PanR: sourceCab.Studio_LeslieCab_Drum_PanR,
			Studio_LeslieCab_Drum_Mute: sourceCab.Studio_LeslieCab_Drum_Mute,
			Studio_LeslieCab_Drum_Solo: sourceCab.Studio_LeslieCab_Drum_Solo,
			Cab:                        sourceCab.Cab,
		}
	case "CabC":
		targetPreset.CabC = CabC{
			Bypass:                     sourceCab.Bypass,
			Mute:                       sourceCab.Mute,
			CabModel:                   sourceCab.CabModel,
			SpeakerModel0:              sourceCab.SpeakerModel0,
			SpeakerModel1:              sourceCab.SpeakerModel1,
			SpeakerModel2:              sourceCab.SpeakerModel2,
			SpeakerModel3:              sourceCab.SpeakerModel3,
			Studio_Mic1_Level:          sourceCab.Studio_Mic1_Level,
			Studio_Mic1_Pan:            sourceCab.Studio_Mic1_Pan,
			Studio_Mic1_Mute:           sourceCab.Studio_Mic1_Mute,
			Studio_Mic1_Angle:          sourceCab.Studio_Mic1_Angle,
			Studio_Mic1_Solo:           sourceCab.Studio_Mic1_Solo,
			Studio_Mic1_Phase:          sourceCab.Studio_Mic1_Phase,
			Studio_Mic2_Level:          sourceCab.Studio_Mic2_Level,
			Studio_Mic2_Pan:            sourceCab.Studio_Mic2_Pan,
			Studio_Mic2_Mute:           sourceCab.Studio_Mic2_Mute,
			Studio_Mic2_Angle:          sourceCab.Studio_Mic2_Angle,
			Studio_Mic2_Solo:           sourceCab.Studio_Mic2_Solo,
			Studio_Mic2_Phase:          sourceCab.Studio_Mic2_Phase,
			Studio_Room_Level:          sourceCab.Studio_Room_Level,
			Studio_Room_Width:          sourceCab.Studio_Room_Width,
			Studio_Room_Mute:           sourceCab.Studio_Room_Mute,
			Studio_Room_Solo:           sourceCab.Studio_Room_Solo,
			Studio_Room_Phase:          sourceCab.Studio_Room_Phase,
			Studio_Bus_Level:           sourceCab.Studio_Bus_Level,
			Studio_Bus_Pan:             sourceCab.Studio_Bus_Pan,
			Studio_Bus_Mute:            sourceCab.Studio_Bus_Mute,
			Studio_Bus_Solo:            sourceCab.Studio_Bus_Solo,
			Studio_Bus_Phase:           sourceCab.Studio_Bus_Phase,
			Studio_DI_Level:            sourceCab.Studio_DI_Level,
			Studio_DI_Pan:              sourceCab.Studio_DI_Pan,
			Studio_DI_Mute:             sourceCab.Studio_DI_Mute,
			Studio_DI_Solo:             sourceCab.Studio_DI_Solo,
			Studio_DI_Phase:            sourceCab.Studio_DI_Phase,
			DI_PhaseDelay:              sourceCab.DI_PhaseDelay,
			Studio_LeslieCab_Horn_VolL: sourceCab.Studio_LeslieCab_Horn_VolL,
			Studio_LeslieCab_Horn_VolR: sourceCab.Studio_LeslieCab_Horn_VolR,
			Studio_LeslieCab_Horn_PanL: sourceCab.Studio_LeslieCab_Horn_PanL,
			Studio_LeslieCab_Horn_PanR: sourceCab.Studio_LeslieCab_Horn_PanR,
			Studio_LeslieCab_Horn_Mute: sourceCab.Studio_LeslieCab_Horn_Mute,
			Studio_LeslieCab_Horn_Solo: sourceCab.Studio_LeslieCab_Horn_Solo,
			Studio_LeslieCab_Drum_VolL: sourceCab.Studio_LeslieCab_Drum_VolL,
			Studio_LeslieCab_Drum_VolR: sourceCab.Studio_LeslieCab_Drum_VolR,
			Studio_LeslieCab_Drum_PanL: sourceCab.Studio_LeslieCab_Drum_PanL,
			Studio_LeslieCab_Drum_PanR: sourceCab.Studio_LeslieCab_Drum_PanR,
			Studio_LeslieCab_Drum_Mute: sourceCab.Studio_LeslieCab_Drum_Mute,
			Studio_LeslieCab_Drum_Solo: sourceCab.Studio_LeslieCab_Drum_Solo,
			Cab:                        sourceCab.Cab,
		}
	default:
		return errors.New("incompatible slots")
	}
	return nil
}

func removeStomps(sourceStomps GenericStomp, targetPreset *PresetXMLV5, targetStompName string, targetSlot string) error {
	if targetSlot == "" || targetSlot == "Slot0" {
		sourceStomps.Stomp0 = EmptySlotGUID
		sourceStomps.Slot0 = Slot0{}
	}
	if targetSlot == "" || targetSlot == "Slot1" {
		sourceStomps.Stomp1 = EmptySlotGUID
		sourceStomps.Slot1 = Slot1{}
	}
	if targetSlot == "" || targetSlot == "Slot2" {
		sourceStomps.Stomp2 = EmptySlotGUID
		sourceStomps.Slot2 = Slot2{}
	}
	if targetSlot == "" || targetSlot == "Slot3" {
		sourceStomps.Stomp3 = EmptySlotGUID
		sourceStomps.Slot3 = Slot3{}
	}
	if targetSlot == "" || targetSlot == "Slot4" {
		sourceStomps.Stomp4 = EmptySlotGUID
		sourceStomps.Slot4 = Slot4{}
	}
	if targetSlot == "" || targetSlot == "Slot5" {
		sourceStomps.Stomp5 = EmptySlotGUID
		sourceStomps.Slot5 = Slot5{}
	}
	return replaceStomps(sourceStomps, targetPreset, targetStompName)
}

func appendOrInsertStomps(sourceStomps GenericStomp, targetPreset *PresetXMLV5, targetStompName string, insert bool) error {
	sourceSlotList := newSlotList()
	sourceGuidList := newGUIDList()
	for i := 0; i < sourceStomps.StompCount; i++ {
		var attrs []xml.Attr
		var guid string
		switch i {
		case 0:
			guid = sourceStomps.Stomp0
			attrs = sourceStomps.Slot0.Attrs
		case 1:
			guid = sourceStomps.Stomp1
			attrs = sourceStomps.Slot1.Attrs
		case 2:
			guid = sourceStomps.Stomp2
			attrs = sourceStomps.Slot2.Attrs
		case 3:
			guid = sourceStomps.Stomp3
			attrs = sourceStomps.Slot3.Attrs
		case 4:
			guid = sourceStomps.Stomp4
			attrs = sourceStomps.Slot4.Attrs
		case 5:
			guid = sourceStomps.Stomp5
			attrs = sourceStomps.Slot5.Attrs
		}
		sourceGuidList.append(guid)
		sourceSlotList.append(GenericSlot{attrs})
	}
	targetStomp := fromStompType(*targetPreset, targetStompName)
	slotList := newSlotList()
	guidList := newGUIDList()
	for i := 0; i < targetStomp.StompCount; i++ {
		var attrs []xml.Attr
		var guid string
		switch i {
		case 0:
			guid = targetStomp.Stomp0
			attrs = targetStomp.Slot0.Attrs
		case 1:
			guid = targetStomp.Stomp1
			attrs = targetStomp.Slot1.Attrs
		case 2:
			guid = targetStomp.Stomp2
			attrs = targetStomp.Slot2.Attrs
		case 3:
			guid = targetStomp.Stomp3
			attrs = targetStomp.Slot3.Attrs
		case 4:
			guid = targetStomp.Stomp4
			attrs = targetStomp.Slot4.Attrs
		case 5:
			guid = targetStomp.Stomp5
			attrs = targetStomp.Slot5.Attrs
		}
		guidList.append(guid)
		slotList.append(GenericSlot{Attrs: attrs})
	}
	if insert {
		for i := len(sourceGuidList.guids) - 1; i >= 0; i-- {
			guidList.insert(sourceGuidList.guids[i])
			slotList.insert(GenericSlot{Attrs: sourceSlotList.slots[i].Attrs})
		}
	} else {
		for i := 0; i < len(sourceGuidList.guids); i++ {
			guidList.append(sourceGuidList.guids[i])
			slotList.append(GenericSlot{Attrs: sourceSlotList.slots[i].Attrs})
		}
	}
	for i := 0; i < 6; i++ {
		if i < len(slotList.slots) {
			switch i {
			case 0:
				sourceStomps.Stomp0 = guidList.guids[i]
				sourceStomps.Slot0 = Slot0{Attrs: slotList.slots[i].Attrs}
			case 1:
				sourceStomps.Stomp1 = guidList.guids[i]
				sourceStomps.Slot1 = Slot1{Attrs: slotList.slots[i].Attrs}
			case 2:
				sourceStomps.Stomp2 = guidList.guids[i]
				sourceStomps.Slot2 = Slot2{Attrs: slotList.slots[i].Attrs}
			case 3:
				sourceStomps.Stomp3 = guidList.guids[i]
				sourceStomps.Slot3 = Slot3{Attrs: slotList.slots[i].Attrs}
			case 4:
				sourceStomps.Stomp4 = guidList.guids[i]
				sourceStomps.Slot4 = Slot4{Attrs: slotList.slots[i].Attrs}
			case 5:
				sourceStomps.Stomp5 = guidList.guids[i]
				sourceStomps.Slot5 = Slot5{Attrs: slotList.slots[i].Attrs}
			}
		} else {
			switch i {
			case 0:
				sourceStomps.Stomp0 = EmptySlotGUID
				sourceStomps.Slot0 = Slot0{}
			case 1:
				sourceStomps.Stomp1 = EmptySlotGUID
				sourceStomps.Slot1 = Slot1{}
			case 2:
				sourceStomps.Stomp2 = EmptySlotGUID
				sourceStomps.Slot2 = Slot2{}
			case 3:
				sourceStomps.Stomp3 = EmptySlotGUID
				sourceStomps.Slot3 = Slot3{}
			case 4:
				sourceStomps.Stomp4 = EmptySlotGUID
				sourceStomps.Slot4 = Slot4{}
			case 5:
				sourceStomps.Stomp5 = EmptySlotGUID
				sourceStomps.Slot5 = Slot5{}
			}
		}
	}
	return replaceStomps(sourceStomps, targetPreset, targetStompName)
}

func replaceStomps(sourceStomps GenericStomp, targetPreset *PresetXMLV5, targetStomps string) error {
	switch targetStomps {
	case "StompA1":
		targetPreset.StompA1 = StompA1{
			Bypass:       sourceStomps.Bypass,
			Mute:         sourceStomps.Mute,
			OutputVolume: sourceStomps.OutputVolume,
			Stomp0:       checkForEmptySlotGUID(sourceStomps.Stomp0, 0, sourceStomps.StompCount),
			Stomp1:       checkForEmptySlotGUID(sourceStomps.Stomp1, 1, sourceStomps.StompCount),
			Stomp2:       checkForEmptySlotGUID(sourceStomps.Stomp2, 2, sourceStomps.StompCount),
			Stomp3:       checkForEmptySlotGUID(sourceStomps.Stomp3, 3, sourceStomps.StompCount),
			Stomp4:       checkForEmptySlotGUID(sourceStomps.Stomp4, 4, sourceStomps.StompCount),
			Stomp5:       checkForEmptySlotGUID(sourceStomps.Stomp5, 5, sourceStomps.StompCount),
			Slot0:        checkForEmptySlot0(sourceStomps.Slot0, sourceStomps.StompCount),
			Slot1:        checkForEmptySlot1(sourceStomps.Slot1, sourceStomps.StompCount),
			Slot2:        checkForEmptySlot2(sourceStomps.Slot2, sourceStomps.StompCount),
			Slot3:        checkForEmptySlot3(sourceStomps.Slot3, sourceStomps.StompCount),
			Slot4:        checkForEmptySlot4(sourceStomps.Slot4, sourceStomps.StompCount),
			Slot5:        checkForEmptySlot5(sourceStomps.Slot5, sourceStomps.StompCount),
		}
	case "StompA2":
		targetPreset.StompA2 = StompA2{
			Bypass:       sourceStomps.Bypass,
			Mute:         sourceStomps.Mute,
			OutputVolume: sourceStomps.OutputVolume,
			Stomp0:       checkForEmptySlotGUID(sourceStomps.Stomp0, 0, sourceStomps.StompCount),
			Stomp1:       checkForEmptySlotGUID(sourceStomps.Stomp1, 1, sourceStomps.StompCount),
			Stomp2:       checkForEmptySlotGUID(sourceStomps.Stomp2, 2, sourceStomps.StompCount),
			Stomp3:       checkForEmptySlotGUID(sourceStomps.Stomp3, 3, sourceStomps.StompCount),
			Stomp4:       checkForEmptySlotGUID(sourceStomps.Stomp4, 4, sourceStomps.StompCount),
			Stomp5:       checkForEmptySlotGUID(sourceStomps.Stomp5, 5, sourceStomps.StompCount),
			Slot0:        checkForEmptySlot0(sourceStomps.Slot0, sourceStomps.StompCount),
			Slot1:        checkForEmptySlot1(sourceStomps.Slot1, sourceStomps.StompCount),
			Slot2:        checkForEmptySlot2(sourceStomps.Slot2, sourceStomps.StompCount),
			Slot3:        checkForEmptySlot3(sourceStomps.Slot3, sourceStomps.StompCount),
			Slot4:        checkForEmptySlot4(sourceStomps.Slot4, sourceStomps.StompCount),
			Slot5:        checkForEmptySlot5(sourceStomps.Slot5, sourceStomps.StompCount),
		}
	case "StompStereo":
		targetPreset.StompStereo = StompStereo{
			Bypass:       sourceStomps.Bypass,
			Mute:         sourceStomps.Mute,
			OutputVolume: sourceStomps.OutputVolume,
			Stomp0:       checkForEmptySlotGUID(sourceStomps.Stomp0, 0, sourceStomps.StompCount),
			Stomp1:       checkForEmptySlotGUID(sourceStomps.Stomp1, 1, sourceStomps.StompCount),
			Stomp2:       checkForEmptySlotGUID(sourceStomps.Stomp2, 2, sourceStomps.StompCount),
			Slot0:        checkForEmptySlot0(sourceStomps.Slot0, sourceStomps.StompCount),
			Slot1:        checkForEmptySlot1(sourceStomps.Slot1, sourceStomps.StompCount),
			Slot2:        checkForEmptySlot2(sourceStomps.Slot2, sourceStomps.StompCount),
		}
	case "StompB1":
		targetPreset.StompB1 = StompB1{
			Bypass:       sourceStomps.Bypass,
			Mute:         sourceStomps.Mute,
			OutputVolume: sourceStomps.OutputVolume,
			Stomp0:       checkForEmptySlotGUID(sourceStomps.Stomp0, 0, sourceStomps.StompCount),
			Stomp1:       checkForEmptySlotGUID(sourceStomps.Stomp1, 1, sourceStomps.StompCount),
			Stomp2:       checkForEmptySlotGUID(sourceStomps.Stomp2, 2, sourceStomps.StompCount),
			Stomp3:       checkForEmptySlotGUID(sourceStomps.Stomp3, 3, sourceStomps.StompCount),
			Stomp4:       checkForEmptySlotGUID(sourceStomps.Stomp4, 4, sourceStomps.StompCount),
			Stomp5:       checkForEmptySlotGUID(sourceStomps.Stomp5, 5, sourceStomps.StompCount),
			Slot0:        checkForEmptySlot0(sourceStomps.Slot0, sourceStomps.StompCount),
			Slot1:        checkForEmptySlot1(sourceStomps.Slot1, sourceStomps.StompCount),
			Slot2:        checkForEmptySlot2(sourceStomps.Slot2, sourceStomps.StompCount),
			Slot3:        checkForEmptySlot3(sourceStomps.Slot3, sourceStomps.StompCount),
			Slot4:        checkForEmptySlot4(sourceStomps.Slot4, sourceStomps.StompCount),
			Slot5:        checkForEmptySlot5(sourceStomps.Slot5, sourceStomps.StompCount),
		}
	case "StompB2":
		targetPreset.StompB2 = StompB2{
			Bypass:       sourceStomps.Bypass,
			Mute:         sourceStomps.Mute,
			OutputVolume: sourceStomps.OutputVolume,
			Stomp0:       checkForEmptySlotGUID(sourceStomps.Stomp0, 0, sourceStomps.StompCount),
			Stomp1:       checkForEmptySlotGUID(sourceStomps.Stomp1, 1, sourceStomps.StompCount),
			Stomp2:       checkForEmptySlotGUID(sourceStomps.Stomp2, 2, sourceStomps.StompCount),
			Stomp3:       checkForEmptySlotGUID(sourceStomps.Stomp3, 3, sourceStomps.StompCount),
			Stomp4:       checkForEmptySlotGUID(sourceStomps.Stomp4, 4, sourceStomps.StompCount),
			Stomp5:       checkForEmptySlotGUID(sourceStomps.Stomp5, 5, sourceStomps.StompCount),
			Slot0:        checkForEmptySlot0(sourceStomps.Slot0, sourceStomps.StompCount),
			Slot1:        checkForEmptySlot1(sourceStomps.Slot1, sourceStomps.StompCount),
			Slot2:        checkForEmptySlot2(sourceStomps.Slot2, sourceStomps.StompCount),
			Slot3:        checkForEmptySlot3(sourceStomps.Slot3, sourceStomps.StompCount),
			Slot4:        checkForEmptySlot4(sourceStomps.Slot4, sourceStomps.StompCount),
			Slot5:        checkForEmptySlot5(sourceStomps.Slot5, sourceStomps.StompCount),
		}
	case "StompB3":
		targetPreset.StompB3 = StompB3{
			Bypass:       sourceStomps.Bypass,
			Mute:         sourceStomps.Mute,
			OutputVolume: sourceStomps.OutputVolume,
			Stomp0:       checkForEmptySlotGUID(sourceStomps.Stomp0, 0, sourceStomps.StompCount),
			Stomp1:       checkForEmptySlotGUID(sourceStomps.Stomp1, 1, sourceStomps.StompCount),
			Stomp2:       checkForEmptySlotGUID(sourceStomps.Stomp2, 2, sourceStomps.StompCount),
			Stomp3:       checkForEmptySlotGUID(sourceStomps.Stomp3, 3, sourceStomps.StompCount),
			Stomp4:       checkForEmptySlotGUID(sourceStomps.Stomp4, 4, sourceStomps.StompCount),
			Stomp5:       checkForEmptySlotGUID(sourceStomps.Stomp5, 5, sourceStomps.StompCount),
			Slot0:        checkForEmptySlot0(sourceStomps.Slot0, sourceStomps.StompCount),
			Slot1:        checkForEmptySlot1(sourceStomps.Slot1, sourceStomps.StompCount),
			Slot2:        checkForEmptySlot2(sourceStomps.Slot2, sourceStomps.StompCount),
			Slot3:        checkForEmptySlot3(sourceStomps.Slot3, sourceStomps.StompCount),
			Slot4:        checkForEmptySlot4(sourceStomps.Slot4, sourceStomps.StompCount),
			Slot5:        checkForEmptySlot5(sourceStomps.Slot5, sourceStomps.StompCount),
		}
	case "LoopFxA":
		targetPreset.LoopFxA = LoopFxA{
			Bypass:       sourceStomps.Bypass,
			Mute:         sourceStomps.Mute,
			OutputVolume: sourceStomps.OutputVolume,
			Stomp0:       checkForEmptySlotGUID(sourceStomps.Stomp0, 0, sourceStomps.StompCount),
			Stomp1:       checkForEmptySlotGUID(sourceStomps.Stomp1, 1, sourceStomps.StompCount),
			Stomp2:       checkForEmptySlotGUID(sourceStomps.Stomp2, 2, sourceStomps.StompCount),
			Stomp3:       checkForEmptySlotGUID(sourceStomps.Stomp3, 3, sourceStomps.StompCount),
			Slot0:        checkForEmptySlot0(sourceStomps.Slot0, sourceStomps.StompCount),
			Slot1:        checkForEmptySlot1(sourceStomps.Slot1, sourceStomps.StompCount),
			Slot2:        checkForEmptySlot2(sourceStomps.Slot2, sourceStomps.StompCount),
			Slot3:        checkForEmptySlot3(sourceStomps.Slot3, sourceStomps.StompCount),
		}
	case "LoopFxB":
		targetPreset.LoopFxB = LoopFxB{
			Bypass:       sourceStomps.Bypass,
			Mute:         sourceStomps.Mute,
			OutputVolume: sourceStomps.OutputVolume,
			Stomp0:       checkForEmptySlotGUID(sourceStomps.Stomp0, 0, sourceStomps.StompCount),
			Stomp1:       checkForEmptySlotGUID(sourceStomps.Stomp1, 1, sourceStomps.StompCount),
			Stomp2:       checkForEmptySlotGUID(sourceStomps.Stomp2, 2, sourceStomps.StompCount),
			Stomp3:       checkForEmptySlotGUID(sourceStomps.Stomp3, 3, sourceStomps.StompCount),
			Slot0:        checkForEmptySlot0(sourceStomps.Slot0, sourceStomps.StompCount),
			Slot1:        checkForEmptySlot1(sourceStomps.Slot1, sourceStomps.StompCount),
			Slot2:        checkForEmptySlot2(sourceStomps.Slot2, sourceStomps.StompCount),
			Slot3:        checkForEmptySlot3(sourceStomps.Slot3, sourceStomps.StompCount),
		}
	case "LoopFxC":
		targetPreset.LoopFxC = LoopFxC{
			Bypass:       sourceStomps.Bypass,
			Mute:         sourceStomps.Mute,
			OutputVolume: sourceStomps.OutputVolume,
			Stomp0:       checkForEmptySlotGUID(sourceStomps.Stomp0, 0, sourceStomps.StompCount),
			Stomp1:       checkForEmptySlotGUID(sourceStomps.Stomp1, 1, sourceStomps.StompCount),
			Stomp2:       checkForEmptySlotGUID(sourceStomps.Stomp2, 2, sourceStomps.StompCount),
			Stomp3:       checkForEmptySlotGUID(sourceStomps.Stomp3, 3, sourceStomps.StompCount),
			Slot0:        checkForEmptySlot0(sourceStomps.Slot0, sourceStomps.StompCount),
			Slot1:        checkForEmptySlot1(sourceStomps.Slot1, sourceStomps.StompCount),
			Slot2:        checkForEmptySlot2(sourceStomps.Slot2, sourceStomps.StompCount),
			Slot3:        checkForEmptySlot3(sourceStomps.Slot3, sourceStomps.StompCount),
		}
	case "RackA":
		targetPreset.RackA = RackA{
			Bypass:       sourceStomps.Bypass,
			Mute:         sourceStomps.Mute,
			OutputVolume: sourceStomps.OutputVolume,
			Stomp0:       checkForEmptySlotGUID(sourceStomps.Stomp0, 0, sourceStomps.StompCount),
			Stomp1:       checkForEmptySlotGUID(sourceStomps.Stomp1, 1, sourceStomps.StompCount),
			Slot0:        checkForEmptySlot0(sourceStomps.Slot0, sourceStomps.StompCount),
			Slot1:        checkForEmptySlot1(sourceStomps.Slot1, sourceStomps.StompCount),
		}
	case "RackB":
		targetPreset.RackB = RackB{
			Bypass:       sourceStomps.Bypass,
			Mute:         sourceStomps.Mute,
			OutputVolume: sourceStomps.OutputVolume,
			Stomp0:       checkForEmptySlotGUID(sourceStomps.Stomp0, 0, sourceStomps.StompCount),
			Stomp1:       checkForEmptySlotGUID(sourceStomps.Stomp1, 1, sourceStomps.StompCount),
			Slot0:        checkForEmptySlot0(sourceStomps.Slot0, sourceStomps.StompCount),
			Slot1:        checkForEmptySlot1(sourceStomps.Slot1, sourceStomps.StompCount),
		}
	case "RackC":
		targetPreset.RackC = RackC{
			Bypass:       sourceStomps.Bypass,
			Mute:         sourceStomps.Mute,
			OutputVolume: sourceStomps.OutputVolume,
			Stomp0:       checkForEmptySlotGUID(sourceStomps.Stomp0, 0, sourceStomps.StompCount),
			Stomp1:       checkForEmptySlotGUID(sourceStomps.Stomp1, 1, sourceStomps.StompCount),
			Slot0:        checkForEmptySlot0(sourceStomps.Slot0, sourceStomps.StompCount),
			Slot1:        checkForEmptySlot1(sourceStomps.Slot1, sourceStomps.StompCount),
		}
	case "RackDI":
		targetPreset.RackDI = RackDI{
			Bypass:       sourceStomps.Bypass,
			Mute:         sourceStomps.Mute,
			OutputVolume: sourceStomps.OutputVolume,
			Stomp0:       checkForEmptySlotGUID(sourceStomps.Stomp0, 0, sourceStomps.StompCount),
			Stomp1:       checkForEmptySlotGUID(sourceStomps.Stomp1, 1, sourceStomps.StompCount),
			Slot0:        checkForEmptySlot0(sourceStomps.Slot0, sourceStomps.StompCount),
			Slot1:        checkForEmptySlot1(sourceStomps.Slot1, sourceStomps.StompCount),
		}
	case "RackMaster":
		targetPreset.RackMaster = RackMaster{
			Bypass:       sourceStomps.Bypass,
			Mute:         sourceStomps.Mute,
			OutputVolume: sourceStomps.OutputVolume,
			Stomp0:       checkForEmptySlotGUID(sourceStomps.Stomp0, 0, sourceStomps.StompCount),
			Stomp1:       checkForEmptySlotGUID(sourceStomps.Stomp1, 1, sourceStomps.StompCount),
			Stomp2:       checkForEmptySlotGUID(sourceStomps.Stomp2, 2, sourceStomps.StompCount),
			Stomp3:       checkForEmptySlotGUID(sourceStomps.Stomp3, 3, sourceStomps.StompCount),
			Stomp4:       checkForEmptySlotGUID(sourceStomps.Stomp4, 4, sourceStomps.StompCount),
			Stomp5:       checkForEmptySlotGUID(sourceStomps.Stomp5, 5, sourceStomps.StompCount),
			Slot0:        checkForEmptySlot0(sourceStomps.Slot0, sourceStomps.StompCount),
			Slot1:        checkForEmptySlot1(sourceStomps.Slot1, sourceStomps.StompCount),
			Slot2:        checkForEmptySlot2(sourceStomps.Slot2, sourceStomps.StompCount),
			Slot3:        checkForEmptySlot3(sourceStomps.Slot3, sourceStomps.StompCount),
			Slot4:        checkForEmptySlot4(sourceStomps.Slot4, sourceStomps.StompCount),
			Slot5:        checkForEmptySlot5(sourceStomps.Slot5, sourceStomps.StompCount),
		}
	default:
		return errors.New("incompatible slots")
	}
	return nil
}

func checkForEmptySlotGUID(guid string, index int, total int) string {
	if index >= total {
		return EmptySlotGUID
	}
	return guid
}

func checkForEmptySlot0(slot Slot0, total int) Slot0 {
	if 0 >= total {
		return Slot0{}
	}
	return slot
}

func checkForEmptySlot1(slot Slot1, total int) Slot1 {
	if 1 >= total {
		return Slot1{}
	}
	return slot
}

func checkForEmptySlot2(slot Slot2, total int) Slot2 {
	if 2 >= total {
		return Slot2{}
	}
	return slot
}

func checkForEmptySlot3(slot Slot3, total int) Slot3 {
	if 3 >= total {
		return Slot3{}
	}
	return slot
}

func checkForEmptySlot4(slot Slot4, total int) Slot4 {
	if 4 >= total {
		return Slot4{}
	}
	return slot
}

func checkForEmptySlot5(slot Slot5, total int) Slot5 {
	if 5 >= total {
		return Slot5{}
	}
	return slot
}

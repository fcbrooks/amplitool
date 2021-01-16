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
	"github.com/google/uuid"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

const PresetsFolder = "Presets"
const PresetExtension = ".at5p"
const PresetExtension4 = ".at4p"
const XMLPrefix = "<?xml version=\"1.0\" ?>\n"
const EmptySlotGUID = "773b8ea7-b54a-4a3c-99df-ffbbf6d29271"

type PresetXMLFormatOnly struct {
	XMLName xml.Name `xml:"Preset"`
	Format  string   `xml:",attr"`
}

type PresetXMLRootOnlyV5 struct {
	XMLName       xml.Name `xml:"Preset"`
	Version       int      `xml:",attr"`
	Format        string   `xml:",attr"`
	GUID          string   `xml:",attr"`
	PresetBPM     int      `xml:",attr"`
	ProgramChange int      `xml:",attr"`
	Inner         []byte   `xml:",innerxml"`
}

type PresetXMLRootOnlyV4 struct {
	XMLName       xml.Name `xml:"Preset"`
	Version       float32  `xml:",attr"`
	Format        string   `xml:",attr"`
	GUID          string   `xml:",attr"`
	PresetBPM     float32  `xml:",attr"`
	ProgramChange int      `xml:",attr"`
	Inner         []byte   `xml:",innerxml"`
}

type PresetXMLV5Raw struct {
	XMLName         xml.Name
	Attrs           []xml.Attr         `xml:",any,attr"`
	Chain           ChainRaw           `xml:""`
	Input           InputRaw           `xml:""`
	Tuner           TunerRaw           `xml:""`
	StompA1         StompA1Raw         `xml:""`
	StompA2         StompA2Raw         `xml:""`
	StompStereo     StompStereoRaw     `xml:""`
	StompB1         StompB1Raw         `xml:""`
	StompB2         StompB2Raw         `xml:""`
	StompB3         StompB3Raw         `xml:""`
	AmpA            AmpARaw            `xml:""`
	AmpB            AmpBRaw            `xml:""`
	AmpC            AmpCRaw            `xml:""`
	LoopFxA         LoopFxARaw         `xml:""`
	LoopFxB         LoopFxBRaw         `xml:""`
	LoopFxC         LoopFxCRaw         `xml:""`
	CabA            CabARaw            `xml:""`
	CabB            CabBRaw            `xml:""`
	CabC            CabCRaw            `xml:""`
	Studio          StudioRaw          `xml:""`
	RackA           RackARaw           `xml:""`
	RackB           RackBRaw           `xml:""`
	RackC           RackCRaw           `xml:""`
	RackDI          RackDIRaw          `xml:""`
	RackMaster      RackMasterRaw      `xml:""`
	Output          OutputRaw          `xml:""`
	MidiAssignments MidiAssignmentsRaw `xml:""`
	MetaInfo        MetaInfoRaw        `xml:""`
}

type PresetXMLV5 struct {
	XMLName         xml.Name        `xml:"Preset"`
	Version         int             `xml:",attr"`
	Format          string          `xml:",attr"`
	GUID            string          `xml:",attr"`
	PresetBPM       int             `xml:",attr"`
	ProgramChange   int             `xml:",attr"`
	Chain           Chain           `xml:""`
	Input           Input           `xml:""`
	Tuner           Tuner           `xml:""`
	StompA1         StompA1         `xml:""`
	StompA2         StompA2         `xml:""`
	StompStereo     StompStereo     `xml:""`
	StompB1         StompB1         `xml:""`
	StompB2         StompB2         `xml:""`
	StompB3         StompB3         `xml:""`
	AmpA            AmpA            `xml:""`
	AmpB            AmpB            `xml:""`
	AmpC            AmpC            `xml:""`
	LoopFxA         LoopFxA         `xml:""`
	LoopFxB         LoopFxB         `xml:""`
	LoopFxC         LoopFxC         `xml:""`
	CabA            CabA            `xml:""`
	CabB            CabB            `xml:""`
	CabC            CabC            `xml:""`
	Studio          Studio          `xml:""`
	RackA           RackA           `xml:""`
	RackB           RackB           `xml:""`
	RackC           RackC           `xml:""`
	RackDI          RackDI          `xml:""`
	RackMaster      RackMaster      `xml:""`
	Output          Output          `xml:""`
	MidiAssignments MidiAssignments `xml:""`
	MetaInfo        MetaInfo        `xml:""`
}

type CabAttrs struct {
	XMLName xml.Name  `xml:"Preset"`
	CabA    CabAAttrs `xml:""`
	CabB    CabBAttrs `xml:""`
	CabC    CabCAttrs `xml:""`
}

type ChainRaw struct {
	XMLName xml.Name   `xml:"Chain"`
	Attrs   []xml.Attr `xml:",any,attr"`
}

type Chain struct {
	XMLName          xml.Name `xml:"Chain"`
	Preset           string   `xml:",attr"`
	MonoChainDualCab int      `xml:",attr"`
	DIBeforeAmp      int      `xml:",attr"`
}

type InputRaw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
}

type Input struct {
	XMLName xml.Name `xml:"Input"`
	Input   int      `xml:",attr"`
}

type TunerRaw struct {
	XMLName xml.Name
	Attrs   []xml.Attr    `xml:",any,attr"`
	Tuner   TunerTunerRaw `xml:""`
}

type TunerTunerRaw struct {
	XMLName xml.Name   `xml:"Tuner"`
	Attrs   []xml.Attr `xml:",any,attr"`
}

type Tuner struct {
	XMLName      xml.Name `xml:"Tuner"`
	Bypass       int      `xml:",attr"`
	Mute         int      `xml:",attr"`
	OutputVolume int      `xml:",attr"`
	TunerType    string   `xml:",attr"`
	Inner        []byte   `xml:",innerxml"`
}

type Slot0 struct {
	XMLName xml.Name   `xml:"Slot0"`
	Attrs   []xml.Attr `xml:",any,attr"`
}

type Slot1 struct {
	XMLName xml.Name   `xml:"Slot1"`
	Attrs   []xml.Attr `xml:",any,attr"`
}

type Slot2 struct {
	XMLName xml.Name   `xml:"Slot2"`
	Attrs   []xml.Attr `xml:",any,attr"`
}

type Slot3 struct {
	XMLName xml.Name   `xml:"Slot3"`
	Attrs   []xml.Attr `xml:",any,attr"`
}

type Slot4 struct {
	XMLName xml.Name   `xml:"Slot4"`
	Attrs   []xml.Attr `xml:",any,attr"`
}

type Slot5 struct {
	XMLName xml.Name   `xml:"Slot5"`
	Attrs   []xml.Attr `xml:",any,attr"`
}

type GenericSlot struct {
	Attrs []xml.Attr
}

type GenericStomp struct {
	Bypass       int
	Mute         int
	OutputVolume int
	Stomp0       string
	Stomp1       string
	Stomp2       string
	Stomp3       string
	Stomp4       string
	Stomp5       string
	Slot0        Slot0
	Slot1        Slot1
	Slot2        Slot2
	Slot3        Slot3
	Slot4        Slot4
	Slot5        Slot5
	StompCount   int
}

func fromStompType(preset PresetXMLV5, name string) GenericStomp {
	var stomp GenericStomp
	switch name {
	case "StompA1":
		stomp = fromStompA1(preset.StompA1)
	case "StompA2":
		stomp = fromStompA2(preset.StompA2)
	case "StompStereo":
		stomp = fromStompStereo(preset.StompStereo)
	case "StompB1":
		stomp = fromStompB1(preset.StompB1)
	case "StompB2":
		stomp = fromStompB2(preset.StompB2)
	case "StompB3":
		stomp = fromStompB3(preset.StompB3)
	case "LoopFxA":
		stomp = fromLoopFxA(preset.LoopFxA)
	case "LoopFxB":
		stomp = fromLoopFxB(preset.LoopFxB)
	case "LoopFxC":
		stomp = fromLoopFxC(preset.LoopFxC)
	case "RackA":
		stomp = fromRackA(preset.RackA)
	case "RackB":
		stomp = fromRackB(preset.RackB)
	case "RackC":
		stomp = fromRackC(preset.RackC)
	case "RackDI":
		stomp = fromRackDI(preset.RackDI)
	case "RackMaster":
		stomp = fromRackMaster(preset.RackMaster)
	}
	return stomp
}

func fromStompA1(stomp StompA1) GenericStomp {
	return GenericStomp{
		Bypass:       stomp.Bypass,
		Mute:         stomp.Mute,
		OutputVolume: stomp.OutputVolume,
		Stomp0:       stomp.Stomp0,
		Stomp1:       stomp.Stomp1,
		Stomp2:       stomp.Stomp2,
		Stomp3:       stomp.Stomp3,
		Stomp4:       stomp.Stomp4,
		Stomp5:       stomp.Stomp5,
		Slot0:        stomp.Slot0,
		Slot1:        stomp.Slot1,
		Slot2:        stomp.Slot2,
		Slot3:        stomp.Slot3,
		Slot4:        stomp.Slot4,
		Slot5:        stomp.Slot5,
		StompCount:   6,
	}
}

type StompA1Raw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Slot0   Slot0      `xml:""`
	Slot1   Slot1      `xml:""`
	Slot2   Slot2      `xml:""`
	Slot3   Slot3      `xml:""`
	Slot4   Slot4      `xml:""`
	Slot5   Slot5      `xml:""`
}

type StompA1 struct {
	XMLName      xml.Name `xml:"StompA1"`
	Bypass       int      `xml:",attr"`
	Mute         int      `xml:",attr"`
	OutputVolume int      `xml:",attr"`
	Stomp0       string   `xml:",attr"`
	Stomp1       string   `xml:",attr"`
	Stomp2       string   `xml:",attr"`
	Stomp3       string   `xml:",attr"`
	Stomp4       string   `xml:",attr"`
	Stomp5       string   `xml:",attr"`
	Slot0        Slot0    `xml:""`
	Slot1        Slot1    `xml:""`
	Slot2        Slot2    `xml:""`
	Slot3        Slot3    `xml:""`
	Slot4        Slot4    `xml:""`
	Slot5        Slot5    `xml:""`
}

func fromStompA2(stomp StompA2) GenericStomp {
	return GenericStomp{
		Bypass:       stomp.Bypass,
		Mute:         stomp.Mute,
		OutputVolume: stomp.OutputVolume,
		Stomp0:       stomp.Stomp0,
		Stomp1:       stomp.Stomp1,
		Stomp2:       stomp.Stomp2,
		Stomp3:       stomp.Stomp3,
		Stomp4:       stomp.Stomp4,
		Stomp5:       stomp.Stomp5,
		Slot0:        stomp.Slot0,
		Slot1:        stomp.Slot1,
		Slot2:        stomp.Slot2,
		Slot3:        stomp.Slot3,
		Slot4:        stomp.Slot4,
		Slot5:        stomp.Slot5,
		StompCount:   6,
	}
}

type StompA2Raw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Slot0   Slot0      `xml:""`
	Slot1   Slot1      `xml:""`
	Slot2   Slot2      `xml:""`
	Slot3   Slot3      `xml:""`
	Slot4   Slot4      `xml:""`
	Slot5   Slot5      `xml:""`
}

type StompA2 struct {
	XMLName      xml.Name `xml:"StompA2"`
	Bypass       int      `xml:",attr"`
	Mute         int      `xml:",attr"`
	OutputVolume int      `xml:",attr"`
	Stomp0       string   `xml:",attr"`
	Stomp1       string   `xml:",attr"`
	Stomp2       string   `xml:",attr"`
	Stomp3       string   `xml:",attr"`
	Stomp4       string   `xml:",attr"`
	Stomp5       string   `xml:",attr"`
	Slot0        Slot0    `xml:""`
	Slot1        Slot1    `xml:""`
	Slot2        Slot2    `xml:""`
	Slot3        Slot3    `xml:""`
	Slot4        Slot4    `xml:""`
	Slot5        Slot5    `xml:""`
}

func fromStompStereo(stomp StompStereo) GenericStomp {
	return GenericStomp{
		Bypass:       stomp.Bypass,
		Mute:         stomp.Mute,
		OutputVolume: stomp.OutputVolume,
		Stomp0:       stomp.Stomp0,
		Stomp1:       stomp.Stomp1,
		Stomp2:       stomp.Stomp2,
		Stomp3:       "",
		Stomp4:       "",
		Stomp5:       "",
		Slot0:        stomp.Slot0,
		Slot1:        stomp.Slot1,
		Slot2:        stomp.Slot2,
		Slot3:        Slot3{},
		Slot4:        Slot4{},
		Slot5:        Slot5{},
		StompCount:   3,
	}
}

type StompStereoRaw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Slot0   Slot0      `xml:""`
	Slot1   Slot1      `xml:""`
	Slot2   Slot2      `xml:""`
}

type StompStereo struct {
	XMLName      xml.Name `xml:"StompStereo"`
	Bypass       int      `xml:",attr"`
	Mute         int      `xml:",attr"`
	OutputVolume int      `xml:",attr"`
	Stomp0       string   `xml:",attr"`
	Stomp1       string   `xml:",attr"`
	Stomp2       string   `xml:",attr"`
	Slot0        Slot0    `xml:""`
	Slot1        Slot1    `xml:""`
	Slot2        Slot2    `xml:""`
}

func fromStompB1(stomp StompB1) GenericStomp {
	return GenericStomp{
		Bypass:       stomp.Bypass,
		Mute:         stomp.Mute,
		OutputVolume: stomp.OutputVolume,
		Stomp0:       stomp.Stomp0,
		Stomp1:       stomp.Stomp1,
		Stomp2:       stomp.Stomp2,
		Stomp3:       stomp.Stomp3,
		Stomp4:       stomp.Stomp4,
		Stomp5:       stomp.Stomp5,
		Slot0:        stomp.Slot0,
		Slot1:        stomp.Slot1,
		Slot2:        stomp.Slot2,
		Slot3:        stomp.Slot3,
		Slot4:        stomp.Slot4,
		Slot5:        stomp.Slot5,
		StompCount:   6,
	}
}

type StompB1Raw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Slot0   Slot0      `xml:""`
	Slot1   Slot1      `xml:""`
	Slot2   Slot2      `xml:""`
	Slot3   Slot3      `xml:""`
	Slot4   Slot4      `xml:""`
	Slot5   Slot5      `xml:""`
}

type StompB1 struct {
	XMLName      xml.Name `xml:"StompB1"`
	Bypass       int      `xml:",attr"`
	Mute         int      `xml:",attr"`
	OutputVolume int      `xml:",attr"`
	Stomp0       string   `xml:",attr"`
	Stomp1       string   `xml:",attr"`
	Stomp2       string   `xml:",attr"`
	Stomp3       string   `xml:",attr"`
	Stomp4       string   `xml:",attr"`
	Stomp5       string   `xml:",attr"`
	Slot0        Slot0    `xml:""`
	Slot1        Slot1    `xml:""`
	Slot2        Slot2    `xml:""`
	Slot3        Slot3    `xml:""`
	Slot4        Slot4    `xml:""`
	Slot5        Slot5    `xml:""`
}

func fromStompB2(stomp StompB2) GenericStomp {
	return GenericStomp{
		Bypass:       stomp.Bypass,
		Mute:         stomp.Mute,
		OutputVolume: stomp.OutputVolume,
		Stomp0:       stomp.Stomp0,
		Stomp1:       stomp.Stomp1,
		Stomp2:       stomp.Stomp2,
		Stomp3:       stomp.Stomp3,
		Stomp4:       stomp.Stomp4,
		Stomp5:       stomp.Stomp5,
		Slot0:        stomp.Slot0,
		Slot1:        stomp.Slot1,
		Slot2:        stomp.Slot2,
		Slot3:        stomp.Slot3,
		Slot4:        stomp.Slot4,
		Slot5:        stomp.Slot5,
		StompCount:   6,
	}
}

type StompB2Raw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Slot0   Slot0      `xml:""`
	Slot1   Slot1      `xml:""`
	Slot2   Slot2      `xml:""`
	Slot3   Slot3      `xml:""`
	Slot4   Slot4      `xml:""`
	Slot5   Slot5      `xml:""`
}

type StompB2 struct {
	XMLName      xml.Name `xml:"StompB2"`
	Bypass       int      `xml:",attr"`
	Mute         int      `xml:",attr"`
	OutputVolume int      `xml:",attr"`
	Stomp0       string   `xml:",attr"`
	Stomp1       string   `xml:",attr"`
	Stomp2       string   `xml:",attr"`
	Stomp3       string   `xml:",attr"`
	Stomp4       string   `xml:",attr"`
	Stomp5       string   `xml:",attr"`
	Slot0        Slot0    `xml:""`
	Slot1        Slot1    `xml:""`
	Slot2        Slot2    `xml:""`
	Slot3        Slot3    `xml:""`
	Slot4        Slot4    `xml:""`
	Slot5        Slot5    `xml:""`
}

func fromStompB3(stomp StompB3) GenericStomp {
	return GenericStomp{
		Bypass:       stomp.Bypass,
		Mute:         stomp.Mute,
		OutputVolume: stomp.OutputVolume,
		Stomp0:       stomp.Stomp0,
		Stomp1:       stomp.Stomp1,
		Stomp2:       stomp.Stomp2,
		Stomp3:       stomp.Stomp3,
		Stomp4:       stomp.Stomp4,
		Stomp5:       stomp.Stomp5,
		Slot0:        stomp.Slot0,
		Slot1:        stomp.Slot1,
		Slot2:        stomp.Slot2,
		Slot3:        stomp.Slot3,
		Slot4:        stomp.Slot4,
		Slot5:        stomp.Slot5,
		StompCount:   6,
	}
}

type StompB3Raw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Slot0   Slot0      `xml:""`
	Slot1   Slot1      `xml:""`
	Slot2   Slot2      `xml:""`
	Slot3   Slot3      `xml:""`
	Slot4   Slot4      `xml:""`
	Slot5   Slot5      `xml:""`
}

type StompB3 struct {
	XMLName      xml.Name `xml:"StompB3"`
	Bypass       int      `xml:",attr"`
	Mute         int      `xml:",attr"`
	OutputVolume int      `xml:",attr"`
	Stomp0       string   `xml:",attr"`
	Stomp1       string   `xml:",attr"`
	Stomp2       string   `xml:",attr"`
	Stomp3       string   `xml:",attr"`
	Stomp4       string   `xml:",attr"`
	Stomp5       string   `xml:",attr"`
	Slot0        Slot0    `xml:""`
	Slot1        Slot1    `xml:""`
	Slot2        Slot2    `xml:""`
	Slot3        Slot3    `xml:""`
	Slot4        Slot4    `xml:""`
	Slot5        Slot5    `xml:""`
}

type GenericAmp struct {
	Bypass       int
	Mute         int
	OutputVolume int
	Model        string
	Amp          Amp
}

func fromAmpA(amp AmpA) GenericAmp {
	return GenericAmp{
		Bypass:       amp.Bypass,
		Mute:         amp.Mute,
		OutputVolume: amp.OutputVolume,
		Model:        amp.Model,
		Amp:          amp.Amp,
	}
}

func fromAmpB(amp AmpB) GenericAmp {
	return GenericAmp{
		Bypass:       amp.Bypass,
		Mute:         amp.Mute,
		OutputVolume: amp.OutputVolume,
		Model:        amp.Model,
		Amp:          amp.Amp,
	}
}

func fromAmpC(amp AmpC) GenericAmp {
	return GenericAmp{
		Bypass:       amp.Bypass,
		Mute:         amp.Mute,
		OutputVolume: amp.OutputVolume,
		Model:        amp.Model,
		Amp:          amp.Amp,
	}
}

type Amp struct {
	XMLName xml.Name   `xml:"Amp"`
	Attrs   []xml.Attr `xml:",any,attr"`
}

type AmpARaw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Amp     Amp        `xml:""`
}

type AmpA struct {
	XMLName      xml.Name `xml:"AmpA"`
	Bypass       int      `xml:",attr"`
	Mute         int      `xml:",attr"`
	OutputVolume int      `xml:",attr"`
	Model        string   `xml:",attr"`
	Amp          Amp      `xml:""`
}

type AmpBRaw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Amp     Amp        `xml:""`
}

type AmpB struct {
	XMLName      xml.Name `xml:"AmpB"`
	Bypass       int      `xml:",attr"`
	Mute         int      `xml:",attr"`
	OutputVolume int      `xml:",attr"`
	Model        string   `xml:",attr"`
	Amp          Amp      `xml:""`
}

type AmpCRaw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Amp     Amp        `xml:""`
}

type AmpC struct {
	XMLName      xml.Name `xml:"AmpC"`
	Bypass       int      `xml:",attr"`
	Mute         int      `xml:",attr"`
	OutputVolume int      `xml:",attr"`
	Model        string   `xml:",attr"`
	Amp          Amp      `xml:""`
}

func fromLoopFxA(stomp LoopFxA) GenericStomp {
	return GenericStomp{
		Bypass:       stomp.Bypass,
		Mute:         stomp.Mute,
		OutputVolume: stomp.OutputVolume,
		Stomp0:       stomp.Stomp0,
		Stomp1:       stomp.Stomp1,
		Stomp2:       stomp.Stomp2,
		Stomp3:       stomp.Stomp3,
		Stomp4:       "",
		Stomp5:       "",
		Slot0:        stomp.Slot0,
		Slot1:        stomp.Slot1,
		Slot2:        stomp.Slot2,
		Slot3:        stomp.Slot3,
		Slot4:        Slot4{},
		Slot5:        Slot5{},
		StompCount:   4,
	}
}

type LoopFxARaw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Slot0   Slot0      `xml:""`
	Slot1   Slot1      `xml:""`
	Slot2   Slot2      `xml:""`
	Slot3   Slot3      `xml:""`
}

type LoopFxA struct {
	XMLName      xml.Name `xml:"LoopFxA"`
	Bypass       int      `xml:",attr"`
	Mute         int      `xml:",attr"`
	OutputVolume int      `xml:",attr"`
	Stomp0       string   `xml:",attr"`
	Stomp1       string   `xml:",attr"`
	Stomp2       string   `xml:",attr"`
	Stomp3       string   `xml:",attr"`
	Slot0        Slot0    `xml:""`
	Slot1        Slot1    `xml:""`
	Slot2        Slot2    `xml:""`
	Slot3        Slot3    `xml:""`
}

func fromLoopFxB(stomp LoopFxB) GenericStomp {
	return GenericStomp{
		Bypass:       stomp.Bypass,
		Mute:         stomp.Mute,
		OutputVolume: stomp.OutputVolume,
		Stomp0:       stomp.Stomp0,
		Stomp1:       stomp.Stomp1,
		Stomp2:       stomp.Stomp2,
		Stomp3:       stomp.Stomp3,
		Stomp4:       "",
		Stomp5:       "",
		Slot0:        stomp.Slot0,
		Slot1:        stomp.Slot1,
		Slot2:        stomp.Slot2,
		Slot3:        stomp.Slot3,
		Slot4:        Slot4{},
		Slot5:        Slot5{},
		StompCount:   4,
	}
}

type LoopFxBRaw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Slot0   Slot0      `xml:""`
	Slot1   Slot1      `xml:""`
	Slot2   Slot2      `xml:""`
	Slot3   Slot3      `xml:""`
}

type LoopFxB struct {
	XMLName      xml.Name `xml:"LoopFxB"`
	Bypass       int      `xml:",attr"`
	Mute         int      `xml:",attr"`
	OutputVolume int      `xml:",attr"`
	Stomp0       string   `xml:",attr"`
	Stomp1       string   `xml:",attr"`
	Stomp2       string   `xml:",attr"`
	Stomp3       string   `xml:",attr"`
	Slot0        Slot0    `xml:""`
	Slot1        Slot1    `xml:""`
	Slot2        Slot2    `xml:""`
	Slot3        Slot3    `xml:""`
}

func fromLoopFxC(stomp LoopFxC) GenericStomp {
	return GenericStomp{
		Bypass:       stomp.Bypass,
		Mute:         stomp.Mute,
		OutputVolume: stomp.OutputVolume,
		Stomp0:       stomp.Stomp0,
		Stomp1:       stomp.Stomp1,
		Stomp2:       stomp.Stomp2,
		Stomp3:       stomp.Stomp3,
		Stomp4:       "",
		Stomp5:       "",
		Slot0:        stomp.Slot0,
		Slot1:        stomp.Slot1,
		Slot2:        stomp.Slot2,
		Slot3:        stomp.Slot3,
		Slot4:        Slot4{},
		Slot5:        Slot5{},
		StompCount:   4,
	}
}

type LoopFxCRaw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Slot0   Slot0      `xml:""`
	Slot1   Slot1      `xml:""`
	Slot2   Slot2      `xml:""`
	Slot3   Slot3      `xml:""`
}

type LoopFxC struct {
	XMLName      xml.Name `xml:"LoopFxC"`
	Bypass       int      `xml:",attr"`
	Mute         int      `xml:",attr"`
	OutputVolume int      `xml:",attr"`
	Stomp0       string   `xml:",attr"`
	Stomp1       string   `xml:",attr"`
	Stomp2       string   `xml:",attr"`
	Stomp3       string   `xml:",attr"`
	Slot0        Slot0    `xml:""`
	Slot1        Slot1    `xml:""`
	Slot2        Slot2    `xml:""`
	Slot3        Slot3    `xml:""`
}

type Cab struct {
	XMLName xml.Name   `xml:"Cab"`
	Attrs   []xml.Attr `xml:",any,attr"`
}

type GenericCab struct {
	Bypass                     int
	Mute                       int
	CabModel                   string
	SpeakerModel0              string
	SpeakerModel1              string
	SpeakerModel2              string
	SpeakerModel3              string
	Studio_Mic1_Level          float32
	Studio_Mic1_Pan            float32
	Studio_Mic1_Mute           int
	Studio_Mic1_Angle          int
	Studio_Mic1_Solo           int
	Studio_Mic1_Phase          int
	Studio_Mic2_Level          float32
	Studio_Mic2_Pan            float32
	Studio_Mic2_Mute           int
	Studio_Mic2_Angle          int
	Studio_Mic2_Solo           int
	Studio_Mic2_Phase          int
	Studio_Room_Level          float32
	Studio_Room_Width          float32
	Studio_Room_Mute           int
	Studio_Room_Solo           int
	Studio_Room_Phase          int
	Studio_Bus_Level           float32
	Studio_Bus_Pan             float32
	Studio_Bus_Mute            int
	Studio_Bus_Solo            int
	Studio_Bus_Phase           int
	Studio_DI_Level            float32
	Studio_DI_Pan              float32
	Studio_DI_Mute             int
	Studio_DI_Solo             int
	Studio_DI_Phase            int
	DI_PhaseDelay              float32
	Studio_LeslieCab_Horn_VolL float32
	Studio_LeslieCab_Horn_VolR float32
	Studio_LeslieCab_Horn_PanL float32
	Studio_LeslieCab_Horn_PanR float32
	Studio_LeslieCab_Horn_Mute int
	Studio_LeslieCab_Horn_Solo int
	Studio_LeslieCab_Drum_VolL float32
	Studio_LeslieCab_Drum_VolR float32
	Studio_LeslieCab_Drum_PanL float32
	Studio_LeslieCab_Drum_PanR float32
	Studio_LeslieCab_Drum_Mute int
	Studio_LeslieCab_Drum_Solo int
	Cab                        Cab
}

func fromCabA(cab CabA) GenericCab {
	return GenericCab{
		Bypass:                     cab.Bypass,
		Mute:                       cab.Mute,
		CabModel:                   cab.CabModel,
		SpeakerModel0:              cab.SpeakerModel0,
		SpeakerModel1:              cab.SpeakerModel1,
		SpeakerModel2:              cab.SpeakerModel2,
		SpeakerModel3:              cab.SpeakerModel3,
		Studio_Mic1_Level:          cab.Studio_Mic1_Level,
		Studio_Mic1_Pan:            cab.Studio_Mic1_Pan,
		Studio_Mic1_Mute:           cab.Studio_Mic1_Mute,
		Studio_Mic1_Angle:          cab.Studio_Mic1_Angle,
		Studio_Mic1_Solo:           cab.Studio_Mic1_Solo,
		Studio_Mic1_Phase:          cab.Studio_Mic1_Phase,
		Studio_Mic2_Level:          cab.Studio_Mic2_Level,
		Studio_Mic2_Pan:            cab.Studio_Mic2_Pan,
		Studio_Mic2_Mute:           cab.Studio_Mic2_Mute,
		Studio_Mic2_Angle:          cab.Studio_Mic2_Angle,
		Studio_Mic2_Solo:           cab.Studio_Mic2_Solo,
		Studio_Mic2_Phase:          cab.Studio_Mic2_Phase,
		Studio_Room_Level:          cab.Studio_Room_Level,
		Studio_Room_Width:          cab.Studio_Room_Width,
		Studio_Room_Mute:           cab.Studio_Room_Mute,
		Studio_Room_Solo:           cab.Studio_Room_Solo,
		Studio_Room_Phase:          cab.Studio_Room_Phase,
		Studio_Bus_Level:           cab.Studio_Bus_Level,
		Studio_Bus_Pan:             cab.Studio_Bus_Pan,
		Studio_Bus_Mute:            cab.Studio_Bus_Mute,
		Studio_Bus_Solo:            cab.Studio_Bus_Solo,
		Studio_Bus_Phase:           cab.Studio_Bus_Phase,
		Studio_DI_Level:            cab.Studio_DI_Level,
		Studio_DI_Pan:              cab.Studio_DI_Pan,
		Studio_DI_Mute:             cab.Studio_DI_Mute,
		Studio_DI_Solo:             cab.Studio_DI_Solo,
		Studio_DI_Phase:            cab.Studio_DI_Phase,
		DI_PhaseDelay:              cab.DI_PhaseDelay,
		Studio_LeslieCab_Horn_VolL: cab.Studio_LeslieCab_Horn_VolL,
		Studio_LeslieCab_Horn_VolR: cab.Studio_LeslieCab_Horn_VolR,
		Studio_LeslieCab_Horn_PanL: cab.Studio_LeslieCab_Horn_PanL,
		Studio_LeslieCab_Horn_PanR: cab.Studio_LeslieCab_Horn_PanR,
		Studio_LeslieCab_Horn_Mute: cab.Studio_LeslieCab_Horn_Mute,
		Studio_LeslieCab_Horn_Solo: cab.Studio_LeslieCab_Horn_Solo,
		Studio_LeslieCab_Drum_VolL: cab.Studio_LeslieCab_Drum_VolL,
		Studio_LeslieCab_Drum_VolR: cab.Studio_LeslieCab_Drum_VolR,
		Studio_LeslieCab_Drum_PanL: cab.Studio_LeslieCab_Drum_PanL,
		Studio_LeslieCab_Drum_PanR: cab.Studio_LeslieCab_Drum_PanR,
		Studio_LeslieCab_Drum_Mute: cab.Studio_LeslieCab_Drum_Mute,
		Studio_LeslieCab_Drum_Solo: cab.Studio_LeslieCab_Drum_Solo,
		Cab:                        cab.Cab,
	}
}

type CabARaw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Cab     Cab        `xml:""`
}

type CabA struct {
	XMLName                    xml.Name `xml:"CabA"`
	Bypass                     int      `xml:",attr"`
	Mute                       int      `xml:",attr"`
	CabModel                   string   `xml:",attr"`
	SpeakerModel0              string   `xml:",attr"`
	SpeakerModel1              string   `xml:",attr"`
	SpeakerModel2              string   `xml:",attr"`
	SpeakerModel3              string   `xml:",attr"`
	Studio_Mic1_Level          float32  `xml:",attr"`
	Studio_Mic1_Pan            float32  `xml:",attr"`
	Studio_Mic1_Mute           int      `xml:",attr"`
	Studio_Mic1_Angle          int      `xml:",attr"`
	Studio_Mic1_Solo           int      `xml:",attr"`
	Studio_Mic1_Phase          int      `xml:",attr"`
	Studio_Mic2_Level          float32  `xml:",attr"`
	Studio_Mic2_Pan            float32  `xml:",attr"`
	Studio_Mic2_Mute           int      `xml:",attr"`
	Studio_Mic2_Angle          int      `xml:",attr"`
	Studio_Mic2_Solo           int      `xml:",attr"`
	Studio_Mic2_Phase          int      `xml:",attr"`
	Studio_Room_Level          float32  `xml:",attr"`
	Studio_Room_Width          float32  `xml:",attr"`
	Studio_Room_Mute           int      `xml:",attr"`
	Studio_Room_Solo           int      `xml:",attr"`
	Studio_Room_Phase          int      `xml:",attr"`
	Studio_Bus_Level           float32  `xml:",attr"`
	Studio_Bus_Pan             float32  `xml:",attr"`
	Studio_Bus_Mute            int      `xml:",attr"`
	Studio_Bus_Solo            int      `xml:",attr"`
	Studio_Bus_Phase           int      `xml:",attr"`
	Studio_DI_Level            float32  `xml:",attr"`
	Studio_DI_Pan              float32  `xml:",attr"`
	Studio_DI_Mute             int      `xml:",attr"`
	Studio_DI_Solo             int      `xml:",attr"`
	Studio_DI_Phase            int      `xml:",attr"`
	DI_PhaseDelay              float32  `xml:",attr"`
	Studio_LeslieCab_Horn_VolL float32  `xml:",attr"`
	Studio_LeslieCab_Horn_VolR float32  `xml:",attr"`
	Studio_LeslieCab_Horn_PanL float32  `xml:",attr"`
	Studio_LeslieCab_Horn_PanR float32  `xml:",attr"`
	Studio_LeslieCab_Horn_Mute int      `xml:",attr"`
	Studio_LeslieCab_Horn_Solo int      `xml:",attr"`
	Studio_LeslieCab_Drum_VolL float32  `xml:",attr"`
	Studio_LeslieCab_Drum_VolR float32  `xml:",attr"`
	Studio_LeslieCab_Drum_PanL float32  `xml:",attr"`
	Studio_LeslieCab_Drum_PanR float32  `xml:",attr"`
	Studio_LeslieCab_Drum_Mute int      `xml:",attr"`
	Studio_LeslieCab_Drum_Solo int      `xml:",attr"`
	Cab                        Cab      `xml:""`
}

func fromCabB(cab CabB) GenericCab {
	return GenericCab{
		Bypass:                     cab.Bypass,
		Mute:                       cab.Mute,
		CabModel:                   cab.CabModel,
		SpeakerModel0:              cab.SpeakerModel0,
		SpeakerModel1:              cab.SpeakerModel1,
		SpeakerModel2:              cab.SpeakerModel2,
		SpeakerModel3:              cab.SpeakerModel3,
		Studio_Mic1_Level:          cab.Studio_Mic1_Level,
		Studio_Mic1_Pan:            cab.Studio_Mic1_Pan,
		Studio_Mic1_Mute:           cab.Studio_Mic1_Mute,
		Studio_Mic1_Angle:          cab.Studio_Mic1_Angle,
		Studio_Mic1_Solo:           cab.Studio_Mic1_Solo,
		Studio_Mic1_Phase:          cab.Studio_Mic1_Phase,
		Studio_Mic2_Level:          cab.Studio_Mic2_Level,
		Studio_Mic2_Pan:            cab.Studio_Mic2_Pan,
		Studio_Mic2_Mute:           cab.Studio_Mic2_Mute,
		Studio_Mic2_Angle:          cab.Studio_Mic2_Angle,
		Studio_Mic2_Solo:           cab.Studio_Mic2_Solo,
		Studio_Mic2_Phase:          cab.Studio_Mic2_Phase,
		Studio_Room_Level:          cab.Studio_Room_Level,
		Studio_Room_Width:          cab.Studio_Room_Width,
		Studio_Room_Mute:           cab.Studio_Room_Mute,
		Studio_Room_Solo:           cab.Studio_Room_Solo,
		Studio_Room_Phase:          cab.Studio_Room_Phase,
		Studio_Bus_Level:           cab.Studio_Bus_Level,
		Studio_Bus_Pan:             cab.Studio_Bus_Pan,
		Studio_Bus_Mute:            cab.Studio_Bus_Mute,
		Studio_Bus_Solo:            cab.Studio_Bus_Solo,
		Studio_Bus_Phase:           cab.Studio_Bus_Phase,
		Studio_DI_Level:            cab.Studio_DI_Level,
		Studio_DI_Pan:              cab.Studio_DI_Pan,
		Studio_DI_Mute:             cab.Studio_DI_Mute,
		Studio_DI_Solo:             cab.Studio_DI_Solo,
		Studio_DI_Phase:            cab.Studio_DI_Phase,
		DI_PhaseDelay:              cab.DI_PhaseDelay,
		Studio_LeslieCab_Horn_VolL: cab.Studio_LeslieCab_Horn_VolL,
		Studio_LeslieCab_Horn_VolR: cab.Studio_LeslieCab_Horn_VolR,
		Studio_LeslieCab_Horn_PanL: cab.Studio_LeslieCab_Horn_PanL,
		Studio_LeslieCab_Horn_PanR: cab.Studio_LeslieCab_Horn_PanR,
		Studio_LeslieCab_Horn_Mute: cab.Studio_LeslieCab_Horn_Mute,
		Studio_LeslieCab_Horn_Solo: cab.Studio_LeslieCab_Horn_Solo,
		Studio_LeslieCab_Drum_VolL: cab.Studio_LeslieCab_Drum_VolL,
		Studio_LeslieCab_Drum_VolR: cab.Studio_LeslieCab_Drum_VolR,
		Studio_LeslieCab_Drum_PanL: cab.Studio_LeslieCab_Drum_PanL,
		Studio_LeslieCab_Drum_PanR: cab.Studio_LeslieCab_Drum_PanR,
		Studio_LeslieCab_Drum_Mute: cab.Studio_LeslieCab_Drum_Mute,
		Studio_LeslieCab_Drum_Solo: cab.Studio_LeslieCab_Drum_Solo,
		Cab:                        cab.Cab,
	}
}

type CabBRaw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Cab     Cab        `xml:""`
}

type CabB struct {
	XMLName                    xml.Name `xml:"CabB"`
	Bypass                     int      `xml:",attr"`
	Mute                       int      `xml:",attr"`
	CabModel                   string   `xml:",attr"`
	SpeakerModel0              string   `xml:",attr"`
	SpeakerModel1              string   `xml:",attr"`
	SpeakerModel2              string   `xml:",attr"`
	SpeakerModel3              string   `xml:",attr"`
	Studio_Mic1_Level          float32  `xml:",attr"`
	Studio_Mic1_Pan            float32  `xml:",attr"`
	Studio_Mic1_Mute           int      `xml:",attr"`
	Studio_Mic1_Angle          int      `xml:",attr"`
	Studio_Mic1_Solo           int      `xml:",attr"`
	Studio_Mic1_Phase          int      `xml:",attr"`
	Studio_Mic2_Level          float32  `xml:",attr"`
	Studio_Mic2_Pan            float32  `xml:",attr"`
	Studio_Mic2_Mute           int      `xml:",attr"`
	Studio_Mic2_Angle          int      `xml:",attr"`
	Studio_Mic2_Solo           int      `xml:",attr"`
	Studio_Mic2_Phase          int      `xml:",attr"`
	Studio_Room_Level          float32  `xml:",attr"`
	Studio_Room_Width          float32  `xml:",attr"`
	Studio_Room_Mute           int      `xml:",attr"`
	Studio_Room_Solo           int      `xml:",attr"`
	Studio_Room_Phase          int      `xml:",attr"`
	Studio_Bus_Level           float32  `xml:",attr"`
	Studio_Bus_Pan             float32  `xml:",attr"`
	Studio_Bus_Mute            int      `xml:",attr"`
	Studio_Bus_Solo            int      `xml:",attr"`
	Studio_Bus_Phase           int      `xml:",attr"`
	Studio_DI_Level            float32  `xml:",attr"`
	Studio_DI_Pan              float32  `xml:",attr"`
	Studio_DI_Mute             int      `xml:",attr"`
	Studio_DI_Solo             int      `xml:",attr"`
	Studio_DI_Phase            int      `xml:",attr"`
	DI_PhaseDelay              float32  `xml:",attr"`
	Studio_LeslieCab_Horn_VolL float32  `xml:",attr"`
	Studio_LeslieCab_Horn_VolR float32  `xml:",attr"`
	Studio_LeslieCab_Horn_PanL float32  `xml:",attr"`
	Studio_LeslieCab_Horn_PanR float32  `xml:",attr"`
	Studio_LeslieCab_Horn_Mute int      `xml:",attr"`
	Studio_LeslieCab_Horn_Solo int      `xml:",attr"`
	Studio_LeslieCab_Drum_VolL float32  `xml:",attr"`
	Studio_LeslieCab_Drum_VolR float32  `xml:",attr"`
	Studio_LeslieCab_Drum_PanL float32  `xml:",attr"`
	Studio_LeslieCab_Drum_PanR float32  `xml:",attr"`
	Studio_LeslieCab_Drum_Mute int      `xml:",attr"`
	Studio_LeslieCab_Drum_Solo int      `xml:",attr"`
	Cab                        Cab      `xml:""`
}

func fromCabC(cab CabC) GenericCab {
	return GenericCab{
		Bypass:                     cab.Bypass,
		Mute:                       cab.Mute,
		CabModel:                   cab.CabModel,
		SpeakerModel0:              cab.SpeakerModel0,
		SpeakerModel1:              cab.SpeakerModel1,
		SpeakerModel2:              cab.SpeakerModel2,
		SpeakerModel3:              cab.SpeakerModel3,
		Studio_Mic1_Level:          cab.Studio_Mic1_Level,
		Studio_Mic1_Pan:            cab.Studio_Mic1_Pan,
		Studio_Mic1_Mute:           cab.Studio_Mic1_Mute,
		Studio_Mic1_Angle:          cab.Studio_Mic1_Angle,
		Studio_Mic1_Solo:           cab.Studio_Mic1_Solo,
		Studio_Mic1_Phase:          cab.Studio_Mic1_Phase,
		Studio_Mic2_Level:          cab.Studio_Mic2_Level,
		Studio_Mic2_Pan:            cab.Studio_Mic2_Pan,
		Studio_Mic2_Mute:           cab.Studio_Mic2_Mute,
		Studio_Mic2_Angle:          cab.Studio_Mic2_Angle,
		Studio_Mic2_Solo:           cab.Studio_Mic2_Solo,
		Studio_Mic2_Phase:          cab.Studio_Mic2_Phase,
		Studio_Room_Level:          cab.Studio_Room_Level,
		Studio_Room_Width:          cab.Studio_Room_Width,
		Studio_Room_Mute:           cab.Studio_Room_Mute,
		Studio_Room_Solo:           cab.Studio_Room_Solo,
		Studio_Room_Phase:          cab.Studio_Room_Phase,
		Studio_Bus_Level:           cab.Studio_Bus_Level,
		Studio_Bus_Pan:             cab.Studio_Bus_Pan,
		Studio_Bus_Mute:            cab.Studio_Bus_Mute,
		Studio_Bus_Solo:            cab.Studio_Bus_Solo,
		Studio_Bus_Phase:           cab.Studio_Bus_Phase,
		Studio_DI_Level:            cab.Studio_DI_Level,
		Studio_DI_Pan:              cab.Studio_DI_Pan,
		Studio_DI_Mute:             cab.Studio_DI_Mute,
		Studio_DI_Solo:             cab.Studio_DI_Solo,
		Studio_DI_Phase:            cab.Studio_DI_Phase,
		DI_PhaseDelay:              cab.DI_PhaseDelay,
		Studio_LeslieCab_Horn_VolL: cab.Studio_LeslieCab_Horn_VolL,
		Studio_LeslieCab_Horn_VolR: cab.Studio_LeslieCab_Horn_VolR,
		Studio_LeslieCab_Horn_PanL: cab.Studio_LeslieCab_Horn_PanL,
		Studio_LeslieCab_Horn_PanR: cab.Studio_LeslieCab_Horn_PanR,
		Studio_LeslieCab_Horn_Mute: cab.Studio_LeslieCab_Horn_Mute,
		Studio_LeslieCab_Horn_Solo: cab.Studio_LeslieCab_Horn_Solo,
		Studio_LeslieCab_Drum_VolL: cab.Studio_LeslieCab_Drum_VolL,
		Studio_LeslieCab_Drum_VolR: cab.Studio_LeslieCab_Drum_VolR,
		Studio_LeslieCab_Drum_PanL: cab.Studio_LeslieCab_Drum_PanL,
		Studio_LeslieCab_Drum_PanR: cab.Studio_LeslieCab_Drum_PanR,
		Studio_LeslieCab_Drum_Mute: cab.Studio_LeslieCab_Drum_Mute,
		Studio_LeslieCab_Drum_Solo: cab.Studio_LeslieCab_Drum_Solo,
		Cab:                        cab.Cab,
	}
}

type CabCRaw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Cab     Cab        `xml:""`
}

type CabC struct {
	XMLName                    xml.Name `xml:"CabC"`
	Bypass                     int      `xml:",attr"`
	Mute                       int      `xml:",attr"`
	CabModel                   string   `xml:",attr"`
	SpeakerModel0              string   `xml:",attr"`
	SpeakerModel1              string   `xml:",attr"`
	SpeakerModel2              string   `xml:",attr"`
	SpeakerModel3              string   `xml:",attr"`
	Studio_Mic1_Level          float32  `xml:",attr"`
	Studio_Mic1_Pan            float32  `xml:",attr"`
	Studio_Mic1_Mute           int      `xml:",attr"`
	Studio_Mic1_Angle          int      `xml:",attr"`
	Studio_Mic1_Solo           int      `xml:",attr"`
	Studio_Mic1_Phase          int      `xml:",attr"`
	Studio_Mic2_Level          float32  `xml:",attr"`
	Studio_Mic2_Pan            float32  `xml:",attr"`
	Studio_Mic2_Mute           int      `xml:",attr"`
	Studio_Mic2_Angle          int      `xml:",attr"`
	Studio_Mic2_Solo           int      `xml:",attr"`
	Studio_Mic2_Phase          int      `xml:",attr"`
	Studio_Room_Level          float32  `xml:",attr"`
	Studio_Room_Width          float32  `xml:",attr"`
	Studio_Room_Mute           int      `xml:",attr"`
	Studio_Room_Solo           int      `xml:",attr"`
	Studio_Room_Phase          int      `xml:",attr"`
	Studio_Bus_Level           float32  `xml:",attr"`
	Studio_Bus_Pan             float32  `xml:",attr"`
	Studio_Bus_Mute            int      `xml:",attr"`
	Studio_Bus_Solo            int      `xml:",attr"`
	Studio_Bus_Phase           int      `xml:",attr"`
	Studio_DI_Level            float32  `xml:",attr"`
	Studio_DI_Pan              float32  `xml:",attr"`
	Studio_DI_Mute             int      `xml:",attr"`
	Studio_DI_Solo             int      `xml:",attr"`
	Studio_DI_Phase            int      `xml:",attr"`
	DI_PhaseDelay              float32  `xml:",attr"`
	Studio_LeslieCab_Horn_VolL float32  `xml:",attr"`
	Studio_LeslieCab_Horn_VolR float32  `xml:",attr"`
	Studio_LeslieCab_Horn_PanL float32  `xml:",attr"`
	Studio_LeslieCab_Horn_PanR float32  `xml:",attr"`
	Studio_LeslieCab_Horn_Mute int      `xml:",attr"`
	Studio_LeslieCab_Horn_Solo int      `xml:",attr"`
	Studio_LeslieCab_Drum_VolL float32  `xml:",attr"`
	Studio_LeslieCab_Drum_VolR float32  `xml:",attr"`
	Studio_LeslieCab_Drum_PanL float32  `xml:",attr"`
	Studio_LeslieCab_Drum_PanR float32  `xml:",attr"`
	Studio_LeslieCab_Drum_Mute int      `xml:",attr"`
	Studio_LeslieCab_Drum_Solo int      `xml:",attr"`
	Cab                        Cab      `xml:""`
}

type CabAAttrs struct {
	XMLName xml.Name   `xml:"CabA"`
	Attrs   []xml.Attr `xml:",any,attr"`
}

type CabBAttrs struct {
	XMLName xml.Name   `xml:"CabB"`
	Attrs   []xml.Attr `xml:",any,attr"`
}

type CabCAttrs struct {
	XMLName xml.Name   `xml:"CabC"`
	Attrs   []xml.Attr `xml:",any,attr"`
}

type StudioRaw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
}

type Studio struct {
	XMLName                xml.Name `xml:"Studio"`
	Bypass                 int      `xml:",attr"`
	Mute                   int      `xml:",attr"`
	OutputVolume           int      `xml:",attr"`
	OutputPan              float32  `xml:",attr"`
	DI_Level               float32  `xml:",attr"`
	DI_Pan                 float32  `xml:",attr"`
	DI_Mute                int      `xml:",attr"`
	DI_Solo                int      `xml:",attr"`
	DI_Phase               int      `xml:",attr"`
	DI_PhaseDelay          float32  `xml:",attr"`
	Cab1_Mic1_Level        float32  `xml:",attr"`
	Cab1_Mic1_Pan          float32  `xml:",attr"`
	Cab1_Mic1_Mute         int      `xml:",attr"`
	Cab1_Mic1_Solo         int      `xml:",attr"`
	Cab1_Mic1_Phase        int      `xml:",attr"`
	Cab1_Mic2_Level        float32  `xml:",attr"`
	Cab1_Mic2_Pan          float32  `xml:",attr"`
	Cab1_Mic2_Mute         int      `xml:",attr"`
	Cab1_Mic2_Solo         int      `xml:",attr"`
	Cab1_Mic2_Phase        int      `xml:",attr"`
	Cab1_Room_Level        float32  `xml:",attr"`
	Cab1_Room_Width        float32  `xml:",attr"`
	Cab1_Room_Mute         int      `xml:",attr"`
	Cab1_Room_Solo         int      `xml:",attr"`
	Cab1_Room_Phase        int      `xml:",attr"`
	Cab1_Bus_Level         float32  `xml:",attr"`
	Cab1_Bus_Pan           float32  `xml:",attr"`
	Cab1_Bus_Mute          int      `xml:",attr"`
	Cab1_Bus_Solo          int      `xml:",attr"`
	Cab2_Mic1_Level        float32  `xml:",attr"`
	Cab2_Mic1_Pan          float32  `xml:",attr"`
	Cab2_Mic1_Mute         int      `xml:",attr"`
	Cab2_Mic1_Solo         int      `xml:",attr"`
	Cab2_Mic1_Phase        int      `xml:",attr"`
	Cab2_Mic2_Level        float32  `xml:",attr"`
	Cab2_Mic2_Pan          float32  `xml:",attr"`
	Cab2_Mic2_Mute         int      `xml:",attr"`
	Cab2_Mic2_Solo         int      `xml:",attr"`
	Cab2_Mic2_Phase        int      `xml:",attr"`
	Cab2_Room_Level        float32  `xml:",attr"`
	Cab2_Room_Width        float32  `xml:",attr"`
	Cab2_Room_Mute         int      `xml:",attr"`
	Cab2_Room_Solo         int      `xml:",attr"`
	Cab2_Room_Phase        int      `xml:",attr"`
	Cab2_Bus_Level         float32  `xml:",attr"`
	Cab2_Bus_Pan           float32  `xml:",attr"`
	Cab2_Bus_Mute          int      `xml:",attr"`
	Cab2_Bus_Solo          int      `xml:",attr"`
	Cab3_Mic1_Level        float32  `xml:",attr"`
	Cab3_Mic1_Pan          float32  `xml:",attr"`
	Cab3_Mic1_Mute         int      `xml:",attr"`
	Cab3_Mic1_Solo         int      `xml:",attr"`
	Cab3_Mic1_Phase        int      `xml:",attr"`
	Cab3_Mic2_Level        float32  `xml:",attr"`
	Cab3_Mic2_Pan          float32  `xml:",attr"`
	Cab3_Mic2_Mute         int      `xml:",attr"`
	Cab3_Mic2_Solo         int      `xml:",attr"`
	Cab3_Mic2_Phase        int      `xml:",attr"`
	Cab3_Room_Level        float32  `xml:",attr"`
	Cab3_Room_Width        float32  `xml:",attr"`
	Cab3_Room_Mute         int      `xml:",attr"`
	Cab3_Room_Solo         int      `xml:",attr"`
	Cab3_Room_Phase        int      `xml:",attr"`
	Cab3_Bus_Level         float32  `xml:",attr"`
	Cab3_Bus_Pan           float32  `xml:",attr"`
	Cab3_Bus_Mute          int      `xml:",attr"`
	Cab3_Bus_Solo          int      `xml:",attr"`
	MasterLevel            float32  `xml:",attr"`
	Cab1_Leslie_Horn_Width float32  `xml:",attr"`
	Cab1_Leslie_Drum_Width float32  `xml:",attr"`
	Cab2_Leslie_Horn_Width float32  `xml:",attr"`
	Cab2_Leslie_Drum_Width float32  `xml:",attr"`
	Cab3_Leslie_Horn_Width float32  `xml:",attr"`
	Cab3_Leslie_Drum_Width float32  `xml:",attr"`
}

func fromRackA(stomp RackA) GenericStomp {
	return GenericStomp{
		Bypass:       stomp.Bypass,
		Mute:         stomp.Mute,
		OutputVolume: stomp.OutputVolume,
		Stomp0:       stomp.Stomp0,
		Stomp1:       stomp.Stomp1,
		Stomp2:       "",
		Stomp3:       "",
		Stomp4:       "",
		Stomp5:       "",
		Slot0:        stomp.Slot0,
		Slot1:        stomp.Slot1,
		Slot2:        Slot2{},
		Slot3:        Slot3{},
		Slot4:        Slot4{},
		Slot5:        Slot5{},
		StompCount:   2,
	}
}

type RackARaw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Slot0   Slot0      `xml:""`
	Slot1   Slot1      `xml:""`
}

type RackA struct {
	XMLName      xml.Name `xml:"RackA"`
	Bypass       int      `xml:",attr"`
	Mute         int      `xml:",attr"`
	OutputVolume int      `xml:",attr"`
	Stomp0       string   `xml:",attr"`
	Stomp1       string   `xml:",attr"`
	Slot0        Slot0    `xml:""`
	Slot1        Slot1    `xml:""`
}

func fromRackB(stomp RackB) GenericStomp {
	return GenericStomp{
		Bypass:       stomp.Bypass,
		Mute:         stomp.Mute,
		OutputVolume: stomp.OutputVolume,
		Stomp0:       stomp.Stomp0,
		Stomp1:       stomp.Stomp1,
		Stomp2:       "",
		Stomp3:       "",
		Stomp4:       "",
		Stomp5:       "",
		Slot0:        stomp.Slot0,
		Slot1:        stomp.Slot1,
		Slot2:        Slot2{},
		Slot3:        Slot3{},
		Slot4:        Slot4{},
		Slot5:        Slot5{},
		StompCount:   2,
	}
}

type RackBRaw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Slot0   Slot0      `xml:""`
	Slot1   Slot1      `xml:""`
}

type RackB struct {
	XMLName      xml.Name `xml:"RackB"`
	Bypass       int      `xml:",attr"`
	Mute         int      `xml:",attr"`
	OutputVolume int      `xml:",attr"`
	Stomp0       string   `xml:",attr"`
	Stomp1       string   `xml:",attr"`
	Slot0        Slot0    `xml:""`
	Slot1        Slot1    `xml:""`
}

func fromRackC(stomp RackC) GenericStomp {
	return GenericStomp{
		Bypass:       stomp.Bypass,
		Mute:         stomp.Mute,
		OutputVolume: stomp.OutputVolume,
		Stomp0:       stomp.Stomp0,
		Stomp1:       stomp.Stomp1,
		Stomp2:       "",
		Stomp3:       "",
		Stomp4:       "",
		Stomp5:       "",
		Slot0:        stomp.Slot0,
		Slot1:        stomp.Slot1,
		Slot2:        Slot2{},
		Slot3:        Slot3{},
		Slot4:        Slot4{},
		Slot5:        Slot5{},
		StompCount:   2,
	}
}

type RackCRaw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Slot0   Slot0      `xml:""`
	Slot1   Slot1      `xml:""`
}

type RackC struct {
	XMLName      xml.Name `xml:"RackC"`
	Bypass       int      `xml:",attr"`
	Mute         int      `xml:",attr"`
	OutputVolume int      `xml:",attr"`
	Stomp0       string   `xml:",attr"`
	Stomp1       string   `xml:",attr"`
	Slot0        Slot0    `xml:""`
	Slot1        Slot1    `xml:""`
}

func fromRackDI(stomp RackDI) GenericStomp {
	return GenericStomp{
		Bypass:       stomp.Bypass,
		Mute:         stomp.Mute,
		OutputVolume: stomp.OutputVolume,
		Stomp0:       stomp.Stomp0,
		Stomp1:       stomp.Stomp1,
		Stomp2:       "",
		Stomp3:       "",
		Stomp4:       "",
		Stomp5:       "",
		Slot0:        stomp.Slot0,
		Slot1:        stomp.Slot1,
		Slot2:        Slot2{},
		Slot3:        Slot3{},
		Slot4:        Slot4{},
		Slot5:        Slot5{},
		StompCount:   2,
	}
}

type RackDIRaw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Slot0   Slot0      `xml:""`
	Slot1   Slot1      `xml:""`
}

type RackDI struct {
	XMLName      xml.Name `xml:"RackDI"`
	Bypass       int      `xml:",attr"`
	Mute         int      `xml:",attr"`
	OutputVolume int      `xml:",attr"`
	Stomp0       string   `xml:",attr"`
	Stomp1       string   `xml:",attr"`
	Slot0        Slot0    `xml:""`
	Slot1        Slot1    `xml:""`
}

func fromRackMaster(stomp RackMaster) GenericStomp {
	return GenericStomp{
		Bypass:       stomp.Bypass,
		Mute:         stomp.Mute,
		OutputVolume: stomp.OutputVolume,
		Stomp0:       stomp.Stomp0,
		Stomp1:       stomp.Stomp1,
		Stomp2:       stomp.Stomp2,
		Stomp3:       stomp.Stomp3,
		Stomp4:       stomp.Stomp4,
		Stomp5:       stomp.Stomp5,
		Slot0:        stomp.Slot0,
		Slot1:        stomp.Slot1,
		Slot2:        stomp.Slot2,
		Slot3:        stomp.Slot3,
		Slot4:        stomp.Slot4,
		Slot5:        stomp.Slot5,
		StompCount:   6,
	}
}

type RackMasterRaw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Slot0   Slot0      `xml:""`
	Slot1   Slot1      `xml:""`
	Slot2   Slot2      `xml:""`
	Slot3   Slot3      `xml:""`
	Slot4   Slot4      `xml:""`
	Slot5   Slot5      `xml:""`
}

type RackMaster struct {
	XMLName      xml.Name `xml:"RackMaster"`
	Bypass       int      `xml:",attr"`
	Mute         int      `xml:",attr"`
	OutputVolume int      `xml:",attr"`
	Stomp0       string   `xml:",attr"`
	Stomp1       string   `xml:",attr"`
	Stomp2       string   `xml:",attr"`
	Stomp3       string   `xml:",attr"`
	Stomp4       string   `xml:",attr"`
	Stomp5       string   `xml:",attr"`
	Slot0        Slot0    `xml:""`
	Slot1        Slot1    `xml:""`
	Slot2        Slot2    `xml:""`
	Slot3        Slot3    `xml:""`
	Slot4        Slot4    `xml:""`
	Slot5        Slot5    `xml:""`
}

type OutputRaw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
}

type Output struct {
	XMLName xml.Name `xml:"Output"`
	Output  int      `xml:",attr"`
}

type MidiAssignmentsRaw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
}

type MidiAssignments struct {
	XMLName xml.Name `xml:"MidiAssignments"`
}

type MetaInfoRaw struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
}

type MetaInfo struct {
	XMLName              xml.Name `xml:"MetaInfo"`
	Description          string   `xml:",attr"`
	Style                string   `xml:",attr"`
	SoundCharacter       string   `xml:",attr"`
	Instrument           string   `xml:",attr"`
	Body                 string   `xml:",attr"`
	PickUpPosition       string   `xml:",attr"`
	Artist               string   `xml:",attr"`
	Band                 string   `xml:",attr"`
	Song                 string   `xml:",attr"`
	SongStructureElement string   `xml:",attr"`
	KeyWords             string   `xml:",attr"`
	Type                 string   `xml:",attr"`
}

type SlotList struct {
	slots []GenericSlot
}

func newSlotList() SlotList {
	return SlotList{
		slots: []GenericSlot{},
	}
}

func (s *SlotList) insert(slot GenericSlot) {
	if slot.Attrs != nil {
		s.slots = append([]GenericSlot{slot}, s.slots...)
	}
}

func (s *SlotList) append(slot GenericSlot) {
	if slot.Attrs != nil {
		s.slots = append(s.slots, slot)
	}
}

type GUIDList struct {
	guids []string
}

func newGUIDList() GUIDList {
	return GUIDList{guids: []string{}}
}

func (g *GUIDList) insert(guid string) {
	if guid != EmptySlotGUID {
		g.guids = append([]string{guid}, g.guids...)
	}
}

func (g *GUIDList) append(guid string) {
	if guid != EmptySlotGUID {
		g.guids = append(g.guids, guid)
	}
}

func presetFormatVersion(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	var format PresetXMLFormatOnly
	err = xml.Unmarshal(data, &format)
	return format.Format, err
}

func writeNewGuidToFile(file string) error {
	data, err := ioutil.ReadFile(file)
	format, err := presetFormatVersion(file)
	var out []byte
	switch format {
	case "at4p":
		var p PresetXMLRootOnlyV4
		err = xml.Unmarshal(data, &p)
		newId, _ := uuid.NewRandom()
		p.GUID = newId.String()
		out, err = xml.MarshalIndent(p, "<?xml version=\"1.0\" ?>\n", "    ")
		break
	case "at5p":
		var p PresetXMLRootOnlyV5
		err = xml.Unmarshal(data, &p)
		newId, _ := uuid.NewRandom()
		p.GUID = newId.String()
		out, err = xml.MarshalIndent(p, "<?xml version=\"1.0\" ?>\n", "    ")
		break
	default:
		err = errors.New("unknown preset file format")
	}
	err = ioutil.WriteFile(file, append(out, []byte("\n")...), 0664)
	return err
}

func isInPresetsFolder(startPath string) bool {

	isPresetPath := false
	exhausted := false
	currPath := startPath

	for !isPresetPath && !exhausted {
		isPresetPath = isProfileFolder(filepath.Dir(currPath)) && filepath.Base(currPath) == PresetsFolder
		if !isPresetPath {
			lastPath := currPath
			currPath = filepath.Dir(currPath)
			exhausted = lastPath == currPath
		}
	}

	return isPresetPath
}

func isValidPresetName(path string) bool {
	return !containsWildcards(path) && (strings.LastIndex(path, PresetExtension) == len(path)-len(PresetExtension) || strings.LastIndex(path, PresetExtension4) == len(path)-len(PresetExtension4))
}

func isValidPresetFolderName(path string) bool {
	return strings.Index(path, PresetExtension) == -1 && strings.Index(path, PresetExtension4) == -1
}

func makePresetPath(path string) string {
	if strings.Index(path, PresetsFolder) == 0 {
		path = path[strings.Index(path, PresetsFolder)+len(PresetsFolder):]
	}
	if strings.Index(path, string(filepath.Separator)) == 0 {
		path = path[1:]
	}
	if strings.Index(path, PresetExtension) > -1 && (strings.Index(path, PresetExtension) == len(path)-5 || strings.Index(path, PresetExtension4) == len(path)-5) {
		path = path[:len(path)-5]
	}
	return path
}

func selfClose(xmldata []byte) []byte {
	var closed []byte
	var skip bool
	for i, b := range xmldata {
		// 60 < 62 > 47 /
		if skip && b != 60 && xmldata[i-1] == 62 {
			skip = false
		}
		if b == 62 && len(xmldata) > i+2 && xmldata[i+1] == 60 && xmldata[i+2] == 47 {
			if !skip {
				closed = append(closed, []byte{' ', '/', '>'}...)
			}
			skip = true
		}
		if !skip {
			closed = append(closed, b)
		}
	}
	return closed
}

func truncateSlots(xmldata []byte, total int) []byte {
	var slots []byte
	count := 0
	for _, b := range xmldata {
		if b == 60 {
			count++
		}
		if count <= total {
			slots = append(slots, b)
		}
	}
	i := len(slots) - 1
	maxRemove := len(slots) - 5
	for i > 0 {
		if slots[i] == 32 && i > maxRemove {
			slots = slots[:i]
			i--
		} else {
			i = -1
		}
	}
	return slots
}

func padSlots(xmldata []byte, total int) []byte {
	var slots []byte
	count := 0
	for _, b := range xmldata {
		if b == 60 {
			count++
		}
		slots = append(slots, b)
	}
	first := true
	for count < total {
		if !first {
			slots = append(slots, []byte("    ")...)
		}
		slots = append(slots, []byte("    <Slot"+strconv.Itoa(count)+" />\n")...)
		count++
		first = false
	}
	return slots
}

func truncateOrPadSlots(xmldata []byte, limit int) []byte {
	count := 0
	for _, b := range xmldata {
		if b == 60 {
			count++
		}
	}
	if count > limit {
		return truncateSlots(xmldata, limit)
	}
	if count < limit {
		return padSlots(xmldata, limit)
	}
	return xmldata
}

func splitSlots(xmldata []byte) []string {
	var result []string
	var buf []byte
	inSlot := false
	for _, b := range xmldata {
		if b == 60 {
			inSlot = true
		}
		if inSlot {
			buf = append(buf, b)
		}
		if b == 62 {
			inSlot = false
			result = append(result, string(buf))
			buf = []byte{}
		}
	}
	return result
}

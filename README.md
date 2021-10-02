# AmpliTool

Command line tool for manipulating IK Multimedia Amplitube 5 Presets.

**Make sure to backup your Amplitube Profile directory before using Amplitool.
Or better still, use git to version control your presets.**

## Possible Uses for Amplitool

- Copying, moving, or deleting a single preset is simple in the Amplitube interface,
  but, if you have a large number of presets to copy, move, or delete, this becomes time-consuming.
  Use Amplitool to perform bulk copy, move, or delete operations on multiple presets
  at the same time.
- When you add a piece of gear to a preset in Amplitube, it always has the default settings.  If you
  have another preset with the piece of gear already configured the way you like it, you can use the Amplitool
  copy gear (cpg) command to copy the already configured piece of gear into the 
  new preset.

## Build

Using [Golang](https://golang.org/) 1.15 or higher.

```
go build -o ampt.exe .
```

## Install

Place ampt.exe on the PATH of your system.

## Usage

### List Presets

In Amplitube 5 Profile directory

```
ampt ls Presets
```

List all subfolders as well

```
ampt ls -r Presets
```

### Copy Presets

Copy presets or preset folders

#### Examples

Copy one preset

```
ampt cp Presets/Default.at5p Presets/DefaultCopy.at5p
```

Copy folder and all subfolders

```
ampt cp -r Presets/Defaults Presets/DefaultsCopy
```

Note this will copy the root folder as well, resulting in a new folder named
Presets/DefaultsCopy/Defaults.  To copy just the contents use wildcard matching

```
ampt cp Presets/Defaults/* Presets/DefaultsCopy
```

### Move Presets

Move presets or preset folders

#### Examples

Move one preset

```
ampt mv Presets/Default.at5p Presets/Moved.at5p
```

Move a folder

```
ampt mv Presets/Defaults Presets/Other
```

This results in the entire folder being moved to the target folder, ie. the 
result is a folder named Presets/Other/Defaults.  To move the contents 
of Presets/Defaults use wildcard matching.

```
ampt mv Presets/Defaults/* Presets/Other
```

### Remove Presets

Remove presets and/or preset folders

#### Examples

Remove one preset

```
ampt rm Presets/Default.at5p
```

Remove an entire folder presets

```
ampt rm -r Presets/Defaults
```

### Make Folder

Make a new preset folder.

#### Example

```
ampt mkdir Presets/NewFolder
```

### List Gear

List gear in a preset.

```
ampt lsg Presets/Default.at5p
```

List with all details

```
ampt lsg -d Presets/Default.at5p
```

List raw preset file contents

```
ampt lsg -r Presets/Default.at5p
```

### Copy Gear

Copy a block of gear from one preset to one or more other presets

#### Examples

Append StompB1 from one preset to StompB1 in another

```
ampt cpg Presets/Default.at5p Presets/Other.at5p StompB1
```

Append StompB1 effects to StompA1 in another

```
ampt cpg Presets/Default.at5p Presets/Other.at5p StompB1 StompA1
```

Overwrite StompB1 effects

```
ampt cpg -o Presets/Default.at5p Presets/Other.at5p StompB1
```

Insert StompB1 effects

```
ampt cpg -i Presets/Default.at5p Presets/Other.at5p StompB1
```

Append StompB1 from one preset to StompB1 in all presets in folder

```
ampt cpg Presets/Default.at5p Presets/Amps StompB1
```

Append StompB1 from one preset to StompB1 in all presets in folder and subfolders

```
ampt cpg -r Presets/Default.at5p Presets/Amps StompB1
```

Overwrite Amp and Cab

```
ampt cpg Presets/Default.at5p Presets/Amps AmpA
```

Overwrite only Amp

```
ampt cpg -c Presets/Default.at5p Presets/Amps AmpA
```

Overwrite all amps

```
ampt cpg -aa Presets/Default.at5p Presets/Other.at5p
```

Overwrite all cabs

```
ampt cpg -ac Presets/Default.at5p Presets/Other.at5p
```

Overwrite all effects

```
ampt cpg -ae Presets/Default.at5p Presets/Other.at5p
```

### Remove Gear

Remove effects from StompB1

```
ampt rmg Presets/Default.at5p StompB1
```

Remove effect in Slot0 from StompB1

```
ampt rmg Presets/Default.at5p StompB1 Slot0
```

### Set Gear Attribute

Set an attribute value on an element in a preset.

Set the bypass on Slot0 of StompA1 to 1

```
ampt sg Presets/Default.at5p Preset.StompA1.Slot0.Bypass=1
```

### Import Preset

Import presets from another Amplitube profile directory.  Imported presets
are placed in an "Import" subfolder.

```
ampt import Profile1 Profile2
```

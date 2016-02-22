// Copyright 2016, Timothy Bogdala <tdb@animal-machine.com>
// See the LICENSE file for more details.

package voxfile

import (
	"testing"
)

const (
	voxfileCharacter = "testdata/chr_sword.vox"
)

// TestFileLoad tests the loading and decoding of a vox file
func TestFileLoad(t *testing.T) {
	voxFile, err := DecodeFile(voxfileCharacter)
	if err != nil || voxFile == nil {
		t.Fatalf("Failed to load the VOX file %s.\n%v", voxfileCharacter, err)
	}

	if voxFile.Version != 150 {
		t.Errorf("File version doesn't match (%d).", voxFile.Version)
	}

	if voxFile.SizeX != 20 {
		t.Errorf("File sizeX doesn't match (%d).", voxFile.SizeX)
	}

	if voxFile.SizeY != 21 {
		t.Errorf("File sizeY doesn't match (%d).", voxFile.SizeY)
	}

	if voxFile.SizeZ != 20 {
		t.Errorf("File sizeZ doesn't match (%d).", voxFile.SizeZ)
	}

	// 334 is the magic number for voxels in the test file
	if len(voxFile.Voxels) != 334 {
		t.Errorf("The number of voxels doesn't match (%d).", len(voxFile.Voxels))
	}

	// test a couple colors from the default palette
	if voxFile.Palette == nil {
		t.Errorf("File doesn't have a palette attached to it.")
	}

	if voxFile.Palette[0].R != 0 || voxFile.Palette[0].G != 0 || voxFile.Palette[0].B != 0 || voxFile.Palette[0].A != 0 {
		t.Errorf("Color #0 of the palette doesn't match what it should.")
	}

	if voxFile.Palette[1].R != 0xff || voxFile.Palette[1].G != 0xff || voxFile.Palette[1].B != 0xff || voxFile.Palette[1].A != 0xff {
		t.Errorf("Color #1 of the palette doesn't match what it should.")
	}

	if voxFile.Palette[2].R != 0xff || voxFile.Palette[2].G != 0xcc || voxFile.Palette[2].B != 0xff || voxFile.Palette[2].A != 0xff {
		t.Errorf("Color #2 of the palette doesn't match what it should.")
	}
}

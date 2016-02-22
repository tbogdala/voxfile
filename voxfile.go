// Copyright 2016, Timothy Bogdala <tdb@animal-machine.com>
// See the LICENSE file for more details.

package voxfile

/*

The VOX file format can be found at the following web page:

https://ephtracy.github.io/index.html?page=mv_vox_format

*/

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
)

const (
	// CurrentVersion specifies the supported version for the file loader.
	CurrentVersion = 150
)

// Voxel is the data type that represents one voxel in the file. It has Location
// data as well as an index value pointing to the palette color.
type Voxel struct {
	X     uint8
	Y     uint8
	Z     uint8
	Index uint8
}

// Color represents a color on a palette.
type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

// VoxFile is the struct containing the data from/for a VOX file, which
// represents colored voxels in 3D space.
type VoxFile struct {
	// Version is the version number of the loaded file
	Version uint32

	// SizeX is the size of the X-axis of the voxel data
	SizeX uint32

	// SizeY is the size of the X-axis of the voxel data
	SizeY uint32

	// SizeZ is the size of the X-axis of the voxel data
	SizeZ uint32

	// Voxels is a slice of all loaded voxels from the file
	Voxels []*Voxel

	// Palette is a 256-size palette of colors that is referenced
	// in the Voxel structs
	Palette []*Color
}

// DecodeFile opens the file specified and reads it in as a VOX file.
func DecodeFile(fn string) (*VoxFile, error) {
	// open the file
	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// make a new reader for the file
	reader := bufio.NewReader(file)

	return Decode(reader)
}

// Decode takes the bytes from the stream passed in and attempts
// to create the vox structures.
func Decode(r *bufio.Reader) (*VoxFile, error) {
	voxelFile := new(VoxFile)

	// Read in and test the 'magic' string
	var magic [4]byte
	c, err := r.Read(magic[:])
	if err != nil {
		return nil, fmt.Errorf("File doesn't appear to be a VOX file. %v", err)
	}
	if c != 4 || magic[0] != 'V' || magic[1] != 'O' || magic[2] != 'X' || magic[3] != ' ' {
		return nil, fmt.Errorf("File doesn't appear to be a VOX file. (Magic: %v)", magic)
	}

	// Read in the version number of the file
	var version uint32
	err = binary.Read(r, binary.LittleEndian, &version)
	if err != nil {
		return nil, fmt.Errorf("Couldn't read the version number from the file. %v", err)
	}
	if version != CurrentVersion {
		return nil, fmt.Errorf("Version number from the file (%d) doesn't match the current version (%d).", version, CurrentVersion)
	}
	voxelFile.Version = version

	_, err = readChunk(r, voxelFile)

	// if we didn't have a custom palette, make an instance of the default one
	voxelFile.Palette = instancePalette(defaultPalette)

	return voxelFile, err
}

// readChunk reads a chunk from the VOX file.
func readChunk(r *bufio.Reader, voxFile *VoxFile) (bytesRead uint32, err error) {
	// get the ID
	var chunkID [4]byte
	_, err = r.Read(chunkID[:])
	if err != nil {
		return 0, fmt.Errorf("Failed to read the chunk ID. %v", err)
	}
	chunkIDStr := string(chunkID[:4])

	// get the chunk Size
	var chunkSize uint32
	err = binary.Read(r, binary.LittleEndian, &chunkSize)
	if err != nil {
		return 0, fmt.Errorf("Failed to read the %s chunk size. %v", chunkIDStr, err)
	}

	// get the children size
	var chunkChildrenSize uint32
	err = binary.Read(r, binary.LittleEndian, &chunkChildrenSize)
	if err != nil {
		return 0, fmt.Errorf("Failed to read the %s chunk's children size. %v", chunkIDStr, err)
	}

	// read in the chunk if necessary
	if chunkSize > 0 {
		// Some chunks are handled separately
		if chunkIDStr == "SIZE" {
			// read the dimensions of the chunk
			if chunkSize != 12 {
				return 0, fmt.Errorf("Failed to read the %s chunk. Size should have been 12 but is %d.", chunkIDStr, chunkSize)
			}
			var sizeX uint32
			err = binary.Read(r, binary.LittleEndian, &sizeX)
			if err != nil {
				return 0, fmt.Errorf("Failed to read the %s chunk X-axis size. %v", chunkIDStr, err)
			}
			var sizeY uint32
			err = binary.Read(r, binary.LittleEndian, &sizeY)
			if err != nil {
				return 0, fmt.Errorf("Failed to read the %s chunk Y-axis size. %v", chunkIDStr, err)
			}
			var sizeZ uint32
			err = binary.Read(r, binary.LittleEndian, &sizeZ)
			if err != nil {
				return 0, fmt.Errorf("Failed to read the %s chunk Z-axis size. %v", chunkIDStr, err)
			}

			voxFile.SizeX = sizeX
			voxFile.SizeY = sizeY
			voxFile.SizeZ = sizeZ
		} else if chunkIDStr == "XYZI" {
			var voxelCount uint32
			err = binary.Read(r, binary.LittleEndian, &voxelCount)
			if err != nil {
				return 0, fmt.Errorf("Failed to read the %s chunk voxel count. %v", chunkIDStr, err)
			}

			// make the slice of Voxels for the number of voxels in the file
			voxels := make([]*Voxel, voxelCount)
			for i := uint32(0); i < voxelCount; i++ {
				var vX, vY, vZ, vI uint8
				err = binary.Read(r, binary.LittleEndian, &vX)
				if err != nil {
					return 0, fmt.Errorf("Failed to read the %s chunk voxel #%d. %v", chunkIDStr, i, err)
				}
				err = binary.Read(r, binary.LittleEndian, &vY)
				if err != nil {
					return 0, fmt.Errorf("Failed to read the %s chunk voxel #%d. %v", chunkIDStr, i, err)
				}
				err = binary.Read(r, binary.LittleEndian, &vZ)
				if err != nil {
					return 0, fmt.Errorf("Failed to read the %s chunk voxel #%d. %v", chunkIDStr, i, err)
				}
				err = binary.Read(r, binary.LittleEndian, &vI)
				if err != nil {
					return 0, fmt.Errorf("Failed to read the %s chunk voxel #%d. %v", chunkIDStr, i, err)
				}

				v := &Voxel{vX, vY, vZ, vI}
				voxels[i] = v
			}

			voxFile.Voxels = voxels
		} else if chunkIDStr == "RGBA" {
			// we have a fancy lad here with a fancy custom palette.
			// note: the last index isn't used but we'll read it anyway.
			const paletteSize = 256
			customPalette := make([]*Color, 256)
			for i := 0; i < paletteSize; i++ {
				var vR, vG, vB, vA uint8
				err = binary.Read(r, binary.LittleEndian, &vR)
				if err != nil {
					return 0, fmt.Errorf("Failed to read the %s chunk color #%d. %v", chunkIDStr, i, err)
				}
				err = binary.Read(r, binary.LittleEndian, &vG)
				if err != nil {
					return 0, fmt.Errorf("Failed to read the %s chunk color #%d. %v", chunkIDStr, i, err)
				}
				err = binary.Read(r, binary.LittleEndian, &vB)
				if err != nil {
					return 0, fmt.Errorf("Failed to read the %s chunk color #%d. %v", chunkIDStr, i, err)
				}
				err = binary.Read(r, binary.LittleEndian, &vA)
				if err != nil {
					return 0, fmt.Errorf("Failed to read the %s chunk color #%d. %v", chunkIDStr, i, err)
				}

				color := &Color{vR, vG, vB, vA}
				customPalette[i] = color
			}

			voxFile.Palette = customPalette
		} else {
			// this is a chunk that must have been deprecated in the new file format.
			// just read it in and ditch the data
			totalRead := 0
			contents := make([]byte, 256)
			_ = contents
			for totalRead < int(chunkSize) {
				c, err := r.Read(chunkID[:])
				if err != nil {
					return 0, fmt.Errorf("Failed to read the %s chunk contents. %v", chunkIDStr, err)
				}
				totalRead += c
			}
		}
	}

	// read children if necessary
	remainingBytes := chunkChildrenSize
	for remainingBytes > 0 {
		childReadSize, err := readChunk(r, voxFile)
		if err != nil {
			return 0, err
		}
		remainingBytes = remainingBytes - childReadSize
	}

	return chunkSize + 12, nil // +12 bytes for id, size, childSize
}

func instancePalette(p [256]uint32) []*Color {
	const paletteSize = 256
	palette := make([]*Color, 256)

	for i := 0; i < paletteSize; i++ {
		color := &Color{}

		value := p[i]
		color.R = uint8(value & 0x000000ff)
		color.G = uint8((value & 0x0000ff00) >> 8)
		color.B = uint8((value & 0x00ff0000) >> 16)
		color.A = uint8((value & 0xff000000) >> 24)

		palette[i] = color
	}

	return palette
}

var defaultPalette = [256]uint32{0x00000000, 0xffffffff, 0xffccffff, 0xff99ffff, 0xff66ffff, 0xff33ffff, 0xff00ffff, 0xffffccff, 0xffccccff, 0xff99ccff, 0xff66ccff, 0xff33ccff, 0xff00ccff, 0xffff99ff, 0xffcc99ff, 0xff9999ff,
	0xff6699ff, 0xff3399ff, 0xff0099ff, 0xffff66ff, 0xffcc66ff, 0xff9966ff, 0xff6666ff, 0xff3366ff, 0xff0066ff, 0xffff33ff, 0xffcc33ff, 0xff9933ff, 0xff6633ff, 0xff3333ff, 0xff0033ff, 0xffff00ff,
	0xffcc00ff, 0xff9900ff, 0xff6600ff, 0xff3300ff, 0xff0000ff, 0xffffffcc, 0xffccffcc, 0xff99ffcc, 0xff66ffcc, 0xff33ffcc, 0xff00ffcc, 0xffffcccc, 0xffcccccc, 0xff99cccc, 0xff66cccc, 0xff33cccc,
	0xff00cccc, 0xffff99cc, 0xffcc99cc, 0xff9999cc, 0xff6699cc, 0xff3399cc, 0xff0099cc, 0xffff66cc, 0xffcc66cc, 0xff9966cc, 0xff6666cc, 0xff3366cc, 0xff0066cc, 0xffff33cc, 0xffcc33cc, 0xff9933cc,
	0xff6633cc, 0xff3333cc, 0xff0033cc, 0xffff00cc, 0xffcc00cc, 0xff9900cc, 0xff6600cc, 0xff3300cc, 0xff0000cc, 0xffffff99, 0xffccff99, 0xff99ff99, 0xff66ff99, 0xff33ff99, 0xff00ff99, 0xffffcc99,
	0xffcccc99, 0xff99cc99, 0xff66cc99, 0xff33cc99, 0xff00cc99, 0xffff9999, 0xffcc9999, 0xff999999, 0xff669999, 0xff339999, 0xff009999, 0xffff6699, 0xffcc6699, 0xff996699, 0xff666699, 0xff336699,
	0xff006699, 0xffff3399, 0xffcc3399, 0xff993399, 0xff663399, 0xff333399, 0xff003399, 0xffff0099, 0xffcc0099, 0xff990099, 0xff660099, 0xff330099, 0xff000099, 0xffffff66, 0xffccff66, 0xff99ff66,
	0xff66ff66, 0xff33ff66, 0xff00ff66, 0xffffcc66, 0xffcccc66, 0xff99cc66, 0xff66cc66, 0xff33cc66, 0xff00cc66, 0xffff9966, 0xffcc9966, 0xff999966, 0xff669966, 0xff339966, 0xff009966, 0xffff6666,
	0xffcc6666, 0xff996666, 0xff666666, 0xff336666, 0xff006666, 0xffff3366, 0xffcc3366, 0xff993366, 0xff663366, 0xff333366, 0xff003366, 0xffff0066, 0xffcc0066, 0xff990066, 0xff660066, 0xff330066,
	0xff000066, 0xffffff33, 0xffccff33, 0xff99ff33, 0xff66ff33, 0xff33ff33, 0xff00ff33, 0xffffcc33, 0xffcccc33, 0xff99cc33, 0xff66cc33, 0xff33cc33, 0xff00cc33, 0xffff9933, 0xffcc9933, 0xff999933,
	0xff669933, 0xff339933, 0xff009933, 0xffff6633, 0xffcc6633, 0xff996633, 0xff666633, 0xff336633, 0xff006633, 0xffff3333, 0xffcc3333, 0xff993333, 0xff663333, 0xff333333, 0xff003333, 0xffff0033,
	0xffcc0033, 0xff990033, 0xff660033, 0xff330033, 0xff000033, 0xffffff00, 0xffccff00, 0xff99ff00, 0xff66ff00, 0xff33ff00, 0xff00ff00, 0xffffcc00, 0xffcccc00, 0xff99cc00, 0xff66cc00, 0xff33cc00,
	0xff00cc00, 0xffff9900, 0xffcc9900, 0xff999900, 0xff669900, 0xff339900, 0xff009900, 0xffff6600, 0xffcc6600, 0xff996600, 0xff666600, 0xff336600, 0xff006600, 0xffff3300, 0xffcc3300, 0xff993300,
	0xff663300, 0xff333300, 0xff003300, 0xffff0000, 0xffcc0000, 0xff990000, 0xff660000, 0xff330000, 0xff0000ee, 0xff0000dd, 0xff0000bb, 0xff0000aa, 0xff000088, 0xff000077, 0xff000055, 0xff000044,
	0xff000022, 0xff000011, 0xff00ee00, 0xff00dd00, 0xff00bb00, 0xff00aa00, 0xff008800, 0xff007700, 0xff005500, 0xff004400, 0xff002200, 0xff001100, 0xffee0000, 0xffdd0000, 0xffbb0000, 0xffaa0000,
	0xff880000, 0xff770000, 0xff550000, 0xff440000, 0xff220000, 0xff110000, 0xffeeeeee, 0xffdddddd, 0xffbbbbbb, 0xffaaaaaa, 0xff888888, 0xff777777, 0xff555555, 0xff444444, 0xff222222, 0xff111111,
}

package mld

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"fmt"
	"os"
)

type LoadSeg string

const (
	ExitSuccess = iota
	ExitFail
)

type SpecificSection string

// 特殊セクションのlist
const (
	DATA     SpecificSection = ".data"
	BSS      SpecificSection = ".bss"
	SYMTAB   SpecificSection = ".symtab"
	STRTAB   SpecificSection = ".strtab"
	SHSTRTAB SpecificSection = ".shstrtab"
)

func setSection() {}

// sectionを走査して、
// 	1 data, bssは消す
// 	2 sym, str系のsectionを書き換える
// などの特殊セクション(SectionType)に対するを行う.
func ReWriteSection(sections []*elf.Section, secfile *os.File) error {
	for _, sec := range sections {
		fmt.Printf("secname: %s  ", sec.Name)
		d, err := sec.Data()

		switch sec.Name {
		case ".data":
		case ".bss":
		case ".symtab":
			d = rewriteSection(".symtab", d)
			_, err = secfile.Write(d)
		case ".strtab":
			d = rewriteSection(".symtab", d)
			_, err = secfile.Write(d)
		case ".shstrtab":
			d = rewriteSection(".symtab", d)
			_, err = secfile.Write(d)
		default:
		}
		if err != nil {
			return err
		}
		fmt.Println()
	}
	return nil
}

func Run() int {
	fileName := os.Args[1]
	fmt.Println("filename:", fileName)

	f, err := elf.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERR: %s", err.Error())
		return ExitFail
	}

	secFile, err := os.Create("tmp_sec")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERR: %s", err.Error())
		return ExitFail
	}
	// defer os.Remove(tmpFile.Name())

	// 各sectionに対してloop, tmp_secに書き込む
	err = ReWriteSection(f.Sections, secFile)

	_, err = createProgramHeader(f.Sections)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERR: %s", err.Error())
		return ExitFail
	}

	proFile, err := os.Create("tmp_pro")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERR: %s", err.Error())
		return ExitFail
	}

	writeAll(secFile, proFile)

	return ExitSuccess
}

// wip
func writeAll(secfile, profile *os.File) {
	// write elf header

	// write program header

	// write each section

	// write section header table
}

var progTable = map[string]struct{}{".text": {}}

func getProgHdrNum(sections []*elf.Section) int {
	count := 0
	for _, v := range sections {
		if _, ok := progTable[v.Name]; ok {
			count++
		}
	}
	// for .note.gnu.build-id segment
	count++
	return count
}

// sectionを書き換える。
// 変更後のbyte数は、後にtmpに書き出す(Write())時の返り値で判断する
func rewriteSection(section SpecificSection, data []byte) []byte {
	switch section {
	case SYMTAB:

	case STRTAB:
	case SHSTRTAB:
	default:
	}
	return data
}

type SegmentType string

type SegEnt struct {
	segType      SegmentType    // segment type
	section      []*elf.Section // そのsegmentに属してるsection
	segEntMemsz  uint64
	segEntFilesz uint64
}

const (
	// .text .note.gnu.build-id
	TEXT SegmentType = ".text"
	// .note.gnu.build-id
	NOTE_GNU_BUILD_ID SegmentType = ".note.gnu.build-id"
)

// sectionの配列を受け取って、それに応じたSegEntの配列を返す
func getSegEnt(sections []*elf.Section) ([]*SegEnt, error) {
	// すでにSegEntを作成したsegmentはこのtableに入れておく.
	segTypeTable := make(map[SegmentType]*SegEnt)

	for _, v := range sections {
		switch v.Name {
		case ".text":
			if _, ok := segTypeTable[TEXT]; !ok {
				segTypeTable[TEXT] = &SegEnt{
					segType:      TEXT,
					section:      []*elf.Section{v},
					segEntMemsz:  v.Size,
					segEntFilesz: v.Size,
				}
				continue
			}
			segTypeTable[TEXT].section = append(segTypeTable[TEXT].section, v)
			segTypeTable[TEXT].segEntMemsz += v.Size
			segTypeTable[TEXT].segEntFilesz += v.Size
		case ".note.gnu.build-id":
			// NOTE_GNU_BUILD_IDに関して
			if _, ok := segTypeTable[NOTE_GNU_BUILD_ID]; !ok {
				segTypeTable[NOTE_GNU_BUILD_ID] = &SegEnt{
					segType:      NOTE_GNU_BUILD_ID,
					section:      []*elf.Section{v},
					segEntMemsz:  v.Size,
					segEntFilesz: v.Size,
				}
			} else {
				segTypeTable[NOTE_GNU_BUILD_ID].section = append(segTypeTable[NOTE_GNU_BUILD_ID].section, v)
				segTypeTable[NOTE_GNU_BUILD_ID].segEntMemsz += v.Size
				segTypeTable[NOTE_GNU_BUILD_ID].segEntFilesz += v.Size
			}
			// TEXTに関して
			if _, ok := segTypeTable[TEXT]; !ok {
				segTypeTable[TEXT] = &SegEnt{
					segType:      TEXT,
					section:      []*elf.Section{v},
					segEntMemsz:  v.Size,
					segEntFilesz: v.Size,
				}
				continue
			}
			segTypeTable[TEXT].section = append(segTypeTable[TEXT].section, v)
			segTypeTable[TEXT].segEntMemsz += v.Size
			segTypeTable[TEXT].segEntFilesz += v.Size
		}
	}

	segEnt := []*SegEnt{}
	for _, v := range segTypeTable {
		segEnt = append(segEnt, v)
	}

	return segEnt, nil
}

// sectionを走査して、load可能sectionに対しては
// それに対応するprogram headerを作成して返す
func createProgramHeader(sections []*elf.Section) ([]*ProgHeader, error) {
	// TODO: hardcode
	var addrOff uint64 = 0x400000

	// sectionを一通り走査して、segEntの配列を得る
	_, err := getSegEnt(sections)
	if err != nil {
		return nil, err
	}

	for _, v := range sections {
		switch v.Name {
		case ".text":
			// textはload先segmentの中でも、先頭に配置する
			p := ProgHeader{
				Type:  PT_LOAD,
				Flags: PF_R | PF_X,
				Off:   0x0,
				Vaddr: addrOff,
				Paddr: addrOff,
				// Filesz: uint64(progNum) * 56, // TODO: hard
				// // Memsz:  uint64(progNum) * 56, // TODO: hard
				Align: 0x200000, // TODO: hard
			}
			buf := new(bytes.Buffer)
			// bufにbinaryを書く
			err := binary.Write(buf, binary.LittleEndian, p)
			if err != nil {
				return nil, err
			}
			// bufのbinaryをそのままproFileに書く
			// _, err = proFile.WriteAt(buf.Bytes(), 0)
			// if err != nil {
			// 	return nil, err
			// }
		}
	}
	return nil, nil
}

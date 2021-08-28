package mld

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"fmt"
	"os"
)

const (
	ExitSuccess = iota
	ExitFail
)

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
	for _, sec := range f.Sections {
		fmt.Printf("secname: %s  ", sec.Name)
		d, err := sec.Data()

		switch sec.Name {
		case ".data":
		case ".bss":
		case ".symtab":
			d = rewriteSection(".symtab", d)
			_, err = secFile.Write(d)
		case ".strtab":
			d = rewriteSection(".symtab", d)
			_, err = secFile.Write(d)
		case ".shstrtab":
			d = rewriteSection(".symtab", d)
			_, err = secFile.Write(d)
		default:
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERR: %s", err.Error())
			return ExitFail
		}
		fmt.Println()
	}

	// program headerの数を予め割り出しておく必要がある
	progNum := getProgHdrNum(f.Sections)

	proFile, err := os.Create("tmp_pro")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERR: %s", err.Error())
		return ExitFail
	}
	err = createProgramHeader(f.Sections, proFile, progNum)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERR: %s", err.Error())
		return ExitFail
	}

	return ExitSuccess
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

// TODO
func rewriteSection(section string, data []byte) []byte {
	switch section {
	case ".symtab":
	case ".strtab":
	case ".shstrtab":
	default:
	}
	return data
}

// sectionを走査して、load可能sectionに対しては
// それに対応するprogram headerを作成する
func createProgramHeader(sections []*elf.Section, proFile *os.File, progNum int) error {
	// TODO: hardcode
	var addrOff uint64 = 0x400000

	for _, v := range sections {
		switch v.Name {
		case ".text":
			// textはload先segmentの中でも、先頭に配置する
			p := ProgHeader{
				Type:   PT_LOAD,
				Flags:  PF_R | PF_X,
				Off:    0x0,
				Vaddr:  addrOff,
				Paddr:  addrOff,
				Filesz: uint64(progNum) * 56, // TODO: hard
				Memsz:  uint64(progNum) * 56, // TODO: hard
				Align:  0x200000,             // TODO: hard
			}
			buf := new(bytes.Buffer)
			err := binary.Write(buf, binary.LittleEndian, p)
			if err != nil {
				return err
			}
			proFile.WriteAt(buf.Bytes(), 0)
		}
	}
	return nil
}

package mld

type Elf64_Addr = uint64
type Elf64_Off = uint64
type Elf64_Half = uint16
type Elf64_Word = uint32
type Elf64_Sword = int32
type Elf64_Xword = uint64
type Elf64_Sxword = int64

// 4+1+1+2+8+8=24byte
type Elf64_Sym struct {
	st_name  Elf64_Word  // Symbol name (index into string table)
	st_info  uint8       // Symbol's type and binding attributes
	st_other uint8       // Must be zero; reserved
	st_shndx Elf64_Half  // Which section (header tbl index) it's defined in
	st_value Elf64_Addr  // Value or address associated with the symbol
	st_size  Elf64_Xword // Size of the symbol
}

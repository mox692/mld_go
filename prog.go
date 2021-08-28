package mld

// A ProgHeader represents a single ELF program header.
type ProgHeader struct {
	Type   ProgType
	Flags  ProgFlag
	Off    uint64
	Vaddr  uint64
	Paddr  uint64
	Filesz uint64
	Memsz  uint64
	Align  uint64
}

type ProgFlag uint32

const (
	PF_X        ProgFlag = 0x1        /* Executable. */
	PF_W        ProgFlag = 0x2        /* Writable. */
	PF_R        ProgFlag = 0x4        /* Readable. */
	PF_MASKOS   ProgFlag = 0x0ff00000 /* Operating system-specific. */
	PF_MASKPROC ProgFlag = 0xf0000000 /* Processor-specific. */
)

// for Linux
type ProgType int32

const (
	PT_NULL    ProgType = 0          /* Unused entry. */
	PT_LOAD    ProgType = 1          /* Loadable segment. */
	PT_DYNAMIC ProgType = 2          /* Dynamic linking information segment. */
	PT_INTERP  ProgType = 3          /* Pathname of interpreter. */
	PT_NOTE    ProgType = 4          /* Auxiliary information. */
	PT_SHLIB   ProgType = 5          /* Reserved (not used). */
	PT_PHDR    ProgType = 6          /* Location of program header itself. */
	PT_TLS     ProgType = 7          /* Thread local storage segment */
	PT_LOOS    ProgType = 0x60000000 /* First OS-specific. */
	PT_HIOS    ProgType = 0x6fffffff /* Last OS-specific. */
	PT_LOPROC  ProgType = 0x70000000 /* First processor-specific type. */
	PT_HIPROC  ProgType = 0x7fffffff /* Last processor-specific type. */
)

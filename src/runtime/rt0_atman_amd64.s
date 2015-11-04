#define _PAGE_ROUND_UP(REGISTER)		\
	ADDQ	$0x0000000000000fff, REGISTER	\
	ANDQ	$0xfffffffffffff000, REGISTER

#define CALL_RBX                \
	BYTE $0xff; BYTE $0xd3	// callq *%rbx

#define CRASH_ON_NONZERO	\
	CMPQ	AX, $-1		\
	JNE	2(PC)		\
	ANDQ	$0xdeadbeef, 0xdeadbeef

#define _HYPERCALL(OFFSET)				\
	MOVQ	$runtime·_atman_hypercall_page(SB), BX	\
	_PAGE_ROUND_UP(BX)				\
	ADDQ	OFFSET, BX				\
	CALL_RBX                                        \
	CRASH_ON_NONZERO

#define _HYPERVISOR_console_io(OP, SIZE, DATA_PTR) \
	MOVQ	OP, DI		\
	MOVQ	SIZE, SI	\
	MOVQ	DATA_PTR, DX	\
	_HYPERCALL($0x240)

TEXT _rt0_amd64_atman(SB),NOSPLIT,$-8
	CLD
	MOVQ	$runtime·_atman_stack+0x8000(SB), SP
	MOVQ	SI, runtime·_atman_start_info(SB)

	_HYPERVISOR_console_io($0, $7, $runtime·_atman_hello(SB))

load_phys_to_machine_mapping:
	MOVQ	runtime·_atman_start_info(SB), AX
	ADDQ	$104, (AX)
	MOVQ	AX, runtime·_atman_phys_to_machine_mapping(SB)

	MOVQ	$main(SB), AX
	JMP	AX

TEXT main(SB),NOSPLIT,$-8
	MOVQ	$runtime·rt0_go(SB), AX
	JMP	AX

DATA runtime·_atman_hello(SB)/8, $"hello\n"
GLOBL runtime·_atman_hello(SB), NOPTR, $8

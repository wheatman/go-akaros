// Created by cgo -cdefs - DO NOT EDIT
// cgo -cdefs types_akaros.go


enum {
	FUTEX_WAIT	= 0x0,
	FUTEX_WAKE	= 0x1,
};
enum {
	UINFO = 0x7f800000,
};

typedef struct Vcore Vcore;
typedef struct Pcore Pcore;
typedef struct Procinfo Procinfo;
typedef struct Timespec Timespec;
typedef struct Timeval Timeval;

#pragma pack on

struct Vcore {
	byte	*dummy_ptr1;
	byte	*dummy_ptr2;
	uint32	pcoreid;
	bool	valid;
	byte	Pad_cgo_0[3];
	uint32	nr_preempts_sent;
	uint32	nr_preempts_done;
	uint64	preempt_pending;
};
struct Pcore {
	uint32	vcoreid;
	bool	valid;
	byte	Pad_cgo_0[3];
};
struct Procinfo {
	int32	pid;
	int32	ppid;
	uint32	max_vcores;
	uint64	tsc_freq;
	uint64	timing_overhead;
	byte	*heap_bottom;
	int8	*argp[32];
	int8	argbuf[3072];
	bool	is_mcp;
	byte	Pad_cgo_0[3];
	uint32	res_grant[3];
	Vcore	vcoremap[255];
	uint32	num_vcores;
	Pcore	pcoremap[255];
	uint8	coremap_seqctr;
	byte	Pad_cgo_1[3];
};
struct Timespec {
	int32	tv_sec;
	int32	tv_nsec;
};
struct Timeval {
	int32	tv_sec;
	int32	tv_usec;
};


#pragma pack off

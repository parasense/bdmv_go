package mobj

// Branch commands, when CommandGroup == 0
const (
	// CommandSubGroup == 0x00
	MOBJ_BRANCH_OPTION_SUB0_NOP   uint8 = 0x00 // BranchOption == 0b0000
	MOBJ_BRANCH_OPTION_SUB0_GOTO  uint8 = 0x01 // BranchOption == 0b0001
	MOBJ_BRANCH_OPTION_SUB0_BREAK uint8 = 0x02 // BranchOption == 0b0010

	// CommandSubGroup == 0x01
	MOBJ_BRANCH_OPTION_SUB1_JUMP_OBJECT uint8 = 0x00 // BranchOption == 0b0000
	MOBJ_BRANCH_OPTION_SUB1_JUMP_TITLE  uint8 = 0x01 // BranchOption == 0b0001
	MOBJ_BRANCH_OPTION_SUB1_CALL_OBJECT uint8 = 0x02 // BranchOption == 0b0010
	MOBJ_BRANCH_OPTION_SUB1_CALL_TITLE  uint8 = 0x03 // BranchOption == 0b0011
	MOBJ_BRANCH_OPTION_SUB1_RESUME      uint8 = 0x04 // BranchOption == 0b0100

	// CommandSubGroup == 0x02
	MOBJ_BRANCH_OPTION_SUB2_PLAYLIST  uint8 = 0x00 // BranchOption == 0b0000
	MOBJ_BRANCH_OPTION_SUB2_PLAYITEM  uint8 = 0x01 // BranchOption == 0b0001
	MOBJ_BRANCH_OPTION_SUB2_PLAYMARK  uint8 = 0x02 // BranchOption == 0b0010
	MOBJ_BRANCH_OPTION_SUB2_TERMINATE uint8 = 0x03 // BranchOption == 0b0011
	MOBJ_BRANCH_OPTION_SUB2_LINKITEM  uint8 = 0x04 // BranchOption == 0b0100
	MOBJ_BRANCH_OPTION_SUB2_LINKMARK  uint8 = 0x05 // BranchOption == 0b0101
)

// Compare commands, when CommandGroup == 1
const (
	MOBJ_COMPARE_OPTION_BC uint8 = 0x01 // CompareOption == 0b0001
	MOBJ_COMPARE_OPTION_EQ uint8 = 0x02 // CompareOption == 0b0010
	MOBJ_COMPARE_OPTION_NE uint8 = 0x03 // CompareOption == 0b0011
	MOBJ_COMPARE_OPTION_GE uint8 = 0x04 // CompareOption == 0b0100
	MOBJ_COMPARE_OPTION_GT uint8 = 0x05 // CompareOption == 0b0101
	MOBJ_COMPARE_OPTION_LE uint8 = 0x06 // CompareOption == 0b0110
	MOBJ_COMPARE_OPTION_LT uint8 = 0x07 // CompareOption == 0b0111
)

// Set commands, when CommandGroup == 2
const (
	// CommandSubGroup == 0x00
	MOBJ_SET_OPTION_SUB0_MOVE       uint8 = 0x01 // SetOption == 0b00001
	MOBJ_SET_OPTION_SUB0_SWAP       uint8 = 0x02 // SetOption == 0b00010
	MOBJ_SET_OPTION_SUB0_ADD        uint8 = 0x03 // SetOption == 0b00011
	MOBJ_SET_OPTION_SUB0_SUB        uint8 = 0x04 // SetOption == 0b00100
	MOBJ_SET_OPTION_SUB0_MUL        uint8 = 0x05 // SetOption == 0b00101
	MOBJ_SET_OPTION_SUB0_DIV        uint8 = 0x06 // SetOption == 0b00110
	MOBJ_SET_OPTION_SUB0_MOD        uint8 = 0x07 // SetOption == 0b00111
	MOBJ_SET_OPTION_SUB0_RND        uint8 = 0x08 // SetOption == 0b01000
	MOBJ_SET_OPTION_SUB0_AND        uint8 = 0x09 // SetOption == 0b01001
	MOBJ_SET_OPTION_SUB0_OR         uint8 = 0x0A // SetOption == 0b01010
	MOBJ_SET_OPTION_SUB0_XOR        uint8 = 0x0B // SetOption == 0b01011
	MOBJ_SET_OPTION_SUB0_BITSET     uint8 = 0x0C // SetOption == 0b01100
	MOBJ_SET_OPTION_SUB0_BITCLR     uint8 = 0x0D // SetOption == 0b01101
	MOBJ_SET_OPTION_SUB0_SHIFTLEFT  uint8 = 0x0E // SetOption == 0b01110
	MOBJ_SET_OPTION_SUB0_SHIFTRIGHT uint8 = 0x0F // SetOption == 0b01111

	// Set System?
	// CommandSubGroup == 0x01
	MOBJ_SET_OPTION_SUB1_SETSTREAM          uint8 = 0x01 // SetOption == 0b00001
	MOBJ_SET_OPTION_SUB1_SETNVTIMER         uint8 = 0x02 // SetOption == 0b00010
	MOBJ_SET_OPTION_SUB1_BUTTONPAGE         uint8 = 0x03 // SetOption == 0b00011
	MOBJ_SET_OPTION_SUB1_ENABLEBUTTON       uint8 = 0x04 // SetOption == 0b00100
	MOBJ_SET_OPTION_SUB1_DISABLEBUTTON      uint8 = 0x05 // SetOption == 0b00101
	MOBJ_SET_OPTION_SUB1_SETSECONDARYSTREAM uint8 = 0x06 // SetOption == 0b00110
	MOBJ_SET_OPTION_SUB1_POPUPOFF           uint8 = 0x07 // SetOption == 0b00111
	MOBJ_SET_OPTION_SUB1_STILLON            uint8 = 0x08 // SetOption == 0b01000
	MOBJ_SET_OPTION_SUB1_STILLOFF           uint8 = 0x09 // SetOption == 0b01001
	MOBJ_SET_OPTION_SUB1_OUTPUTMODE         uint8 = 0x0A // SetOption == 0b01010 // XXX - 3d option?
	MOBJ_SET_OPTION_SUB1_STREAMSS           uint8 = 0x0B // SetOption == 0b01011 // XXX - 3d option?
)

func GetCommand(cmdGrp, cmdSubGrp, branchOpt, compareOpt, setOpt uint8) (cmd string) {

	switch {
	case cmdGrp == 0 && cmdSubGrp == 0:
		switch branchOpt {
		case MOBJ_BRANCH_OPTION_SUB0_NOP:
			cmd = "NOP"
		case MOBJ_BRANCH_OPTION_SUB0_GOTO:
			cmd = "GOTO"
		case MOBJ_BRANCH_OPTION_SUB0_BREAK:
			cmd = "BREAK"
		}
	case cmdGrp == 0 && cmdSubGrp == 1:
		switch branchOpt {
		case MOBJ_BRANCH_OPTION_SUB1_JUMP_OBJECT:
			cmd = "JUMP OBJ"
		case MOBJ_BRANCH_OPTION_SUB1_JUMP_TITLE:
			cmd = "JUMP TITLE"
		case MOBJ_BRANCH_OPTION_SUB1_CALL_OBJECT:
			cmd = "CALL OBJECT"
		case MOBJ_BRANCH_OPTION_SUB1_CALL_TITLE:
			cmd = "CALL TITLE"
		case MOBJ_BRANCH_OPTION_SUB1_RESUME:
			cmd = "RESUME"
		}
	case cmdGrp == 0 && cmdSubGrp == 2:
		switch branchOpt {
		case MOBJ_BRANCH_OPTION_SUB2_PLAYLIST:
			cmd = "PLAY LIST"
		case MOBJ_BRANCH_OPTION_SUB2_PLAYITEM:
			cmd = "PLAY ITEM"
		case MOBJ_BRANCH_OPTION_SUB2_PLAYMARK:
			cmd = "PLAY MARK"
		case MOBJ_BRANCH_OPTION_SUB2_TERMINATE:
			cmd = "TERMINATE"
		case MOBJ_BRANCH_OPTION_SUB2_LINKITEM:
			cmd = "LINK ITEM"
		case MOBJ_BRANCH_OPTION_SUB2_LINKMARK:
			cmd = "LINK MARK"
		}
	case cmdGrp == 1 && cmdSubGrp == 0:
		switch compareOpt {
		case MOBJ_COMPARE_OPTION_BC:
			cmd = "BC"
		case MOBJ_COMPARE_OPTION_EQ:
			cmd = "EQ"
		case MOBJ_COMPARE_OPTION_NE:
			cmd = "NE"
		case MOBJ_COMPARE_OPTION_GE:
			cmd = "GE"
		case MOBJ_COMPARE_OPTION_GT:
			cmd = "GT"
		case MOBJ_COMPARE_OPTION_LE:
			cmd = "LE"
		case MOBJ_COMPARE_OPTION_LT:
			cmd = "LT"
		}
	case cmdGrp == 2 && cmdSubGrp == 0:
		switch setOpt {
		case MOBJ_SET_OPTION_SUB0_MOVE:
			cmd = "MOVE"
		case MOBJ_SET_OPTION_SUB0_SWAP:
			cmd = "SWAP"
		case MOBJ_SET_OPTION_SUB0_ADD:
			cmd = "ADD"
		case MOBJ_SET_OPTION_SUB0_SUB:
			cmd = "SUB"
		case MOBJ_SET_OPTION_SUB0_MUL:
			cmd = "MUL"
		case MOBJ_SET_OPTION_SUB0_DIV:
			cmd = "DIV"
		case MOBJ_SET_OPTION_SUB0_MOD:
			cmd = "MOD"
		case MOBJ_SET_OPTION_SUB0_RND:
			cmd = "RND"
		case MOBJ_SET_OPTION_SUB0_AND:
			cmd = "AND"
		case MOBJ_SET_OPTION_SUB0_OR:
			cmd = "OR"
		case MOBJ_SET_OPTION_SUB0_XOR:
			cmd = "XOR"
		case MOBJ_SET_OPTION_SUB0_BITSET:
			cmd = "BIT SET"
		case MOBJ_SET_OPTION_SUB0_BITCLR:
			cmd = "BIT CLEAR"
		case MOBJ_SET_OPTION_SUB0_SHIFTLEFT:
			cmd = "SHIFT LEFT"
		case MOBJ_SET_OPTION_SUB0_SHIFTRIGHT:
			cmd = "SHIFT RIGHT"
		}
	case cmdGrp == 2 && cmdSubGrp == 1:
		switch setOpt {
		case MOBJ_SET_OPTION_SUB1_SETSTREAM:
			cmd = "SET STREAM"
		case MOBJ_SET_OPTION_SUB1_SETNVTIMER:
			cmd = "SET NV TIMER"
		case MOBJ_SET_OPTION_SUB1_BUTTONPAGE:
			cmd = "SET BUTTON PAGE"
		case MOBJ_SET_OPTION_SUB1_ENABLEBUTTON:
			cmd = "SET ENABLE BUTTON"
		case MOBJ_SET_OPTION_SUB1_DISABLEBUTTON:
			cmd = "SET DISABLE BUTTON"
		case MOBJ_SET_OPTION_SUB1_SETSECONDARYSTREAM:
			cmd = "SET SECONDARY STREAM"
		case MOBJ_SET_OPTION_SUB1_POPUPOFF:
			cmd = "SET POPUP OFF"
		case MOBJ_SET_OPTION_SUB1_STILLON:
			cmd = "SET STILL ON"
		case MOBJ_SET_OPTION_SUB1_STILLOFF:
			cmd = "SET STILL OFF"
		case MOBJ_SET_OPTION_SUB1_OUTPUTMODE:
			cmd = "SET OUTPUT MODE"
		case MOBJ_SET_OPTION_SUB1_STREAMSS:
			cmd = "SET STREAM SS"
		}
	}
	return cmd
}

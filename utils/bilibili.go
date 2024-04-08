package utils

import "liveroom"

func BilibiliGuardLevelToPrivilege(level int) int {
	switch level {
	case 1:
		return liveroom.PrivilegeUltimate
	case 2:
		return liveroom.PrivilegeAdvanced
	case 3:
		return liveroom.PrivilegeBasic
	default:
		return liveroom.PrivilegeNone
	}
}

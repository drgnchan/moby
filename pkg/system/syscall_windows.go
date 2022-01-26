package system // import "github.com/docker/docker/pkg/system"

import (
	"unsafe"

	"github.com/sirupsen/logrus"
)

const (
	// Deprecated: use github.com/docker/pkg/idtools.SeTakeOwnershipPrivilege
	SeTakeOwnershipPrivilege = "SeTakeOwnershipPrivilege"
	// Deprecated: use github.com/docker/pkg/idtools.ContainerAdministratorSidString
	ContainerAdministratorSidString = "S-1-5-93-2-1"
	// Deprecated: use github.com/docker/pkg/idtools.ContainerUserSidString
	ContainerUserSidString = "S-1-5-93-2-2"
)

var procGetVersionExW = modkernel32.NewProc("GetVersionExW")

// https://docs.microsoft.com/en-us/windows/win32/api/winnt/ns-winnt-osversioninfoexa
// TODO: use golang.org/x/sys/windows.OsVersionInfoEx (needs OSVersionInfoSize to be exported)
type osVersionInfoEx struct {
	OSVersionInfoSize uint32
	MajorVersion      uint32
	MinorVersion      uint32
	BuildNumber       uint32
	PlatformID        uint32
	CSDVersion        [128]uint16
	ServicePackMajor  uint16
	ServicePackMinor  uint16
	SuiteMask         uint16
	ProductType       byte
	Reserve           byte
}

// IsWindowsClient returns true if the SKU is client. It returns false on
// Windows server, or if an error occurred when making the GetVersionExW
// syscall.
func IsWindowsClient() bool {
	osviex := &osVersionInfoEx{OSVersionInfoSize: 284}
	r1, _, err := procGetVersionExW.Call(uintptr(unsafe.Pointer(osviex)))
	if r1 == 0 {
		logrus.WithError(err).Warn("GetVersionExW failed - assuming server SKU")
		return false
	}
	// VER_NT_WORKSTATION, see https://docs.microsoft.com/en-us/windows/win32/api/winnt/ns-winnt-osversioninfoexa
	const verNTWorkstation = 0x00000001 // VER_NT_WORKSTATION
	return osviex.ProductType == verNTWorkstation
}

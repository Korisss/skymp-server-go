package skymp_wrapper

import (
	"errors"
	"reflect"
	"syscall"
	"unsafe"
)

var (
	scampLib, _                   = syscall.LoadLibrary("scamp_lib")
	createServerProc, _           = syscall.GetProcAddress(scampLib, "CreateServer")
	getUserActorProc, _           = syscall.GetProcAddress(scampLib, "GetUserActor")
	setEnabledProc, _             = syscall.GetProcAddress(scampLib, "SetEnabled")
	setConnectHandlerProc, _      = syscall.GetProcAddress(scampLib, "SetConnectHandler")
	setDisconnectHandlerProc, _   = syscall.GetProcAddress(scampLib, "SetDisconnectHandler")
	setCustomPacketHandlerProc, _ = syscall.GetProcAddress(scampLib, "SetCustomPacketHandler")
	getActorsByProfileIdProc, _   = syscall.GetProcAddress(scampLib, "GetActorsByProfileId")
	setUserActorProc, _           = syscall.GetProcAddress(scampLib, "SetUserActor")
	createActorProc, _            = syscall.GetProcAddress(scampLib, "CreateActor")
	setRaceMenuOpenProc, _        = syscall.GetProcAddress(scampLib, "SetRaceMenuOpen")
)

type ScampServer struct {
	ptr uintptr
}

func Free() {
	syscall.FreeLibrary(scampLib)
}

func CreateServer(port uint16, maxPlayers uint16) ScampServer {
	server, _, _ := syscall.SyscallN(uintptr(createServerProc), uintptr(port), uintptr(maxPlayers))
	return ScampServer{
		ptr: server,
	}
}

func (server ScampServer) GetUserActor(userId UserId) (ActorFormId, error) {
	formId, _, _ := syscall.SyscallN(uintptr(getUserActorProc), server.ptr, uintptr(userId))
	if formId <= 0 {
		return 0, errors.New("actor not found")
	}

	return uint32(formId), nil
}

func (server ScampServer) GetActorsByProfileId(profileId ProfileId) []ActorFormId {
	len := uintptr(0)
	data, _, _ := syscall.SyscallN(uintptr(getActorsByProfileIdProc), server.ptr, uintptr(profileId), uintptr(unsafe.Pointer(&len)))

	var slice []ActorFormId
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Cap = int(len)
	header.Len = int(len)
	header.Data = data

	return slice
}

func (server ScampServer) CreateActor(formId FormId, pos Position, angleZ Angle, cellOrWorld WorldOrCell, profileId ProfileId) ActorFormId {
	actorId, _, _ := syscall.SyscallN(uintptr(createActorProc), server.ptr, uintptr(formId), uintptr(unsafe.Pointer(&pos)), uintptr(angleZ), uintptr(cellOrWorld), uintptr(profileId))
	return ActorFormId(actorId)
}

func (server ScampServer) SetRaceMenuOpen(actorId ActorFormId, open bool) {
	syscall.SyscallN(uintptr(setRaceMenuOpenProc), server.ptr, uintptr(actorId), uintptr(unsafe.Pointer(&open)))
}

func (server ScampServer) SetUserActor(userId UserId, actorId ActorFormId) {
	syscall.SyscallN(uintptr(setUserActorProc), server.ptr, uintptr(userId), uintptr(actorId))
}

func (server ScampServer) SetEnabled(actorId ActorFormId, enabled bool) {
	syscall.SyscallN(uintptr(setEnabledProc), server.ptr, uintptr(actorId), uintptr(unsafe.Pointer(&enabled)))
}

func (server ScampServer) SetConnectHandler(handler OnConnectHandler) {
	syscall.SyscallN(uintptr(setConnectHandlerProc), server.ptr, uintptr(unsafe.Pointer(&handler)))
}

func (server ScampServer) SetDisconnectHandler(handler OnDisonnectHandler) {
	syscall.SyscallN(uintptr(setDisconnectHandlerProc), server.ptr, uintptr(unsafe.Pointer(&handler)))
}

func (server ScampServer) SetCustomPacketHandler(handler OnCustomPacketHandler) {
	syscall.SyscallN(uintptr(setCustomPacketHandlerProc), server.ptr, uintptr(unsafe.Pointer(&handler)))
}

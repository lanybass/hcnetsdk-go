package hcnetsdk

// #include <stdlib.h>
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
	"syscall"
)

var DLL = syscall.MustLoadDLL("HCNetSDK.dll")

const (
	SERIALNO_LEN = 48
	STREAM_ID_LEN = 32
)

const (
 NET_DVR_SYSHEAD = 1    //系统头数据
 NET_DVR_STREAMDATA = 2    //视频流数据（包括复合流和音视频分开的视频流数据）
 NET_DVR_AUDIOSTREAMDATA = 3    //音频流数据
 NET_DVR_STD_VIDEODATA = 4    //标准视频流数据
 NET_DVR_STD_AUDIODATA = 5    //标准音频流数据
 NET_DVR_SDP = 6   //SDP信息(Rstp传输时有效)
 NET_DVR_CHANGE_FORWARD = 10  //码流改变为正放  
 NET_DVR_CHANGE_REVERSE = 11  //码流改变为倒放
 NET_DVR_PLAYBACK_ALLFILEEND = 12  //回放文件结束标记
 NET_DVR_VOD_DRAW_FRAME = 13  //回放抽帧码流
 NET_DVR_VOD_DRAW_DATA = 14  //拖动平滑码流
 NET_DVR_PRIVATE_DATA = 112 //私有数据,包括智能信息
)

type BYTE byte
type LONG int32
type WORD uint16
type DWORD uint32
type HWND uint32

type HCNetSDK struct {
	UserId LONG
	Info *DeviceInfo
}

type DeviceInfo struct {
	SerialNumber [SERIALNO_LEN]byte
	AlarmInPortNum BYTE
	AlarmOutPortNum BYTE
	DiskNum BYTE
	DVRType BYTE
	ChanNum BYTE
	StartChan BYTE
	AudioChanNum BYTE
	IPChanNum BYTE
	ZeroChanNum BYTE
	MainProto BYTE
	SubProto BYTE
	Support BYTE
	Support1 BYTE
	Support2 BYTE
	DevType WORD
	Support3 BYTE
	MultiStreamProto BYTE
	StartDChan BYTE
	StartDTalkChan BYTE
	HighDChanNum BYTE
	Support4 BYTE
	LanguageType BYTE
	VoiceInChanNum BYTE
	StartVoiceInChanNo BYTE
	Support5 BYTE
	Support6 BYTE
	MirrorChanNum BYTE
	StartMirrorChanNo WORD
	Support7 BYTE
	Res2 BYTE
}

type PreviewInfo struct {
	lChannel LONG	//通道号
	StreamType DWORD	// 码流类型，0-主码流，1-子码流，2-码流3，3-码流4, 4-码流5,5-码流6,7-码流7,8-码流8,9-码流9,10-码流10
	LinkMode DWORD	// 0：TCP方式,1：UDP方式,2：多播方式,3 - RTP方式，4-RTP/RTSP,5-RSTP/HTTP ,6- HRUDP（可靠传输） ,7-RTSP/HTTPS
	PlayWnd HWND	//播放窗口的句柄,为NULL表示不播放图象
	Blocked DWORD	//0-非阻塞取流, 1-阻塞取流, 如果阻塞SDK内部connect失败将会有5s的超时才能够返回,不适合于轮询取流操作.
	PassbackRecord DWORD	//0-不启用录像回传,1启用录像回传
	PreviewMode BYTE	//预览模式，0-正常预览，1-延迟预览
	StreamID [STREAM_ID_LEN]byte	//流ID，lChannel为0xffffffff时启用此参数
	ProtoType BYTE	//应用层取流协议，0-私有协议，1-RTSP协议
	Res1 BYTE;
	VideoCodingType BYTE	//码流数据编解码类型 0-通用编码数据 1-热成像探测器产生的原始数据（温度数据的加密信息，通过去加密运算，将原始数据算出真实的温度值）
	DisplayBufNum DWORD	//播放库播放缓冲区最大缓冲帧数，范围1-50，置0时默认为1 
	NPQMode BYTE	//NPQ是直连模式，还是过流媒体 0-直连 1-过流媒体
	Res [215]byte
}

func (sdk *HCNetSDK) Init() bool {
	proc := DLL.MustFindProc("NET_DVR_Init")
	r, _, _ := proc.Call()
	if int(r) == 0 {
		return false
	}
	return true
}

func (sdk *HCNetSDK) Login(ipAddr string, port int, username, password string) error {
	proc := DLL.MustFindProc("NET_DVR_Login_V30")
	info := &DeviceInfo{}
	ip := unsafe.Pointer(C.CString(ipAddr))
	u := unsafe.Pointer(C.CString(username))
	p := unsafe.Pointer(C.CString(password))
	defer C.free(ip)
	defer C.free(u)
	defer C.free(p)
	userId, _, _ := proc.Call(
		uintptr(ip),
		uintptr(port),
		uintptr(u),
		uintptr(p),
		uintptr(unsafe.Pointer(info)))
	if sdk.UserId = LONG(userId); sdk.UserId < 0 {
		return errors.New(fmt.Sprintf("Login Error: %v\n", sdk.GetLastError()))
	}
	sdk.Info = info
	return nil
}

type RealDataCallBack func(playHandle LONG, dataType DWORD, buffer []byte, bufferSize DWORD)

func (sdk *HCNetSDK) RealPlay(info *PreviewInfo, callback RealDataCallBack) (handle LONG, err error) {
	proc := DLL.MustFindProc("NET_DVR_RealPlay_V40")
	realPlayHandle, _, _ := proc.Call(
		uintptr(sdk.UserId),
		uintptr(unsafe.Pointer(&info)),
		uintptr(unsafe.Pointer(&callback)),
		uintptr(0))
	if handle = LONG(realPlayHandle); handle < 0 {
		return -1, errors.New(fmt.Sprintf("RealPlay Error: %v\n", handle))
	}
	return handle, nil
}

func (sdk *HCNetSDK) GetLastError() LONG {
	proc := DLL.MustFindProc("NET_DVR_GetLastError")
	err, _, _ := proc.Call()
	return LONG(err)
}

func (sdk *HCNetSDK) Logout() bool {
	proc := DLL.MustFindProc("NET_DVR_Logout")
	r, _, _ := proc.Call(uintptr(sdk.UserId))
	if int(r) == 0 {
		return false
	}
	return true
}

func (sdk *HCNetSDK) Cleanup() bool {
	proc := DLL.MustFindProc("NET_DVR_Cleanup")
	r, _, _ := proc.Call()
	if int(r) == 0 {
		return false
	}
	return true
}

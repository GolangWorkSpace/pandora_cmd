package main

import (
	"errors"
	"bytes"
	"strconv"
)

func ReferSelect() (*ReferVersionModel, error) {
	refers, err := ReferRequestList()
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	buffer.WriteString("请选择参照(用于跟进):\n")
	for idx, aRefer := range refers {
		buffer.WriteString(strconv.Itoa(idx+1)+"."+aRefer.ReferName+"("+aRefer.ReferVersion + ") ")
		buffer.WriteString("内部版本:"+strconv.Itoa(aRefer.VersionId)+" 生成时间:"+aRefer.TimeString+"\n")
	}
	buffer.WriteString(":")
	selected := SimpleInputSelectNum(buffer.String(), 1, len(refers))
	return refers[selected], nil
}

func ReferRequestList() ([]*ReferVersionModel, error) {
	var referRes *ReferResp
	url := GenURL("/api/ref/version/list")
	if err := GETParse(url, &referRes); err != nil {
		return nil, errors.New(FormtError("查询参照列表失败", err))
	}
	if referRes.Errno != 0 {
		return nil, errors.New(FormatResError("查询参照列表失败", referRes.Msg, referRes.Errno))
	}
	if len(referRes.Refers) == 0 {
		return nil, errors.New("没有任何参照！")
	}
	return referRes.Refers, nil
}

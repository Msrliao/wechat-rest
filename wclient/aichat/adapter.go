package aichat

import (
	"github.com/opentdp/wechat-rest/dbase/llmodel"
	"github.com/opentdp/wechat-rest/dbase/profile"
	"github.com/opentdp/wechat-rest/dbase/setting"
	"github.com/opentdp/wechat-rest/dbase/tables"
)

func Text(id, rid, msg string) string {

	var err error
	var res string

	// 预设模型参数
	CountHistory(id)
	llmc := UserModel(id, rid)

	// 调用接口生成文本
	switch llmc.Provider {
	case "google":
		res, err = GoogleText(id, rid, msg)
	case "openai":
		res, err = OpenaiText(id, rid, msg)
	case "xunfei":
		res, err = XunfeiText(id, rid, msg)
	case "baidu":
		res, err = BaiDuText(id, rid, msg)
	default:
		res = "暂不支持此模型"
	}

	// 返回结果
	if err != nil {
		return err.Error()
	}
	return res

}

// User LLModel

type UserLLModel struct {
	RoleContext string
	*tables.LLModel
}

func UserModel(id, rid string) *UserLLModel {

	var llmc *tables.LLModel

	up, _ := profile.Fetch(&profile.FetchParam{Wxid: id, Roomid: rid})

	if up != nil {
		llmc, _ = llmodel.Fetch(&llmodel.FetchParam{Mid: up.AiModel})
	}

	if llmc == nil {
		llmc, _ = llmodel.Fetch(&llmodel.FetchParam{Mid: setting.ModelDefault})
	}

	if llmc == nil {
		llmc, _ = llmodel.Fetch(&llmodel.FetchParam{})
	}

	return &UserLLModel{LLModel: llmc, RoleContext: setting.ModelContext}

}

// Message History

type MsgHistory struct {
	Content string
	Role    string
}

var msgHistories = map[string][]*MsgHistory{}

func ResetHistory(id string) {

	msgHistories[id] = []*MsgHistory{}

}

func CountHistory(id string) int {

	if _, ok := msgHistories[id]; !ok {
		ResetHistory(id)
	}

	return len(msgHistories[id])

}

func AppendHistory(id string, items ...*MsgHistory) {

	if len(msgHistories[id]) >= setting.ModelHistory {
		msgHistories[id] = msgHistories[id][len(items):]
	}

	msgHistories[id] = append(msgHistories[id], items...)

}

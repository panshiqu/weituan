package wechat

import (
	"fmt"

	"github.com/panshiqu/weituan/db"
	"github.com/panshiqu/weituan/define"
	"github.com/panshiqu/weituan/utils"
)

// SendTemplateMessage 发送模板消息
func SendTemplateMessage(userID interface{}, templateID string, page string, formID string, emphasisKeyword string, args ...string) error {
	accessToken, err := AccessToken()
	if err != nil {
		return err
	}

	wxTemplateMessage := &define.WxTemplateMessage{
		TemplateID:      templateID,
		Page:            page,
		FormID:          formID,
		Data:            make(map[string]*define.WxTemplateValue),
		EmphasisKeyword: emphasisKeyword,
	}

	for k, v := range args {
		wxTemplateMessage.Data[fmt.Sprintf("keyword%d", k+1)] = &define.WxTemplateValue{Value: v}
	}

	if err := db.MySQL.QueryRow("SELECT OpenID FROM user WHERE UserID = ?", userID).Scan(&wxTemplateMessage.ToUser); err != nil {
		return err
	}

	wxError := &define.WxError{}
	if err := utils.HTTPJSONPostJSON(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send?access_token=%s", accessToken), wxTemplateMessage, wxError); err != nil {
		return err
	}

	if wxError.ErrCode != 0 {
		return define.NewFailure(fmt.Sprintf("%d:%s", wxError.ErrCode, wxError.ErrMsg))
	}

	return nil
}

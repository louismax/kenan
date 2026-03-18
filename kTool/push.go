package kTool

// PushFeiShuWebhookCardMsg 推送飞书卡片消息
func PushFeiShuWebhookCardMsg(url, cardId string, bizData map[string]interface{}) ([]byte, error) {
	return HTTPPostJson(url, map[string]interface{}{
		"msg_type": "interactive",
		"card": map[string]interface{}{
			"type": "template",
			"data": map[string]interface{}{
				"template_id":       cardId,
				"template_variable": bizData,
			},
		},
	}, nil)
}

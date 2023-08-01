package commission

import (
	"context"
	"encoding/json"
	"fmt"
	"net/mail"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	applangmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	applangmwpb "github.com/NpoolPlatform/message/npool/g11n/mw/v1/applang"
	commmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/commission"
	tmplmwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template"
	sendmwpb "github.com/NpoolPlatform/message/npool/third/mw/v1/send"
	tmplmwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template"
	sendmwcli "github.com/NpoolPlatform/third-middleware/pkg/client/send"
)

func Prepare(body string) (interface{}, error) {
	req := commmwpb.Commission{}
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func Apply(ctx context.Context, req interface{}) error {
	_req, ok := req.(*commmwpb.Commission)
	if !ok {
		return fmt.Errorf("invalid request")
	}

	user, err := usermwcli.GetUser(ctx, _req.AppID, _req.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("invalid user")
	}
	if _, err := mail.ParseAddress(user.EmailAddress); err != nil {
		return err
	}

	conds := &applangmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: _req.AppID},
	}
	if user.SelectedLangID != nil {
		conds.LangID = &basetypes.StringVal{Op: cruder.EQ, Value: *user.SelectedLangID}
	} else {
		conds.Main = &basetypes.BoolVal{Op: cruder.EQ, Value: true}
	}
	lang, err := applangmwcli.GetLangOnly(ctx, conds)
	if err != nil {
		return err
	}
	if lang == nil {
		return fmt.Errorf("invalid lang")
	}

	text, err := tmplmwcli.GenerateText(ctx, &tmplmwpb.GenerateTextRequest{
		AppID:     _req.AppID,
		LangID:    lang.LangID,
		Channel:   basetypes.NotifChannel_ChannelEmail,
		EventType: basetypes.UsedFor_SetCommission,
	})
	if err != nil {
		return err
	}
	if text == nil {
		return fmt.Errorf("fail generate text")
	}

	return sendmwcli.SendMessage(ctx, &sendmwpb.SendMessageRequest{
		Subject:     text.Subject,
		Content:     text.Content,
		From:        text.From,
		To:          user.EmailAddress,
		ToCCs:       text.ToCCs,
		ReplyTos:    text.ReplyTos,
		AccountType: basetypes.SignMethod_Email,
	})
}

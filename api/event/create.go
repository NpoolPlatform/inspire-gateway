package event

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"
	mgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/event"

	event1 "github.com/NpoolPlatform/inspire-gateway/pkg/event"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

//nolint
func (s *Server) CreateEvent(ctx context.Context, in *npool.CreateEventRequest) (*npool.CreateEventResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("CreateEvent", "AppID", in.GetAppID(), "Error", err)
		return &npool.CreateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	switch in.GetEventType() {
	case basetypes.UsedFor_Signup:
	case basetypes.UsedFor_Signin:
	case basetypes.UsedFor_Update:
	case basetypes.UsedFor_Contact:
	case basetypes.UsedFor_SetWithdrawAddress:
	case basetypes.UsedFor_Withdraw:
	case basetypes.UsedFor_CreateInvitationCode:
	case basetypes.UsedFor_SetCommission:
	case basetypes.UsedFor_SetTransferTargetUser:
	case basetypes.UsedFor_WithdrawalRequest:
	case basetypes.UsedFor_WithdrawalCompleted:
	case basetypes.UsedFor_DepositReceived:
	case basetypes.UsedFor_KYCApproved:
	case basetypes.UsedFor_KYCRejected:
	default:
		logger.Sugar().Errorw("CreateEvent", "EventType", in.GetEventType())
		return &npool.CreateEventResponse{}, status.Error(codes.InvalidArgument, "EventType is invalid")
	}
	for _, id := range in.GetCouponIDs() {
		if _, err := uuid.Parse(id); err != nil {
			logger.Sugar().Errorw("CreateEvent", "CouponIDs", in.GetCouponIDs(), "Error", err)
			return &npool.CreateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	if _, err := decimal.NewFromString(in.GetCredits()); err != nil {
		logger.Sugar().Errorw("CreateEvent", "Credits", in.GetCredits(), "Error", err)
		return &npool.CreateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := decimal.NewFromString(in.GetCreditsPerUSD()); err != nil {
		logger.Sugar().Errorw("CreateEvent", "CreditsPerUSD", in.GetCreditsPerUSD(), "Error", err)
		return &npool.CreateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	req := &mgrpb.EventReq{
		AppID:         &in.AppID,
		EventType:     &in.EventType,
		CouponIDs:     in.CouponIDs,
		Credits:       in.Credits,
		CreditsPerUSD: in.CreditsPerUSD,
	}

	info, err := event1.CreateEvent(ctx, req)
	if err != nil {
		logger.Sugar().Errorw("CreateEvent", "Req", req, "Error", err)
		return &npool.CreateEventResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateEventResponse{
		Info: info,
	}, nil
}

package cmppserver

import (
	"context"
	"errors"
	"fmt"

	protocol "github.com/hujm2023/go-sms-protocol"
	"github.com/hujm2023/go-sms-protocol/cmpp"
	"github.com/hujm2023/go-sms-protocol/cmpp/cmpp20"
	"github.com/hujm2023/hlog"
)

var ErrInvalidPDUAssert = errors.New("invalid pdu assert")

type CommandHandler func(ctx context.Context, pdu protocol.PDU) (resp []byte, err error)

type Dispatcher struct {
	handlers map[protocol.ICommander]CommandHandler
}

func newDisPatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[protocol.ICommander]CommandHandler),
	}
}

func (d *Dispatcher) Register(cmd protocol.ICommander, handler CommandHandler) {
	if _, ok := d.handlers[cmd]; ok {
		panic(fmt.Sprintf("%s has been registered", cmd.String()))
	}
	d.handlers[cmd] = handler
}

func (d *Dispatcher) Dispatch(ctx context.Context, data []byte) (resp []byte, err error) {
	pdu, err := cmpp20.DecodeCMPP20(data)
	if err != nil {
		return nil, fmt.Errorf("decode cmpp20 error: %w", err)
	}
	cmd := pdu.GetCommand()
	handler, ok := d.handlers[cmd]
	if !ok {
		return nil, fmt.Errorf("%s not implemented", cmd.String())
	}
	return handler(ctx, pdu)
}

var cmpp20Dispatcher = newDisPatcher()

func init() {
	cmpp20Dispatcher.Register(cmpp.CommandConnect, cmpp20Connect)
	cmpp20Dispatcher.Register(cmpp.CommandSubmit, cmpp20Submit)
	cmpp20Dispatcher.Register(cmpp.CommandActiveTest, cmpp20ActiveTest)
	cmpp20Dispatcher.Register(cmpp.CommandActiveTestResp, cmpp20ActiveTestResp)
	cmpp20Dispatcher.Register(cmpp.CommandDeliverResp, cmpp20DeliveyResp)
}

func cmpp20Connect(ctx context.Context, pdu protocol.PDU) (resp []byte, err error) {
	connect, ok := pdu.(*cmpp20.PduConnect)
	if !ok {
		return nil, ErrInvalidPDUAssert
	}
	hlog.CtxInfo(ctx, "[cmpp20Connect] user:%s", connect.SourceAddr)

	// TODO: handle auth

	return connect.GenEmptyResponse().IEncode()
}

func cmpp20Submit(ctx context.Context, pdu protocol.PDU) (resp []byte, err error) {
	// handle sumit content
	submit, ok := pdu.(*cmpp20.PduSubmit)
	if !ok {
		return nil, ErrInvalidPDUAssert
	}

	content, err := protocol.DecodeCMPPCContent(ctx, submit.MsgContent, submit.MsgFmt)
	if err != nil {
		hlog.CtxWarn(ctx, "[cmpp20Submit] decode content error: %v", err)
		return nil, nil
	}
	submitResp := submit.GenEmptyResponse().(*cmpp20.PduSubmitResp)
	submitResp.MsgID = GenMsgID()
	submitResp.Result = 0

	hlog.CtxInfo(ctx, "[cmpp20Submit] content:%s, msgID:%d", content, submitResp.MsgID)
	return submitResp.IEncode()
}

func cmpp20ActiveTest(ctx context.Context, pdu protocol.PDU) (resp []byte, err error) {
	return pdu.GenEmptyResponse().IEncode()
}

func cmpp20ActiveTestResp(ctx context.Context, pdu protocol.PDU) (resp []byte, err error) {
	hlog.CtxInfo(ctx, "[cmpp20ActiveTestResp] received an active test resp")
	return nil, nil
}

func cmpp20DeliveyResp(ctx context.Context, pdu protocol.PDU) (resp []byte, err error) {
	return nil, nil
}

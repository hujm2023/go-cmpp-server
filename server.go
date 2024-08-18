package cmppserver

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/cloudwego/netpoll"
	"github.com/hujm2023/go-sms-protocol/codec"
	"github.com/hujm2023/hlog"
	"github.com/samber/lo"
)

type CMPPServer struct {
	*codec.CMPPCodec
	handler   *Dispatcher
	eventLoop netpoll.EventLoop
}

func NewCMPPServer() *CMPPServer {
	c := &CMPPServer{
		CMPPCodec: codec.NewCMPPCodec(),
		handler:   cmpp20Dispatcher,
	}
	c.eventLoop = lo.Must(
		netpoll.NewEventLoop(
			c.onRequest,
			netpoll.WithOnPrepare(c.onPrepare),
			netpoll.WithOnConnect(c.onConnect),
			netpoll.WithReadTimeout(time.Second),
			netpoll.WithOnDisconnect(c.onDisConnect),
		),
	)

	return c
}

// Listen 阻塞监听.
func (c *CMPPServer) Listen(network, address string) error {
	listener, err := net.Listen(network, address)
	if err != nil {
		return err
	}
	hlog.Noticef("===> cmpp2.0 server listened at %s://%s", network, address)
	return c.eventLoop.Serve(listener)
}

func (c *CMPPServer) Shutdown(waitFor time.Duration) error {
	ctx, cancel := context.WithTimeout(context.TODO(), waitFor)
	defer cancel()
	return c.eventLoop.Shutdown(ctx)
}

// onPrepare means connection Connected but not initialized.
// connection is not registered into poller.
func (c *CMPPServer) onPrepare(connection netpoll.Connection) context.Context {
	hlog.Noticef("[onPrepare] remote address: %s", connection.RemoteAddr().String())
	return context.Background()
}

// onConnect means connection has Connected and been initialized.
// This connection is ready for read and write.
func (c *CMPPServer) onConnect(ctx context.Context, connection netpoll.Connection) context.Context {
	hlog.Noticef("[onConnect] remote address: %s", connection.RemoteAddr().String())
	// 注意：这里的connection里面没数据
	return nil
}

// onRequest means the first byte has beed sent to this side.
func (c *CMPPServer) onRequest(ctx context.Context, connection netpoll.Connection) error {
	reader, writer := connection.Reader(), connection.Writer()
	defer func() {
		writer.Flush()
		reader.Release()
	}()

	conn := &Connection{conn: reader}

	// 水平触发，在一个for中读全部的数据
	for {
		data, err := c.CMPPCodec.Decode(conn)
		if err != nil {
			if errors.Is(err, codec.ErrPacketNotComplete) {
				break
			}
			return err
		}

		// TODO: handle pdu
		hlog.CtxInfo(ctx, "[onRequest] read data: %+v", data)
		resp, err := c.handler.Dispatch(ctx, data)
		if err != nil {
			hlog.CtxError(ctx, "[onRequest] dispatch data error: %v", err)
			return err
		}

		if len(resp) > 0 {
			_, _ = writer.WriteBinary(resp)
		}
	}

	return nil
}

func (c *CMPPServer) onDisConnect(ctx context.Context, connection netpoll.Connection) {
	hlog.Noticef("[onDisConnect] remote address: %s", connection.RemoteAddr().String())
}

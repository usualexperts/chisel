//go:build (darwin || linux)
package tunnel

import (
	"context"
)

func (p *Proxy) runTCP(ctx context.Context) error {
	done := make(chan struct{})
	//implements missing net.ListenContext
	go func() {
		select {
		case <-ctx.Done():
			p.tcp.Close()
		case <-done:
		}
	}()

	for {
		err := preAccept(p.tcp)
		if err != nil {
			select {
			case <-ctx.Done():
				//listener closed (to check)
				err = nil
			default:
				p.Infof("PreAccept error: %s", err)
			}
			close(done)
			return err
		}

		dst, l := p.openSshChannel(ctx)
		if dst == nil {
			// SSH channel failed - likely because remote end not connectable
			if err := p.restartTCPListener() ; err != nil {
				p.Infof("Listener restart error: %s", err)
				return err
			}
		    continue
		}

		src, err := p.tcp.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				//listener closed
				err = nil
			default:
				p.Infof("Accept error: %s", err)
			}
			close(done)
			return err
		}
		go p.pipeRemoteSshAlreadyOpened(src, dst, l)
	}
}

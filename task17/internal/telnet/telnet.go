package telnet

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

func Run(ctx context.Context, host string, port int, timeout time.Duration) error {
	connCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	dialer := net.Dialer{}
	conn, err := dialer.DialContext(connCtx, "tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)
	wg := &sync.WaitGroup{}

	go func() {
		ticker := time.NewTicker(timeout / 2)
		defer ticker.Stop()

		for {
			select {
			case <-connCtx.Done():
				return
			case <-ticker.C:
				err := conn.SetDeadline(time.Now().Add(timeout))
				if err != nil {
					cancel()
				}
			}
		}
	}()

	wg.Go(func() {
		err := receive(connCtx, conn)
		if err != nil {
			cancel()
			fmt.Println(err)
		}
	})

	wg.Go(func() {
		err := send(connCtx, conn)
		if err != nil {
			cancel()
			fmt.Println(err)
			err := conn.Close()
			if err != nil {
				fmt.Println(err)
			}
		}
	})

	wg.Wait()
	return nil
}

func receive(ctx context.Context, conn net.Conn) error {
	scanner := bufio.NewScanner(conn)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if !scanner.Scan() {
				return scanner.Err()
			}

			line := scanner.Text()
			fmt.Println(line)
		}
	}
}

func send(ctx context.Context, conn net.Conn) error {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if !scanner.Scan() {
				return errors.New("EOF")
			}

			line := scanner.Text()
			_, err := conn.Write([]byte(line + "\n"))
			if err != nil {
				return err
			}
		}
	}
}
